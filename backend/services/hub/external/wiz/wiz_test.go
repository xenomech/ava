package wiz

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"ava/hub/internal/device"
	"ava/pkg/wire"
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

	if err := light.Apply(context.Background(), wire.TraitPower, wire.Bool(true)); err != nil {
		t.Fatalf("Apply power: %v", err)
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

			if err := light.Apply(context.Background(), wire.TraitBrightness, wire.Number(float64(tc.in))); err != nil {
				t.Fatalf("Apply brightness: %v", err)
			}

			params := lastRequest(t, fake)["params"].(map[string]any)
			if params["dimming"] != tc.want {
				t.Errorf("dimming = %v, want %v", params["dimming"], tc.want)
			}
		})
	}
}

func TestSetBrightnessTurnsTheLightOn(t *testing.T) {
	fake := &fakeTransport{}
	light := newTestLight(fake)

	if err := light.Apply(context.Background(), wire.TraitBrightness, wire.Number(60)); err != nil {
		t.Fatalf("Apply brightness: %v", err)
	}

	params := lastRequest(t, fake)["params"].(map[string]any)
	if params["state"] != true {
		t.Errorf("state = %v, want true; a wiz bulb ignores dimming while it is off", params["state"])
	}
}

func TestSetColorTempClamps(t *testing.T) {
	fake := &fakeTransport{}
	light := newTestLight(fake)

	if err := light.Apply(context.Background(), wire.TraitColorTemp, wire.Number(1000)); err != nil {
		t.Fatalf("Apply color temp: %v", err)
	}

	if got := lastRequest(t, fake)["params"].(map[string]any)["temp"]; got != float64(MinKelvin) {
		t.Errorf("temp = %v, want %d", got, MinKelvin)
	}

	if err := light.Apply(context.Background(), wire.TraitColorTemp, wire.Number(9000)); err != nil {
		t.Fatalf("Apply color temp: %v", err)
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

	assertNumber(t, state, wire.TraitBrightness, 42)
	assertNumber(t, state, wire.TraitColorTemp, 2700)

	if on, ok := state.Get(wire.TraitPower); !ok {
		t.Error("power missing")
	} else if flag, _ := on.Bool(); !flag {
		t.Error("power = false, want true")
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

	err := newTestLight(fake).Apply(context.Background(), wire.TraitPower, wire.Bool(true))
	if err == nil {
		t.Fatal("expected an error")
	}

	if got := err.Error(); got != "wiz: setPilot: Method not found (-32601)" {
		t.Errorf("got %q", got)
	}
}

func TestCapabilitiesExcludeColor(t *testing.T) {
	capabilities := newTestLight(&fakeTransport{}).Capabilities()

	if !capabilities.Has(wire.TraitBrightness) || !capabilities.Has(wire.TraitColorTemp) {
		t.Errorf("capabilities = %v", capabilities)
	}

	if capabilities.Has(wire.TraitColor) {
		t.Error("wiz battens are tunable white only")
	}
}

func TestTraitsComeFromTheModuleName(t *testing.T) {
	tests := []struct {
		module string
		want   []wire.Trait
	}{
		{"ESP25_SHTW_01", []wire.Trait{wire.TraitPower, wire.TraitBrightness, wire.TraitColorTemp}},
		{"ESP01_SHRGB1C_31", []wire.Trait{wire.TraitPower, wire.TraitBrightness, wire.TraitColorTemp, wire.TraitColor}},
		{"ESP01_SHDW1C_31", []wire.Trait{wire.TraitPower, wire.TraitBrightness}},
		{"ESP10_SOCKET_06", []wire.Trait{wire.TraitPower}},
		{"nonsense", []wire.Trait{wire.TraitPower, wire.TraitBrightness, wire.TraitColorTemp}},
	}

	for _, tc := range tests {
		t.Run(tc.module, func(t *testing.T) {
			fake := &fakeTransport{}
			fake.replyFor = func([]byte) string { return `{"error":{"code":-32601,"message":"Method not found"}}` }

			got := newTestLight(fake).describe(context.Background(), tc.module)

			if err := got.Verify(); err != nil {
				t.Fatalf("Verify: %v", err)
			}

			if len(got) != len(tc.want) {
				t.Fatalf("%s -> %d traits, want %d", tc.module, len(got), len(tc.want))
			}

			for _, trait := range tc.want {
				if !got.Has(trait) {
					t.Errorf("%s is missing %s", tc.module, trait)
				}
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

	assertRange(t, light.Capabilities(), wire.TraitColorTemp, 2700, 6500)
	assertRange(t, light.Capabilities(), wire.TraitBrightness, 10, 100)

	if light.Capabilities().Has(wire.TraitColor) {
		t.Error("a tunable white bulb must not claim colour")
	}
}

func TestColorTempClampsToTheDeviceRangeNotTheDefault(t *testing.T) {
	fake := &fakeTransport{}
	light := newTestLight(fake)
	light.capabilities = wire.Capabilities{
		device.Switch(wire.TraitPower),
		device.Bounded(wire.TraitBrightness, 10, 100, "%"),
		device.Bounded(wire.TraitColorTemp, 2700, 6500, "K"),
	}

	if err := light.Apply(context.Background(), wire.TraitColorTemp, wire.Number(2200)); err != nil {
		t.Fatalf("Apply color temp: %v", err)
	}

	if got := lastRequest(t, fake)["params"].(map[string]any)["temp"]; got != float64(2700) {
		t.Errorf("temp = %v, want 2700 (the device floor, not the 2200 default)", got)
	}
}

func TestASocketRejectsBrightness(t *testing.T) {
	light := newTestLight(&fakeTransport{})
	light.capabilities = wire.Capabilities{device.Switch(wire.TraitPower)}

	if err := light.Apply(context.Background(), wire.TraitBrightness, wire.Number(50)); !errors.Is(err, device.ErrUnsupported) {
		t.Errorf("got %v", err)
	}
}

func TestASocketStillReportsOnlyPower(t *testing.T) {
	fake := &fakeTransport{}
	fake.replyFor = func([]byte) string {
		return `{"result":{"state":true,"dimming":100,"temp":2700}}`
	}

	light := newTestLight(fake)
	light.capabilities = wire.Capabilities{device.Switch(wire.TraitPower)}

	state, err := light.State(context.Background())
	if err != nil {
		t.Fatalf("State: %v", err)
	}

	if _, ok := state.Get(wire.TraitBrightness); ok {
		t.Error("a socket reported brightness")
	}
}

func assertNumber(t *testing.T, state wire.State, trait wire.Trait, want float64) {
	t.Helper()

	value, ok := state.Get(trait)
	if !ok {
		t.Errorf("%s missing", trait)

		return
	}

	if got, _ := value.Number(); got != want {
		t.Errorf("%s = %g, want %g", trait, got, want)
	}
}

func assertRange(t *testing.T, capabilities wire.Capabilities, trait wire.Trait, low, high float64) {
	t.Helper()

	capability, ok := capabilities.Find(trait)
	if !ok {
		t.Errorf("%s missing", trait)

		return
	}

	if capability.Min == nil || *capability.Min != low {
		t.Errorf("%s min = %v, want %g", trait, capability.Min, low)
	}

	if capability.Max == nil || *capability.Max != high {
		t.Errorf("%s max = %v, want %g", trait, capability.Max, high)
	}
}
