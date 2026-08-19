package wiz

import (
	"context"
	"net"
	"testing"
	"time"
)

func startBulb(t *testing.T, dropFirst int, reply string) string {
	t.Helper()

	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	t.Cleanup(func() { conn.Close() })

	go func() {
		buf := make([]byte, 1024)
		dropped := 0

		for {
			n, addr, err := conn.ReadFrom(buf)
			if err != nil {
				return
			}

			if dropped < dropFirst {
				dropped++

				continue
			}

			_, _ = conn.WriteTo([]byte(reply), addr)
			_ = n
		}
	}()

	return conn.LocalAddr().String()
}

func TestTransportRetriesUntilTheBulbAnswers(t *testing.T) {
	addr := startBulb(t, 2, `{"result":{"state":true}}`)

	transport := newUDPTransport(300 * time.Millisecond)

	reply, err := transport.Do(context.Background(), addr, []byte(`{"method":"getPilot"}`))
	if err != nil {
		t.Fatalf("Do: %v", err)
	}

	if string(reply) != `{"result":{"state":true}}` {
		t.Errorf("reply = %s", reply)
	}
}

func TestTransportGivesUpWithAClearError(t *testing.T) {
	addr := startBulb(t, 100, "")

	transport := newUDPTransport(120 * time.Millisecond)

	_, err := transport.Do(context.Background(), addr, []byte(`{"method":"getPilot"}`))
	if err == nil {
		t.Fatal("expected an error when the bulb never answers")
	}

	if !contains(err.Error(), "no reply from") || !contains(err.Error(), "4 attempts") {
		t.Errorf("unhelpful error: %v", err)
	}
}

func TestTransportHonoursContextCancellation(t *testing.T) {
	addr := startBulb(t, 100, "")

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	start := time.Now()

	_, err := newUDPTransport(time.Second).Do(ctx, addr, []byte(`{}`))
	if err == nil {
		t.Fatal("expected an error")
	}

	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("kept retrying past the deadline: %s", elapsed)
	}
}

func contains(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}

	return false
}
