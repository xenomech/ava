package inventory

import (
	"testing"

	"ava/hub/internal/device"
	"ava/pkg/wire"
)

func TestADeviceWithNoLightCapabilitiesIsAPlug(t *testing.T) {
	cases := []struct {
		name         string
		capabilities wire.Capabilities
		want         string
	}{
		{
			name:         "socket reports only power",
			capabilities: wire.Capabilities{device.Switch(wire.TraitPower)},
			want:         KindPlug,
		},
		{
			name: "dimmable white",
			capabilities: wire.Capabilities{
				device.Switch(wire.TraitPower),
				device.Bounded(wire.TraitBrightness, 10, 100, "%"),
			},
			want: KindBulb,
		},
		{
			name: "tunable white",
			capabilities: wire.Capabilities{
				device.Switch(wire.TraitPower),
				device.Bounded(wire.TraitBrightness, 10, 100, "%"),
				device.Bounded(wire.TraitColorTemp, 2200, 6500, "K"),
			},
			want: KindBulb,
		},
		{
			name:         "nothing at all",
			capabilities: nil,
			want:         KindPlug,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := kindFor(tc.capabilities); got != tc.want {
				t.Errorf("kindFor = %s, want %s", got, tc.want)
			}
		})
	}
}
