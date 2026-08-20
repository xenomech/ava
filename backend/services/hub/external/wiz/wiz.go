package wiz

import (
	"context"
	"encoding/json"
	"fmt"
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

type response struct {
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type Light struct {
	info      device.Info
	transport transport
}

func New(ip string, timeout time.Duration) *Light {
	return newWith(ip, newUDPTransport(timeout))
}

func newWith(ip string, t transport) *Light {
	return &Light{
		info:      device.Info{Vendor: Vendor, IP: ip},
		transport: t,
	}
}

func (l *Light) Info() device.Info {
	return l.info
}

func (l *Light) Capabilities() device.Capability {
	return device.CapabilityBrightness | device.CapabilityColorTemp
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

	return l.info, nil
}

func (l *Light) SetPower(ctx context.Context, on bool) error {
	return l.setPilot(ctx, pilotParams{State: &on})
}

func (l *Light) SetBrightness(ctx context.Context, percent int) error {
	level := device.Clamp(percent, MinDimming, MaxDimming)

	return l.setPilot(ctx, pilotParams{Dimming: &level})
}

func (l *Light) SetColorTemp(ctx context.Context, kelvin int) error {
	temp := device.Clamp(kelvin, MinKelvin, MaxKelvin)

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
