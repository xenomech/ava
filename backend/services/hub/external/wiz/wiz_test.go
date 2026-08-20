package wiz

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"ava/hub/internal/device"
)

var _ device.Device = (*Light)(nil)

type fakeTransport struct {
	sent      [][]byte
	addrs     []string
	reply     []byte
	err       error
	failFirst int
}

func (f *fakeTransport) Do(_ context.Context, addr string, payload []byte) ([]byte, error) {
	f.addrs = append(f.addrs, addr)
	f.sent = append(f.sent, payload)

	if f.failFirst > 0 {
		f.failFirst--

		return nil, errors.New("timeout")
	}

	if f.err != nil {
		return nil, f.err
	}

	if f.reply == nil {
		return []byte(`{"result":{"success":true}}`), nil
	}

	return f.reply, nil
}

func newTestLight(t *fakeTransport) *Light {
	return newWith("192.168.1.50", t)
}

func lastRequest(t *testing.T, f *fakeTransport) map[string]any {
	t.Helper()

	if len(f.sent) == 0 {
		t.Fatal("nothing was sent")
	}

	var out map[string]any
	if err := json.Unmarshal(f.sent[len(f.sent)-1], &out); err != nil {
		t.Fatalf("sent payload is not json: %v", err)
	}

	return out
}

func TestSetPowerEncodesSetPilot(t *testing.T) {
	fake := &fakeTransport{}
	light := newTestLight(fake)

	if err := light.SetPower(context.Background(), true); err != nil {
		t.Fatalf("SetPower: %v", err)
	}

	req := lastRequest(t, fake)
	if req["method"] != "setPilot" {
		t.Errorf("method = %v", req["method"])
	}

	params, _ := req["params"].(map[string]any)
	if params["state"] != true {
		t.Errorf("state = %v", params["state"])
	}

	if _, ok := params["dimming"]; ok {
		t.Error("power command must not carry dimming")
	}

	if got := fake.addrs[0]; got != "192.168.1.50:38899" {
		t.Errorf("addr = %s", got)
	}
}

func TestSetBrightnessClamps(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want float64
	}{
		{"below range", 0, 10},
		{"in range", 55, 55},
		{"above range", 140, 100},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fake := &fakeTransport{}
			light := newTestLight(fake)

			if err := light.SetBrightness(context.Background(), tc.in); err != nil {
				t.Fatalf("SetBrightness: %v", err)
			}

			params := lastRequest(t, fake)["params"].(map[string]any)
			if params["dimming"] != tc.want {
				t.Errorf("dimming = %v, want %v", params["dimming"], tc.want)
			}
		})
	}
}

func TestSetColorTempClamps(t *testing.T) {
	fake := &fakeTransport{}
	light := newTestLight(fake)

	if err := light.SetColorTemp(context.Background(), 1000); err != nil {
		t.Fatalf("SetColorTemp: %v", err)
	}

	if got := lastRequest(t, fake)["params"].(map[string]any)["temp"]; got != float64(MinKelvin) {
		t.Errorf("temp = %v, want %d", got, MinKelvin)
	}

	if err := light.SetColorTemp(context.Background(), 9000); err != nil {
		t.Fatalf("SetColorTemp: %v", err)
	}

	if got := lastRequest(t, fake)["params"].(map[string]any)["temp"]; got != float64(MaxKelvin) {
		t.Errorf("temp = %v, want %d", got, MaxKelvin)
	}
}

func TestStateParsesGetPilot(t *testing.T) {
	fake := &fakeTransport{
		reply: []byte(`{"method":"getPilot","result":{"mac":"a8bb50","state":true,"dimming":42,"temp":2700}}`),
	}

	state, err := newTestLight(fake).State(context.Background())
	if err != nil {
		t.Fatalf("State: %v", err)
	}

	if !state.Power || state.Brightness != 42 || state.ColorTemp != 2700 {
		t.Errorf("state = %+v", state)
	}

	if lastRequest(t, fake)["method"] != "getPilot" {
		t.Error("expected getPilot")
	}
}

func TestIdentifyFillsIdentity(t *testing.T) {
	fake := &fakeTransport{
		reply: []byte(`{"result":{"mac":"a8bb5006033d","moduleName":"ESP03_SHTW1C_31","fwVersion":"1.25.0"}}`),
	}

	light := newTestLight(fake)

	info, err := light.Identify(context.Background())
	if err != nil {
		t.Fatalf("Identify: %v", err)
	}

	if info.MAC != "a8bb5006033d" || info.Model != "ESP03_SHTW1C_31" || info.ID != "a8bb5006033d" {
		t.Errorf("info = %+v", info)
	}

	if light.Info().MAC != "a8bb5006033d" {
		t.Error("Identify should store what it learned on the light")
	}
}

func TestDeviceErrorIsReported(t *testing.T) {
	fake := &fakeTransport{reply: []byte(`{"error":{"code":-32601,"message":"Method not found"}}`)}

	err := newTestLight(fake).SetPower(context.Background(), true)
	if err == nil {
		t.Fatal("expected an error")
	}

	if got := err.Error(); got != "wiz: setPilot: Method not found (-32601)" {
		t.Errorf("got %q", got)
	}
}

func TestCapabilitiesExcludeColor(t *testing.T) {
	caps := newTestLight(&fakeTransport{}).Capabilities()

	if !caps.Has(device.CapabilityBrightness) || !caps.Has(device.CapabilityColorTemp) {
		t.Errorf("caps = %s", caps)
	}

	if caps.Has(device.CapabilityColor) {
		t.Error("wiz battens are tunable white only")
	}
}
