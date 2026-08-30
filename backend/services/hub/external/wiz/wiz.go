package wiz

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"ava/hub/internal/device"
	"ava/pkg/wire"
)

const (
	Vendor = "wiz"

	MinDimming = 10
	MaxDimming = 100
	MinKelvin  = 2200
	MaxKelvin  = 6500
)

type request struct {
	Method string `json:"method"`
	Params any    `json:"params,omitempty"`
}

type pilotParams struct {
	State   *bool `json:"state,omitempty"`
	Dimming *int  `json:"dimming,omitempty"`
	Temp    *int  `json:"temp,omitempty"`
	R       *int  `json:"r,omitempty"`
	G       *int  `json:"g,omitempty"`
	B       *int  `json:"b,omitempty"`
}

type pilotResult struct {
	MAC     string `json:"mac"`
	State   bool   `json:"state"`
	Dimming int    `json:"dimming"`
	Temp    int    `json:"temp"`
	R       int    `json:"r"`
	G       int    `json:"g"`
	B       int    `json:"b"`
}

type systemConfigResult struct {
	MAC        string `json:"mac"`
	ModuleName string `json:"moduleName"`
	FWVersion  string `json:"fwVersion"`
}

type modelConfigResult struct {
	MinDimLevel int   `json:"minDimLevel"`
	CCTRange    []int `json:"cctRange"`
	WhiteRange  []int `json:"whiteRange"`
	ExtRange    []int `json:"extRange"`
}

