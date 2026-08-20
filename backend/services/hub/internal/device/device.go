package device

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Capability uint8

const (
	CapabilityBrightness Capability = 1 << iota
	CapabilityColorTemp
	CapabilityColor
)

var capabilityNames = []struct {
	capability Capability
	name       string
}{
	{CapabilityBrightness, "brightness"},
	{CapabilityColorTemp, "color_temp"},
	{CapabilityColor, "color"},
}

func (c Capability) Has(other Capability) bool {
	return c&other == other
}

func (c Capability) Names() []string {
	names := make([]string, 0, len(capabilityNames))

	for _, entry := range capabilityNames {
		if c.Has(entry.capability) {
			names = append(names, entry.name)
		}
	}

	return names
}

func (c Capability) String() string {
	names := c.Names()
	if len(names) == 0 {
		return "none"
	}

	out := names[0]
	for _, name := range names[1:] {
		out += "," + name
	}

	return out
}

func (c Capability) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Names())
}

type Info struct {
	ID       string    `json:"id"`
	Vendor   string    `json:"vendor"`
	Name     string    `json:"name"`
	Model    string    `json:"model,omitempty"`
	IP       string    `json:"ip,omitempty"`
	MAC      string    `json:"mac,omitempty"`
	LastSeen time.Time `json:"last_seen,omitempty"`
}

type State struct {
	Power      bool `json:"power"`
	Brightness int  `json:"brightness,omitempty"`
	ColorTemp  int  `json:"color_temp,omitempty"`
}

type Device interface {
	Info() Info
	Capabilities() Capability
	State(ctx context.Context) (State, error)
	SetPower(ctx context.Context, on bool) error
	SetBrightness(ctx context.Context, percent int) error
	SetColorTemp(ctx context.Context, kelvin int) error
}

var ErrUnsupported = errors.New("capability not supported")

func Unsupported(vendor string, capability Capability) error {
	return fmt.Errorf("%s: %s: %w", vendor, capability, ErrUnsupported)
}

func Clamp(value, low, high int) int {
	if value < low {
		return low
	}

	if value > high {
		return high
	}

	return value
}
