package wiz

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"ava/hub/internal/device"
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
}

type pilotResult struct {
	MAC     string `json:"mac"`
	State   bool   `json:"state"`
	Dimming int    `json:"dimming"`
	Temp    int    `json:"temp"`
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
	capabilities device.Capability
	limits       device.Limits
	transport    transport
}

func New(ip string, timeout time.Duration) *Light {
	return newWith(ip, newUDPTransport(timeout))
}

func newWith(ip string, t transport) *Light {
	return &Light{
		info:         device.Info{Vendor: Vendor, IP: ip},
		capabilities: device.CapabilityBrightness | device.CapabilityColorTemp,
		limits:       device.Limits{}.WithDefaults(MinDimming, MaxDimming, MinKelvin, MaxKelvin),
		transport:    t,
	}
}

func (l *Light) Info() device.Info {
	return l.info
}

func (l *Light) Capabilities() device.Capability {
	return l.capabilities
}

func (l *Light) Limits() device.Limits {
	return l.limits
}

func (l *Light) State(ctx context.Context) (device.State, error) {
	var result pilotResult

	if err := l.call(ctx, request{Method: "getPilot", Params: struct{}{}}, &result); err != nil {
		return device.State{}, err
	}

	return device.State{
		Power:      result.State,
		Brightness: result.Dimming,
		ColorTemp:  result.Temp,
	}, nil
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
	l.capabilities = capabilitiesFor(result.ModuleName)

	l.limits = l.readLimits(ctx)

	return l.info, nil
}

func (l *Light) readLimits(ctx context.Context) device.Limits {
	var model modelConfigResult

	if err := l.call(ctx, request{Method: "getModelConfig", Params: struct{}{}}, &model); err != nil {
		return device.Limits{}.WithDefaults(MinDimming, MaxDimming, MinKelvin, MaxKelvin)
	}

	limits := device.Limits{BrightnessMin: model.MinDimLevel, BrightnessMax: MaxDimming}

	if low, high, ok := kelvinBounds(&model); ok {
		limits.KelvinMin = low
		limits.KelvinMax = high
	}

	if !l.capabilities.Has(device.CapabilityColorTemp) {
		limits.KelvinMin = 0
		limits.KelvinMax = 0

		return limits.WithDefaults(MinDimming, MaxDimming, 0, 0)
	}

	return limits.WithDefaults(MinDimming, MaxDimming, MinKelvin, MaxKelvin)
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

func capabilitiesFor(moduleName string) device.Capability {
	parts := strings.Split(moduleName, "_")
	if len(parts) < 2 {
		return device.CapabilityBrightness | device.CapabilityColorTemp
	}

	identifier := strings.ToUpper(parts[1])

	switch {
	case strings.Contains(identifier, "RGB"):
		return device.CapabilityBrightness | device.CapabilityColorTemp | device.CapabilityColor
	case strings.Contains(identifier, "TW"):
		return device.CapabilityBrightness | device.CapabilityColorTemp
	case strings.Contains(identifier, "SOCKET"):
		return 0
	default:
		return device.CapabilityBrightness
	}
}

func (l *Light) SetPower(ctx context.Context, on bool) error {
	return l.setPilot(ctx, pilotParams{State: &on})
}

func (l *Light) SetBrightness(ctx context.Context, percent int) error {
	if !l.capabilities.Has(device.CapabilityBrightness) {
		return device.Unsupported(Vendor, device.CapabilityBrightness)
	}

	level := device.Clamp(percent, l.limits.BrightnessMin, l.limits.BrightnessMax)
	on := true

	return l.setPilot(ctx, pilotParams{State: &on, Dimming: &level})
}

func (l *Light) SetColorTemp(ctx context.Context, kelvin int) error {
	if !l.capabilities.Has(device.CapabilityColorTemp) {
		return device.Unsupported(Vendor, device.CapabilityColorTemp)
	}

	temp := device.Clamp(kelvin, l.limits.KelvinMin, l.limits.KelvinMax)

	return l.setPilot(ctx, pilotParams{Temp: &temp})
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
