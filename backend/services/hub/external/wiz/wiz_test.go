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
	replyFor  func(payload []byte) string
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

	if f.replyFor != nil {
		return []byte(f.replyFor(payload)), nil
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

func TestCapabilitiesComeFromTheModuleName(t *testing.T) {
	tests := []struct {
		module string
		want   device.Capability
	}{
		{"ESP25_SHTW_01", device.CapabilityBrightness | device.CapabilityColorTemp},
		{"ESP01_SHRGB1C_31", device.CapabilityBrightness | device.CapabilityColorTemp | device.CapabilityColor},
		{"ESP01_SHDW1C_31", device.CapabilityBrightness},
		{"ESP10_SOCKET_06", 0},
		{"nonsense", device.CapabilityBrightness | device.CapabilityColorTemp},
	}

	for _, tc := range tests {
		t.Run(tc.module, func(t *testing.T) {
			if got := capabilitiesFor(tc.module); got != tc.want {
				t.Errorf("%s -> %s, want %s", tc.module, got, tc.want)
			}
		})
	}
}

func TestIdentifyLearnsLimitsFromTheDevice(t *testing.T) {
	fake := &fakeTransport{}
	fake.reply = nil

	replies := map[string]string{
		"getSystemConfig": `{"result":{"mac":"a8bb50aabbcc","moduleName":"ESP25_SHTW_01","fwVersion":"1.38.0"}}`,
		"getModelConfig":  `{"result":{"minDimLevel":10,"cctRange":[2700,2700,6500,6500]}}`,
	}

	fake.replyFor = func(payload []byte) string {
		for method, reply := range replies {
			if contains(string(payload), method) {
				return reply
			}
		}

		return `{"result":{"success":true}}`
	}

	light := newTestLight(fake)

	if _, err := light.Identify(context.Background()); err != nil {
		t.Fatalf("Identify: %v", err)
	}

	limits := light.Limits()

	if limits.KelvinMin != 2700 || limits.KelvinMax != 6500 {
		t.Errorf("kelvin range = %d-%d, want 2700-6500", limits.KelvinMin, limits.KelvinMax)
	}

	if limits.BrightnessMin != 10 {
		t.Errorf("brightness min = %d, want 10", limits.BrightnessMin)
	}

	if light.Capabilities().Has(device.CapabilityColor) {
		t.Error("a tunable white bulb must not claim colour")
	}
}

func TestColorTempClampsToTheDeviceRangeNotTheDefault(t *testing.T) {
	fake := &fakeTransport{}
	light := newTestLight(fake)
	light.limits = device.Limits{BrightnessMin: 10, BrightnessMax: 100, KelvinMin: 2700, KelvinMax: 6500}

	if err := light.SetColorTemp(context.Background(), 2200); err != nil {
		t.Fatalf("SetColorTemp: %v", err)
	}

	if got := lastRequest(t, fake)["params"].(map[string]any)["temp"]; got != float64(2700) {
		t.Errorf("temp = %v, want 2700 (the device floor, not the 2200 default)", got)
	}
}

func TestASocketRejectsBrightness(t *testing.T) {
	light := newTestLight(&fakeTransport{})
	light.capabilities = 0

	if err := light.SetBrightness(context.Background(), 50); !errors.Is(err, device.ErrUnsupported) {
		t.Errorf("got %v", err)
	}
}
