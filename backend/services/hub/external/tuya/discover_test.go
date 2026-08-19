package tuya

import (
	"encoding/hex"
	"encoding/json"
	"testing"
)

func TestDiscoveryKeyMatchesTheDocumentedConstant(t *testing.T) {
	const want = "6c1ec8e2bb9bb59ab50b0daf649b410a"

	if got := hex.EncodeToString(discoveryKey()); got != want {
		t.Errorf("md5(%q) = %s, want %s", discoverySecret, got, want)
	}
}

func announcement(t *testing.T, version string) []byte {
	t.Helper()

	raw, err := json.Marshal(broadcast{
		IP:         "192.168.1.60",
		GatewayID:  "bf1234567890abcdef",
		ProductKey: "keyabc123",
		Version:    version,
		Active:     2,
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	return raw
}

func TestDecodePlaintextBroadcast(t *testing.T) {
	found, err := decodeBroadcast(announcement(t, "3.3"))
	if err != nil {
		t.Fatalf("decodeBroadcast: %v", err)
	}

	if found.Info.ID != "bf1234567890abcdef" || found.Info.IP != "192.168.1.60" {
		t.Errorf("info = %+v", found.Info)
	}

	if !found.Supported() {
		t.Error("3.3 must be reported as supported")
	}
}

func TestDecodeEncryptedBroadcast(t *testing.T) {
	ciphertext, err := encrypt(discoveryKey(), announcement(t, "3.3"))
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	framed, err := pack(0, commandDPQuery, append(make([]byte, retcodeLen), ciphertext...))
	if err != nil {
		t.Fatalf("pack: %v", err)
	}

	found, err := decodeBroadcast(framed)
	if err != nil {
		t.Fatalf("decodeBroadcast: %v", err)
	}

	if found.Info.ID != "bf1234567890abcdef" {
		t.Errorf("info = %+v", found.Info)
	}

	if found.Product != "keyabc123" {
		t.Errorf("product = %s", found.Product)
	}
}

func TestNewerProtocolsAreReportedButNotSupported(t *testing.T) {
	for _, version := range []string{"3.4", "3.5"} {
		found, err := decodeBroadcast(announcement(t, version))
		if err != nil {
			t.Fatalf("decodeBroadcast(%s): %v", version, err)
		}

		if found.Supported() {
			t.Errorf("%s must not be reported as supported", version)
		}

		if found.Version != version {
			t.Errorf("version = %s", found.Version)
		}
	}
}

func TestBroadcastWithoutADeviceIDIsRejected(t *testing.T) {
	if _, err := decodeBroadcast([]byte(`{"ip":"192.168.1.60"}`)); err == nil {
		t.Error("expected an error when gwId is missing")
	}
}
