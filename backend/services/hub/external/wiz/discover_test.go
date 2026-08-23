package wiz

import (
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"

	"ava/pkg/wire"
)

func replyingBulb(t *testing.T, replies []string) (net.PacketConn, *net.UDPAddr) {
	t.Helper()

	server, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen server: %v", err)
	}

	t.Cleanup(func() { server.Close() })

	target, err := net.ResolveUDPAddr("udp4", server.LocalAddr().String())
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}

	go func() {
		buf := make([]byte, 512)

		_, from, err := server.ReadFrom(buf)
		if err != nil {
			return
		}

		for _, reply := range replies {
			_, _ = server.WriteTo([]byte(reply), from)
		}
	}()

	return server, target
}

func probeAndCollect(t *testing.T, target *net.UDPAddr) []Found {
	t.Helper()

	client, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen client: %v", err)
	}

	defer client.Close()

	if err := client.SetDeadline(time.Now().Add(400 * time.Millisecond)); err != nil {
		t.Fatalf("deadline: %v", err)
	}

	probe, err := json.Marshal(request{Method: "getPilot", Params: struct{}{}})
	if err != nil {
		t.Fatalf("encode probe: %v", err)
	}

	if _, err := client.WriteTo(probe, target); err != nil {
		t.Fatalf("write: %v", err)
	}

	return collect(context.Background(), client)
}

func TestCollectDeduplicatesAndSkipsGarbage(t *testing.T) {
	const pilot = `{"method":"getPilot","result":{"mac":"a8bb50000001","state":true,"dimming":80,"temp":2700}}`

	_, target := replyingBulb(t, []string{pilot, pilot, `not json at all`})

	found := probeAndCollect(t, target)

	if len(found) != 1 {
		t.Fatalf("expected one device after dedup and filtering, got %d", len(found))
	}

	if found[0].Info.MAC != "a8bb50000001" || found[0].Info.ID != "a8bb50000001" {
		t.Errorf("info = %+v", found[0].Info)
	}

	if found[0].Info.IP != "127.0.0.1" {
		t.Errorf("ip = %s", found[0].Info.IP)
	}

	power, _ := found[0].State.Get(wire.TraitPower)
	on, _ := power.Bool()
	brightness, _ := found[0].State.Get(wire.TraitBrightness)
	level, _ := brightness.Number()
	kelvin, _ := found[0].State.Get(wire.TraitColorTemp)
	temp, _ := kelvin.Number()

	if !on || level != 80 || temp != 2700 {
		t.Errorf("state = %+v", found[0].State)
	}
}

func TestCollectIgnoresErrorReplies(t *testing.T) {
	_, target := replyingBulb(t, []string{`{"error":{"code":-1,"message":"nope"}}`})

	if found := probeAndCollect(t, target); len(found) != 0 {
		t.Errorf("an error reply is not a device: %+v", found)
	}
}

func TestDiscoverReturnsCleanlyWhenNobodyAnswers(t *testing.T) {
	found, err := Discover(context.Background(), 200*time.Millisecond)
	if err != nil {
		t.Fatalf("Discover: %v", err)
	}

	for _, f := range found {
		t.Logf("a real WiZ device answered on this LAN: %s", f.Info.IP)
	}
}