type response struct {
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type Light struct {
	info         device.Info
	capabilities wire.Capabilities
	transport    transport
}

func New(ip string, timeout time.Duration) *Light {
	return newWith(ip, newUDPTransport(timeout))
}

func Open(spec *device.Spec) *Light {
	light := New(spec.IP, spec.Timeout)

	if len(spec.Capabilities) > 0 {
		light.capabilities = spec.Capabilities
	}

	if spec.ID != "" {
		light.info.ID = spec.ID
		light.info.MAC = spec.ID
	}

	return light
}

func newWith(ip string, t transport) *Light {
	return &Light{
		info:         device.Info{Vendor: Vendor, IP: ip},
		capabilities: tunableWhite(MinDimming, MaxDimming, MinKelvin, MaxKelvin),
		transport:    t,
	}
}

func (l *Light) Info() device.Info {
	return l.info
}

func (l *Light) Capabilities() wire.Capabilities {
	return l.capabilities
}

func (l *Light) State(ctx context.Context) (wire.State, error) {
	var result pilotResult

	if err := l.call(ctx, request{Method: "getPilot", Params: struct{}{}}, &result); err != nil {
		return nil, err
	}

	state := wire.State{wire.TraitPower: wire.Bool(result.State)}

	if l.capabilities.Has(wire.TraitBrightness) && result.Dimming > 0 {
		state[wire.TraitBrightness] = wire.Number(float64(result.Dimming))
	}

	if l.capabilities.Has(wire.TraitColorTemp) {
		state[wire.TraitColorTemp] = optionalNumber(result.Temp)
	}

	if l.capabilities.Has(wire.TraitColor) {
		state[wire.TraitColor] = optionalColor(result.R, result.G, result.B)
	}

	return state, nil
}

func optionalNumber(value int) wire.Value {
	if value <= 0 {
		return wire.Value{}
	}

	return wire.Number(float64(value))
}

func optionalColor(red, green, blue int) wire.Value {
	if red+green+blue <= 0 {
		return wire.Value{}
	}

	return wire.Text(formatHex(red, green, blue))
}

func (l *Light) Identify(ctx context.Context) (device.Info, error) {
	var result systemConfigResult

	if err := l.call(ctx, request{Method: "getSystemConfig", Params: struct{}{}}, &result); err != nil {
		return device.Info{}, err
	}

	l.info.MAC = result.MAC
	l.info.Model = result.ModuleName
	l.info.ID = result.MAC
	l.info.LastSeen = time.Now()
	l.capabilities = l.describe(ctx, result.ModuleName)

	return l.info, nil
}

func (l *Light) describe(ctx context.Context, moduleName string) wire.Capabilities {
	class := classOf(moduleName)
	if class == classSocket {
		return wire.Capabilities{device.Switch(wire.TraitPower)}
	}

	dimMin, kelvinMin, kelvinMax := float64(MinDimming), float64(MinKelvin), float64(MaxKelvin)

	var model modelConfigResult
	if err := l.call(ctx, request{Method: "getModelConfig", Params: struct{}{}}, &model); err == nil {
		if model.MinDimLevel > 0 {
			dimMin = float64(model.MinDimLevel)
		}

		if low, high, ok := kelvinBounds(&model); ok {
			kelvinMin, kelvinMax = float64(low), float64(high)
		}
	}

	capabilities := wire.Capabilities{
		device.Switch(wire.TraitPower),
		device.Bounded(wire.TraitBrightness, dimMin, MaxDimming, "%"),
	}

	if class == classDimmable {
		return capabilities
	}

	capabilities = append(capabilities, device.Bounded(wire.TraitColorTemp, kelvinMin, kelvinMax, "K"))

	if class == classColor {
		capabilities = append(capabilities, wire.Capability{
			Trait: wire.TraitColor, Kind: wire.KindColor, Access: wire.AccessReadWrite,
		})
	}

	return capabilities
}

func tunableWhite(dimMin, dimMax, kelvinMin, kelvinMax float64) wire.Capabilities {
	return wire.Capabilities{
		device.Switch(wire.TraitPower),
		device.Bounded(wire.TraitBrightness, dimMin, dimMax, "%"),
		device.Bounded(wire.TraitColorTemp, kelvinMin, kelvinMax, "K"),
	}
}

func kelvinBounds(model *modelConfigResult) (low, high int, ok bool) {
	for _, candidate := range [][]int{model.CCTRange, model.ExtRange, model.WhiteRange} {
		if len(candidate) < 2 {
			continue
		}

		low, high = candidate[0], candidate[len(candidate)-1]
		if low > 0 && high > low {
			return low, high, true
		}
	}

	return 0, 0, false
}

type class int

const (
	classTunable class = iota
	classColor
	classDimmable
	classSocket
)

func classOf(moduleName string) class {
	parts := strings.Split(moduleName, "_")
	if len(parts) < 2 {
		return classTunable
	}

	identifier := strings.ToUpper(parts[1])

	switch {
	case strings.Contains(identifier, "RGB"):
		return classColor
	case strings.Contains(identifier, "TW"):
		return classTunable
	case strings.Contains(identifier, "SOCKET"):
		return classSocket
	default:
		return classDimmable
	}
}

func (l *Light) Apply(ctx context.Context, trait wire.Trait, value wire.Value) error {
	capability, ok := l.capabilities.Find(trait)
	if !ok || !capability.Writable() {
		return device.Unsupported(Vendor, trait)
	}

	value = capability.Clamp(value)

	if err := capability.Validate(value); err != nil {
		return err
	}

	switch trait {
	case wire.TraitPower:
		on, _ := value.Bool()

		return l.setPilot(ctx, pilotParams{State: &on})
	case wire.TraitBrightness:
		level, _ := value.Number()
		dimming := int(level)
		on := true

		return l.setPilot(ctx, pilotParams{State: &on, Dimming: &dimming})
	case wire.TraitColorTemp:
		kelvin, _ := value.Number()
		temp := int(kelvin)

		return l.setPilot(ctx, pilotParams{Temp: &temp})
	case wire.TraitColor:
		hex, _ := value.Text()

		red, green, blue, err := parseHex(hex)
		if err != nil {
			return err
		}

		on := true

		return l.setPilot(ctx, pilotParams{State: &on, R: &red, G: &green, B: &blue})
	default:
		return device.Unsupported(Vendor, trait)
	}
}

func (l *Light) setPilot(ctx context.Context, params pilotParams) error {
	return l.call(ctx, request{Method: "setPilot", Params: params}, nil)
}

func (l *Light) call(ctx context.Context, req request, out any) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("wiz: encode %s: %w", req.Method, err)
	}

	raw, err := l.transport.Do(ctx, address(l.info.IP), payload)
	if err != nil {
		return err
	}

	var envelope response
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return fmt.Errorf("wiz: decode %s: %w", req.Method, err)
	}

	if envelope.Error != nil {
		return fmt.Errorf("wiz: %s: %s (%d)", req.Method, envelope.Error.Message, envelope.Error.Code)
	}

	if out == nil {
		return nil
	}

	if err := json.Unmarshal(envelope.Result, out); err != nil {
		return fmt.Errorf("wiz: decode %s result: %w", req.Method, err)
	}

	return nil
}

func parseHex(hex string) (red, green, blue int, err error) {
	if len(hex) != 7 || hex[0] != '#' {
		return 0, 0, 0, fmt.Errorf("wiz: color must look like #rrggbb, got %q", hex)
	}

	channels := [3]int{}

	for at := range channels {
		value, convErr := strconv.ParseUint(hex[1+at*2:3+at*2], 16, 8)
		if convErr != nil {
			return 0, 0, 0, fmt.Errorf("wiz: color %q is not hexadecimal", hex)
		}

		channels[at] = int(value)
	}

	return channels[0], channels[1], channels[2], nil
}

func formatHex(red, green, blue int) string {
	return fmt.Sprintf("#%02x%02x%02x", red, green, blue)
}
