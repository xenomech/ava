package mqtt

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"
)

// The guard is worth little if it runs after the credentials are already on the wire.
func TestConnectRefusesAPublicPlaintextBrokerBeforeDialling(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	start := time.Now()

	_, err := Connect(ctx, &Options{
		BrokerURL: "tcp://broker.example.com:1883",
		ClientID:  "ava-guard-check",
		Username:  "ava-api",
		Password:  "hunter2",
	})
	if !errors.Is(err, ErrInsecureBroker) {
		t.Fatalf("Connect = %v, want ErrInsecureBroker", err)
	}

	// Reaching a public host would cost orders of magnitude more, so returning at once proves nothing was sent.
	if elapsed := time.Since(start); elapsed > 100*time.Millisecond {
		t.Errorf("refused after %v, which suggests it dialled first", elapsed)
	}
}

func TestPlaintextIsRefusedOnlyWhenItLeavesThePrivateNetwork(t *testing.T) {
	cases := []struct {
		name    string
		url     string
		refused bool
	}{
		{"loopback by name", "tcp://localhost:1883", false},
		{"loopback by address", "tcp://127.0.0.1:1883", false},
		{"ipv6 loopback", "tcp://[::1]:1883", false},
		{"compose service name", "tcp://mosquitto:1883", false},
		{"railway private network", "tcp://mosquitto.railway.internal:1883", false},
		{"rfc1918", "tcp://10.0.0.4:1883", false},
		{"rfc1918 192.168", "tcp://192.168.1.50:1883", false},
		{"mdns", "tcp://raspberrypi.local:1883", false},

		{"public hostname", "tcp://broker.example.com:1883", true},
		{"public address", "tcp://203.0.113.10:1883", true},
		{"railway public domain", "tcp://ava-broker.up.railway.app:1883", true},
		{"websocket without tls", "ws://broker.example.com:9001", true},

		{"native tls", "ssl://broker.example.com:8883", false},
		{"tls alias", "tls://broker.example.com:8883", false},
		{"mqtts alias", "mqtts://broker.example.com:8883", false},
		{"websocket over tls", "wss://broker.example.com:443", false},
		{"scheme case is ignored", "SSL://broker.example.com:8883", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := guardTransport(tc.url, false)

			if tc.refused && !errors.Is(err, ErrInsecureBroker) {
				t.Fatalf("guardTransport(%q) = %v, want ErrInsecureBroker", tc.url, err)
			}

			if !tc.refused && err != nil {
				t.Fatalf("guardTransport(%q) = %v, want nil", tc.url, err)
			}
		})
	}
}

func TestTheOverrideIsTheOnlyWayPastTheGuard(t *testing.T) {
	const public = "tcp://broker.example.com:1883"

	if err := guardTransport(public, false); !errors.Is(err, ErrInsecureBroker) {
		t.Fatalf("without the override = %v, want ErrInsecureBroker", err)
	}

	if err := guardTransport(public, true); err != nil {
		t.Fatalf("with the override = %v, want nil", err)
	}
}

func TestTheSystemTrustStoreIsUsedWhenNoAuthorityIsPinned(t *testing.T) {
	config, err := tlsConfigFor("")
	if err != nil {
		t.Fatalf("tlsConfigFor: %v", err)
	}

	// A nil config is what makes crypto/tls verify the chain and the hostname on its own.
	if config != nil {
		t.Fatal("tlsConfigFor(\"\") returned a config, which would replace the system roots")
	}
}

func TestAnUnreadableAuthorityFailsRatherThanFallingBack(t *testing.T) {
	if _, err := tlsConfigFor("/nonexistent/ca.pem"); err == nil {
		t.Fatal("tlsConfigFor accepted a missing CA file, which would silently widen trust")
	}

	empty := t.TempDir() + "/empty.pem"
	if err := os.WriteFile(empty, []byte("not a certificate"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if _, err := tlsConfigFor(empty); err == nil {
		t.Fatal("tlsConfigFor accepted a file holding no certificates")
	}
}
