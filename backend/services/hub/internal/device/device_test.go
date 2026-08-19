package device_test

import (
	"encoding/json"
	"errors"
	"testing"

	"ava/hub/internal/device"
)

func TestCapabilityHas(t *testing.T) {
	set := device.CapabilityBrightness | device.CapabilityColorTemp

	if !set.Has(device.CapabilityBrightness) {
		t.Error("expected brightness")
	}

	if !set.Has(device.CapabilityBrightness | device.CapabilityColorTemp) {
		t.Error("expected both capabilities to match together")
	}

	if set.Has(device.CapabilityColor) {
		t.Error("did not expect color")
	}

	var none device.Capability
	if none.Has(device.CapabilityBrightness) {
		t.Error("empty set should have nothing")
	}
}

func TestCapabilityMarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		set  device.Capability
		want string
	}{
		{"none", 0, `[]`},
		{"one", device.CapabilityBrightness, `["brightness"]`},
		{"tunable white", device.CapabilityBrightness | device.CapabilityColorTemp, `["brightness","color_temp"]`},
		{"all", device.CapabilityBrightness | device.CapabilityColorTemp | device.CapabilityColor, `["brightness","color_temp","color"]`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := json.Marshal(tc.set)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}

			if string(raw) != tc.want {
				t.Errorf("got %s, want %s", raw, tc.want)
			}
		})
	}
}

func TestUnsupportedIsMatchable(t *testing.T) {
	err := device.Unsupported("wiz", device.CapabilityColor)

	if !errors.Is(err, device.ErrUnsupported) {
		t.Fatal("expected errors.Is to match ErrUnsupported")
	}

	if got := err.Error(); got != "wiz: color: capability not supported" {
		t.Errorf("got %q", got)
	}
}

func TestClamp(t *testing.T) {
	tests := []struct{ in, want int }{{-10, 10}, {5, 10}, {10, 10}, {50, 50}, {100, 100}, {250, 100}}

	for _, tc := range tests {
		if got := device.Clamp(tc.in, 10, 100); got != tc.want {
			t.Errorf("Clamp(%d) = %d, want %d", tc.in, got, tc.want)
		}
	}
}
