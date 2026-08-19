package tuya

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"testing"
)

var testKey = []byte("0123456789abcdef")

func TestPackLayoutMatchesTheProtocol(t *testing.T) {
	payload := []byte("hello")

	got, err := pack(1, commandDPQuery, payload)
	if err != nil {
		t.Fatalf("pack: %v", err)
	}

	if len(got) != headerLen+len(payload)+tailLen {
		t.Fatalf("length = %d", len(got))
	}

	if prefix := binary.BigEndian.Uint32(got[0:4]); prefix != 0x000055AA {
		t.Errorf("prefix = %08x", prefix)
	}

	if seq := binary.BigEndian.Uint32(got[4:8]); seq != 1 {
		t.Errorf("sequence = %d", seq)
	}

	if cmd := binary.BigEndian.Uint32(got[8:12]); cmd != 0x0a {
		t.Errorf("command = %02x", cmd)
	}

	if declared := binary.BigEndian.Uint32(got[12:16]); declared != uint32(len(payload)+tailLen) {
		t.Errorf("declared length = %d, want %d", declared, len(payload)+tailLen)
	}

	if suffix := binary.BigEndian.Uint32(got[len(got)-4:]); suffix != 0x0000AA55 {
		t.Errorf("suffix = %08x", suffix)
	}
}

func TestUnpackReadsBackAPackedFrame(t *testing.T) {
	body := append(make([]byte, retcodeLen), []byte("payload bytes")...)

	raw, err := pack(7, commandStatus, body)
	if err != nil {
		t.Fatalf("pack: %v", err)
	}

	got, unpackErr := unpack(raw, true)
	if unpackErr != nil {
		t.Fatalf("unpack: %v", unpackErr)
	}

	if got.sequence != 7 || got.command != commandStatus {
		t.Errorf("frame = %+v", got)
	}

	if string(got.payload) != "payload bytes" {
		t.Errorf("payload = %q", got.payload)
	}
}

func TestUnpackRejectsCorruptFrames(t *testing.T) {
	body := append(make([]byte, retcodeLen), []byte("payload")...)
	good, err := pack(1, commandStatus, body)
	if err != nil {
		t.Fatalf("pack: %v", err)
	}

	tests := []struct {
		name string
		data []byte
		want error
	}{
		{"truncated", good[:8], errShortFrame},
		{"bad prefix", flip(good, 1), errBadPrefix},
		{"bad suffix", flip(good, len(good)-1), errBadSuffix},
		{"bad checksum", flip(good, headerLen+retcodeLen+1), errBadChecksum},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := unpack(tc.data, true); !errors.Is(err, tc.want) {
				t.Errorf("got %v, want %v", err, tc.want)
			}
		})
	}
}

func TestUnpackRejectsAnImplausibleLength(t *testing.T) {
	raw, err := pack(1, commandStatus, append(make([]byte, retcodeLen), 'x'))
	if err != nil {
		t.Fatalf("pack: %v", err)
	}

	binary.BigEndian.PutUint32(raw[12:16], 1<<20)

	if _, err := unpack(raw, true); !errors.Is(err, errPayloadLimit) {
		t.Errorf("got %v", err)
	}
}

func TestEncryptRoundTrip(t *testing.T) {
	plaintext := []byte(`{"dps":{"1":true}}`)

	ciphertext, err := encrypt(testKey, plaintext)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	if len(ciphertext)%16 != 0 {
		t.Errorf("ciphertext is not block aligned: %d", len(ciphertext))
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Error("ciphertext equals plaintext")
	}

	back, err := decrypt(testKey, ciphertext)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}

	if !bytes.Equal(back, plaintext) {
		t.Errorf("round trip = %q", back)
	}
}

func TestEncryptIsECBSoIdenticalBlocksRepeat(t *testing.T) {
	block := []byte("SIXTEEN BYTES!!!")

	ciphertext, err := encrypt(testKey, append(append([]byte{}, block...), block...))
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	if !bytes.Equal(ciphertext[0:16], ciphertext[16:32]) {
		t.Error("ECB must map identical plaintext blocks to identical ciphertext blocks")
	}
}

func TestEncryptMatchesOpenSSL(t *testing.T) {
	ciphertext, err := encrypt(testKey, []byte("SIXTEEN BYTES!!!"))
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	const want = "de77dfa289d49cdcea0645e2f8c68308"

	if got := hex.EncodeToString(ciphertext[:16]); got != want {
		t.Errorf("AES-ECB block = %s, want %s", got, want)
	}
}

func TestDecryptRejectsBadPadding(t *testing.T) {
	if _, err := decrypt(testKey, make([]byte, 16)); !errors.Is(err, errBadPadding) {
		t.Errorf("got %v", err)
	}

	if _, err := decrypt(testKey, []byte("short")); !errors.Is(err, errBadPadding) {
		t.Errorf("got %v", err)
	}
}

func TestVersionHeaderOnlyOnControl(t *testing.T) {
	plaintext := []byte(`{"dps":{"1":true}}`)

	control, err := encodePayload(testKey, commandControl, plaintext)
	if err != nil {
		t.Fatalf("encode control: %v", err)
	}

	if !bytes.HasPrefix(control, versionHeader) {
		t.Error("CONTROL must carry the 3.3 version header")
	}

	query, err := encodePayload(testKey, commandDPQuery, plaintext)
	if err != nil {
		t.Fatalf("encode query: %v", err)
	}

	if bytes.HasPrefix(query, []byte("3.3")) {
		t.Error("DP_QUERY is in NO_PROTOCOL_HEADER_CMDS and must not carry the header")
	}

	if len(versionHeader) != 15 {
		t.Errorf("version header is %d bytes, want 15", len(versionHeader))
	}
}

func TestDecodePayloadStripsTheHeader(t *testing.T) {
	plaintext := []byte(`{"dps":{"1":false}}`)

	encoded, err := encodePayload(testKey, commandControl, plaintext)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	got, err := decodePayload(testKey, encoded)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	if !bytes.Equal(got, plaintext) {
		t.Errorf("got %q", got)
	}
}

func TestDecodePayloadHandlesShortAndEmptyInput(t *testing.T) {
	if got, err := decodePayload(testKey, nil); err != nil || got != nil {
		t.Errorf("empty payload should decode to nothing, got %q %v", got, err)
	}

	if _, err := decodePayload(testKey, []byte("3.3")); !errors.Is(err, errBadPadding) {
		t.Errorf("a truncated header must not panic, got %v", err)
	}
}

func flip(data []byte, at int) []byte {
	out := append([]byte{}, data...)
	out[at] ^= 0xff

	return out
}

func TestUnpackRequestFramesHaveNoRetcode(t *testing.T) {
	payload := []byte("request body")

	raw, err := pack(3, commandControl, payload)
	if err != nil {
		t.Fatalf("pack: %v", err)
	}

	asRequest, err := unpack(raw, false)
	if err != nil {
		t.Fatalf("unpack request: %v", err)
	}

	if string(asRequest.payload) != "request body" {
		t.Errorf("request payload = %q", asRequest.payload)
	}

	asResponse, err := unpack(raw, true)
	if err != nil {
		t.Fatalf("unpack response: %v", err)
	}

	if string(asResponse.payload) == "request body" {
		t.Error("reading a request as a response should swallow four bytes as a retcode")
	}
}
