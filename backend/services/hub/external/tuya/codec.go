package tuya

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
)

const (
	prefixValue uint32 = 0x000055AA
	suffixValue uint32 = 0x0000AA55

	headerLen  = 16
	retcodeLen = 4
	tailLen    = 8

	maxPayloadLen = 1 << 16
)

const (
	commandControl   uint32 = 0x07
	commandStatus    uint32 = 0x08
	commandHeartbeat uint32 = 0x09
	commandDPQuery   uint32 = 0x0a
)

var versionHeader = append([]byte("3.3"), make([]byte, 12)...)

var noVersionHeader = map[uint32]bool{
	commandDPQuery:   true,
	commandHeartbeat: true,
}

var (
	errShortFrame   = errors.New("tuya: frame is too short")
	errBadPrefix    = errors.New("tuya: frame prefix is wrong")
	errBadSuffix    = errors.New("tuya: frame suffix is wrong")
	errBadChecksum  = errors.New("tuya: frame checksum does not match")
	errBadPadding   = errors.New("tuya: decrypted payload has invalid padding")
	errKeyLength    = errors.New("tuya: local key must be 16 bytes")
	errPayloadLimit = errors.New("tuya: frame claims an implausible payload length")
)

type frame struct {
	sequence uint32
	command  uint32
	retcode  uint32
	payload  []byte
}

func pack(sequence, command uint32, payload []byte) ([]byte, error) {
	length := len(payload) + tailLen
	if length > maxPayloadLen {
		return nil, fmt.Errorf("%w: %d bytes", errPayloadLimit, len(payload))
	}

	out := make([]byte, headerLen, headerLen+len(payload)+tailLen)

	binary.BigEndian.PutUint32(out[0:4], prefixValue)
	binary.BigEndian.PutUint32(out[4:8], sequence)
	binary.BigEndian.PutUint32(out[8:12], command)
	binary.BigEndian.PutUint32(out[12:16], uint32(length))

	out = append(out, payload...)

	checksum := crc32.ChecksumIEEE(out)
	out = binary.BigEndian.AppendUint32(out, checksum)
	out = binary.BigEndian.AppendUint32(out, suffixValue)

	return out, nil
}

func unpack(data []byte, hasRetcode bool) (frame, error) {
	if len(data) < headerLen {
		return frame{}, errShortFrame
	}

	if binary.BigEndian.Uint32(data[0:4]) != prefixValue {
		return frame{}, errBadPrefix
	}

	length := int(binary.BigEndian.Uint32(data[12:16]))
	if length > maxPayloadLen {
		return frame{}, errPayloadLimit
	}

	prefixLen := 0
	if hasRetcode {
		prefixLen = retcodeLen
	}

	total := headerLen + length
	if len(data) < total || length < prefixLen+tailLen {
		return frame{}, errShortFrame
	}

	if binary.BigEndian.Uint32(data[total-4:total]) != suffixValue {
		return frame{}, errBadSuffix
	}

	want := crc32.ChecksumIEEE(data[:total-tailLen])
	if got := binary.BigEndian.Uint32(data[total-tailLen : total-4]); got != want {
		return frame{}, fmt.Errorf("%w: got %08x want %08x", errBadChecksum, got, want)
	}

	parsed := frame{
		sequence: binary.BigEndian.Uint32(data[4:8]),
		command:  binary.BigEndian.Uint32(data[8:12]),
		payload:  data[headerLen+prefixLen : total-tailLen],
	}

	if hasRetcode {
		parsed.retcode = binary.BigEndian.Uint32(data[headerLen : headerLen+retcodeLen])
	}

	return parsed, nil
}

func encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errKeyLength, err)
	}

	padded := pad(plaintext, block.BlockSize())
	out := make([]byte, len(padded))

	for at := 0; at < len(padded); at += block.BlockSize() {
		block.Encrypt(out[at:at+block.BlockSize()], padded[at:at+block.BlockSize()])
	}

	return out, nil
}

func decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errKeyLength, err)
	}

	if len(ciphertext) == 0 || len(ciphertext)%block.BlockSize() != 0 {
		return nil, errBadPadding
	}

	out := make([]byte, len(ciphertext))

	for at := 0; at < len(ciphertext); at += block.BlockSize() {
		block.Decrypt(out[at:at+block.BlockSize()], ciphertext[at:at+block.BlockSize()])
	}

	return unpad(out, block.BlockSize())
}

func pad(data []byte, size int) []byte {
	missing := size - len(data)%size
	if missing < 1 || missing > 0xff {
		return data
	}

	return append(data, bytes.Repeat([]byte{byte(missing)}, missing)...)
}

func unpad(data []byte, size int) ([]byte, error) {
	if len(data) == 0 {
		return nil, errBadPadding
	}

	count := int(data[len(data)-1])
	if count == 0 || count > size || count > len(data) {
		return nil, errBadPadding
	}

	for _, b := range data[len(data)-count:] {
		if int(b) != count {
			return nil, errBadPadding
		}
	}

	return data[:len(data)-count], nil
}

func encodePayload(key []byte, command uint32, plaintext []byte) ([]byte, error) {
	ciphertext, err := encrypt(key, plaintext)
	if err != nil {
		return nil, err
	}

	if noVersionHeader[command] {
		return ciphertext, nil
	}

	return append(append([]byte{}, versionHeader...), ciphertext...), nil
}

func decodePayload(key, payload []byte) ([]byte, error) {
	if len(payload) >= len(versionHeader) && bytes.HasPrefix(payload, []byte("3.3")) {
		payload = payload[len(versionHeader):]
	}

	if len(payload) == 0 {
		return nil, nil
	}

	return decrypt(key, payload)
}
