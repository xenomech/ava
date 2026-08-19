package tuya

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"ava/hub/internal/device"
)

var _ device.Device = (*Device)(nil)

type fakeTransport struct {
	requests [][]byte
	reply    func(request []byte) ([]byte, error)
}

func (f *fakeTransport) Do(_ context.Context, _ string, request []byte) ([]byte, error) {
	f.requests = append(f.requests, request)

	if f.reply == nil {
		return replyWith(testKey, map[string]any{})
	}

	return f.reply(request)
}

func replyWith(key []byte, dps map[string]any) ([]byte, error) {
	plaintext, err := json.Marshal(map[string]any{"dps": dps})
	if err != nil {
		return nil, err
	}

	ciphertext, err := encrypt(key, plaintext)
	if err != nil {
		return nil, err
	}

	return pack(1, commandStatus, append(make([]byte, retcodeLen), ciphertext...))
}

func newTestDevice(t *testing.T, caps device.Capability, fake *fakeTransport) *Device {
	t.Helper()

	dev, err := newWith(&Config{
		ID:           "bf1234567890",
		IP:           "192.168.1.60",
		LocalKey:     string(testKey),
		Capabilities: caps,
	}, fake)
	if err != nil {
		t.Fatalf("newWith: %v", err)
	}

	return dev
}

func sentDPS(t *testing.T, fake *fakeTransport) map[string]any {
	t.Helper()

	if len(fake.requests) == 0 {
		t.Fatal("nothing was sent")
	}

	frame, err := unpack(fake.requests[len(fake.requests)-1], false)
	if err != nil {
		t.Fatalf("our own frame does not unpack: %v", err)
	}

	plaintext, err := decodePayload(testKey, frame.payload)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	var body struct {
		DPS map[string]any `json:"dps"`
	}

	if err := json.Unmarshal(plaintext, &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	return body.DPS
}

func TestNewRejectsAWrongLengthKey(t *testing.T) {
	_, err := newWith(&Config{ID: "x", LocalKey: "too-short"}, &fakeTransport{})

	if !errors.Is(err, errKeyLength) {
		t.Errorf("got %v", err)
	}
}

func TestSetPowerSendsTheSwitchDPS(t *testing.T) {
	fake := &fakeTransport{}

	if err := newTestDevice(t, 0, fake).SetPower(context.Background(), true); err != nil {
		t.Fatalf("SetPower: %v", err)
	}

	if got := sentDPS(t, fake)[DefaultPowerDP]; got != true {
		t.Errorf("dps[%s] = %v", DefaultPowerDP, got)
	}

	frame, err := unpack(fake.requests[0], false)
	if err != nil {
		t.Fatalf("unpack: %v", err)
	}

	if frame.command != commandControl {
		t.Errorf("command = %d, want CONTROL", frame.command)
	}
}

func TestPowerOnlyDeviceRejectsBrightnessAndColorTemp(t *testing.T) {
	dev := newTestDevice(t, 0, &fakeTransport{})

	if err := dev.SetBrightness(context.Background(), 50); !errors.Is(err, device.ErrUnsupported) {
		t.Errorf("SetBrightness: got %v", err)
	}

	if err := dev.SetColorTemp(context.Background(), 3000); !errors.Is(err, device.ErrUnsupported) {
		t.Errorf("SetColorTemp: got %v", err)
	}
}

func TestSetBrightnessScalesPercentToTuyaRange(t *testing.T) {
	tests := []struct {
		percent int
		want    float64
	}{
		{0, brightnessMin},
		{50, 505},
		{100, brightnessMax},
		{140, brightnessMax},
	}

	for _, tc := range tests {
		fake := &fakeTransport{}
		dev := newTestDevice(t, device.CapabilityBrightness, fake)

		if err := dev.SetBrightness(context.Background(), tc.percent); err != nil {
			t.Fatalf("SetBrightness(%d): %v", tc.percent, err)
		}

		if got := sentDPS(t, fake)[DefaultBrightnessDP]; got != tc.want {
			t.Errorf("percent %d -> %v, want %v", tc.percent, got, tc.want)
		}
	}
}

func TestStateParsesTheDPSMap(t *testing.T) {
	fake := &fakeTransport{
		reply: func([]byte) ([]byte, error) {
			return replyWith(testKey, map[string]any{
				DefaultPowerDP:      true,
				DefaultBrightnessDP: 505,
				DefaultColorTempDP:  500,
			})
		},
	}

	dev := newTestDevice(t, device.CapabilityBrightness|device.CapabilityColorTemp, fake)

	state, err := dev.State(context.Background())
	if err != nil {
		t.Fatalf("State: %v", err)
	}

	if !state.Power {
		t.Error("expected power on")
	}

	if state.Brightness != 50 {
		t.Errorf("brightness = %d, want 50", state.Brightness)
	}

	if state.ColorTemp != 4600 {
		t.Errorf("color temp = %d, want 4600", state.ColorTemp)
	}
}

func TestStateIgnoresCapabilitiesTheDeviceDoesNotHave(t *testing.T) {
	fake := &fakeTransport{
		reply: func([]byte) ([]byte, error) {
			return replyWith(testKey, map[string]any{DefaultPowerDP: true, DefaultBrightnessDP: 900})
		},
	}

	state, err := newTestDevice(t, 0, fake).State(context.Background())
	if err != nil {
		t.Fatalf("State: %v", err)
	}

	if state.Brightness != 0 {
		t.Errorf("a power-only plug must not report brightness, got %d", state.Brightness)
	}
}

func TestRawStateExposesEveryDPSForUnknownModels(t *testing.T) {
	fake := &fakeTransport{
		reply: func([]byte) ([]byte, error) {
			return replyWith(testKey, map[string]any{"1": true, "9": 42, "101": "scene"})
		},
	}

	dps, err := newTestDevice(t, 0, fake).RawState(context.Background())
	if err != nil {
		t.Fatalf("RawState: %v", err)
	}

	if len(dps) != 3 || dps["101"] != "scene" {
		t.Errorf("dps = %v", dps)
	}
}

func TestQueryUsesDPQueryWithoutTheVersionHeader(t *testing.T) {
	fake := &fakeTransport{}

	if _, err := newTestDevice(t, 0, fake).RawState(context.Background()); err != nil {
		t.Fatalf("RawState: %v", err)
	}

	frame, err := unpack(fake.requests[0], false)
	if err != nil {
		t.Fatalf("unpack: %v", err)
	}

	if frame.command != commandDPQuery {
		t.Errorf("command = %d, want DP_QUERY", frame.command)
	}

	if len(frame.payload) > 3 && string(frame.payload[:3]) == "3.3" {
		t.Error("DP_QUERY must not carry the version header")
	}
}

func TestSequenceNumberIncrements(t *testing.T) {
	fake := &fakeTransport{}
	dev := newTestDevice(t, 0, fake)

	for range 3 {
		if err := dev.SetPower(context.Background(), true); err != nil {
			t.Fatalf("SetPower: %v", err)
		}
	}

	for i, raw := range fake.requests {
		frame, err := unpack(raw, false)
		if err != nil {
			t.Fatalf("unpack: %v", err)
		}

		if frame.sequence != uint32(i+1) {
			t.Errorf("request %d has sequence %d", i, frame.sequence)
		}
	}
}
