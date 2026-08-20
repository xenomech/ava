package tuya

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"ava/hub/internal/device"
)

const (
	Vendor = "tuya"

	Version = "3.3"

	DefaultPowerDP      = "1"
	DefaultBrightnessDP = "3"
	DefaultColorTempDP  = "4"

	brightnessMin = 10
	brightnessMax = 1000

	colorTempScale = 1000

	kelvinMin = 2700
	kelvinMax = 6500
)

var ErrUnsupportedVersion = errors.New("tuya: only protocol 3.3 is supported")

type Config struct {
	ID           string
	IP           string
	Name         string
	LocalKey     string
	Capabilities device.Capability
	PowerDP      string
	BrightnessDP string
	ColorTempDP  string
	Timeout      time.Duration
}

type Device struct {
	info         device.Info
	key          []byte
	capabilities device.Capability
	powerDP      string
	brightnessDP string
	colorTempDP  string
	transport    transport

	mu       sync.Mutex
	sequence uint32
}

func New(cfg *Config) (*Device, error) {
	if cfg == nil {
		return nil, errors.New("tuya: config is required")
	}

	return newWith(cfg, newTCPTransport(cfg.Timeout))
}

func newWith(cfg *Config, t transport) (*Device, error) {
	if len(cfg.LocalKey) != 16 {
		return nil, fmt.Errorf("%w: got %d characters", errKeyLength, len(cfg.LocalKey))
	}

	return &Device{
		info: device.Info{
			ID:     cfg.ID,
			Vendor: Vendor,
			Name:   cfg.Name,
			IP:     cfg.IP,
		},
		key:          []byte(cfg.LocalKey),
		capabilities: cfg.Capabilities,
		powerDP:      orDefault(cfg.PowerDP, DefaultPowerDP),
		brightnessDP: orDefault(cfg.BrightnessDP, DefaultBrightnessDP),
		colorTempDP:  orDefault(cfg.ColorTempDP, DefaultColorTempDP),
		transport:    t,
	}, nil
}

func orDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}

func (d *Device) Info() device.Info {
	return d.info
}

func (d *Device) Capabilities() device.Capability {
	return d.capabilities
}

func (d *Device) Limits() device.Limits {
	limits := device.Limits{BrightnessMin: 0, BrightnessMax: 100}

	if d.capabilities.Has(device.CapabilityColorTemp) {
		limits.KelvinMin = kelvinMin
		limits.KelvinMax = kelvinMax
	}

	return limits
}

func (d *Device) RawState(ctx context.Context) (map[string]any, error) {
	payload := map[string]any{
		"gwId":  d.info.ID,
		"devId": d.info.ID,
		"uid":   d.info.ID,
		"t":     strconv.FormatInt(time.Now().Unix(), 10),
	}

	reply, err := d.send(ctx, commandDPQuery, payload)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (d *Device) State(ctx context.Context) (device.State, error) {
	dps, err := d.RawState(ctx)
	if err != nil {
		return device.State{}, err
	}

	state := device.State{Power: asBool(dps[d.powerDP])}

	if d.capabilities.Has(device.CapabilityBrightness) {
		if raw, ok := asInt(dps[d.brightnessDP]); ok {
			state.Brightness = scale(raw, brightnessMin, brightnessMax, 0, 100)
		}
	}

	if d.capabilities.Has(device.CapabilityColorTemp) {
		if raw, ok := asInt(dps[d.colorTempDP]); ok {
			state.ColorTemp = scale(raw, 0, colorTempScale, kelvinMin, kelvinMax)
		}
	}

	return state, nil
}

func (d *Device) SetPower(ctx context.Context, on bool) error {
	return d.control(ctx, map[string]any{d.powerDP: on})
}

func (d *Device) SetBrightness(ctx context.Context, percent int) error {
	if !d.capabilities.Has(device.CapabilityBrightness) {
		return device.Unsupported(Vendor, device.CapabilityBrightness)
	}

	level := scale(device.Clamp(percent, 0, 100), 0, 100, brightnessMin, brightnessMax)

	return d.control(ctx, map[string]any{d.brightnessDP: level})
}

func (d *Device) SetColorTemp(ctx context.Context, kelvin int) error {
	if !d.capabilities.Has(device.CapabilityColorTemp) {
		return device.Unsupported(Vendor, device.CapabilityColorTemp)
	}

	value := scale(device.Clamp(kelvin, kelvinMin, kelvinMax), kelvinMin, kelvinMax, 0, colorTempScale)

	return d.control(ctx, map[string]any{d.colorTempDP: value})
}

func (d *Device) control(ctx context.Context, dps map[string]any) error {
	payload := map[string]any{
		"devId": d.info.ID,
		"uid":   d.info.ID,
		"t":     strconv.FormatInt(time.Now().Unix(), 10),
		"dps":   dps,
	}

	_, err := d.send(ctx, commandControl, payload)

	return err
}

func (d *Device) send(ctx context.Context, command uint32, payload map[string]any) (map[string]any, error) {
	plaintext, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("tuya: encode payload: %w", err)
	}

	encoded, err := encodePayload(d.key, command, plaintext)
	if err != nil {
		return nil, err
	}

	request, err := pack(d.next(), command, encoded)
	if err != nil {
		return nil, err
	}

	raw, err := d.transport.Do(ctx, address(d.info.IP), request)
	if err != nil {
		return nil, err
	}

	reply, err := unpack(raw, true)
	if err != nil {
		return nil, err
	}

	if reply.retcode != 0 {
		return nil, fmt.Errorf("tuya: device returned status %d", reply.retcode)
	}

	return decodeDPS(d.key, reply.payload)
}

func decodeDPS(key, payload []byte) (map[string]any, error) {
	plaintext, err := decodePayload(key, payload)
	if err != nil {
		return nil, err
	}

	if len(plaintext) == 0 {
		return map[string]any{}, nil
	}

	var body struct {
		DPS map[string]any `json:"dps"`
	}

	if err := json.Unmarshal(plaintext, &body); err != nil {
		return nil, fmt.Errorf("tuya: decode reply: %w", err)
	}

	if body.DPS == nil {
		return map[string]any{}, nil
	}

	return body.DPS, nil
}

func (d *Device) next() uint32 {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.sequence++

	return d.sequence
}

func asBool(value any) bool {
	on, _ := value.(bool)

	return on
}

func asInt(value any) (int, bool) {
	switch typed := value.(type) {
	case float64:
		return int(typed), true
	case int:
		return typed, true
	default:
		return 0, false
	}
}

func scale(value, fromLow, fromHigh, toLow, toHigh int) int {
	if fromHigh == fromLow {
		return toLow
	}

	span := float64(value-fromLow) / float64(fromHigh-fromLow)

	return device.Clamp(toLow+int(span*float64(toHigh-toLow)+0.5), min(toLow, toHigh), max(toLow, toHigh))
}
