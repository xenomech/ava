package wiz

import (
	"testing"

	"ava/pkg/wire"
)

// What a real bulb in colour mode looks like coming out of discovery, and what
// one in white mode looks like. The sync that consumes this replaces a device's
// stored state, so anything missing here is destroyed rather than merely
// omitted.
func TestDiscoveryDescribesColourAndWhiteCompletely(t *testing.T) {
	cases := map[string]struct {
		result   pilotResult
		wantHas  []wire.Trait
		wantGone []wire.Trait
	}{
		"colour mode reports its colour and no temperature": {
			result:   pilotResult{State: true, Dimming: 80, R: 255, G: 51, B: 102},
			wantHas:  []wire.Trait{wire.TraitPower, wire.TraitBrightness, wire.TraitColor},
			wantGone: []wire.Trait{wire.TraitColorTemp},
		},
		"white mode reports its temperature and no colour": {
			result:   pilotResult{State: true, Dimming: 80, Temp: 2700},
			wantHas:  []wire.Trait{wire.TraitPower, wire.TraitBrightness, wire.TraitColorTemp},
			wantGone: []wire.Trait{wire.TraitColor},
		},
		"an off bulb still says so": {
			result:   pilotResult{State: false},
			wantHas:  []wire.Trait{wire.TraitPower},
			wantGone: []wire.Trait{wire.TraitBrightness, wire.TraitColor, wire.TraitColorTemp},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			state := pilotState(&tc.result)

			for _, trait := range tc.wantHas {
				if value, ok := state.Get(trait); !ok || !value.IsSet() {
					t.Errorf("%s is missing — the sync would erase it", trait)
				}
			}

			for _, trait := range tc.wantGone {
				if _, ok := state.Get(trait); ok {
					t.Errorf("%s should not be reported in this mode", trait)
				}
			}
		})
	}
}

// The colour a bulb reports has to survive the round trip, or the room comes
// back a slightly different shade every half minute.
func TestDiscoveryRoundTripsTheExactColour(t *testing.T) {
	state := pilotState(&pilotResult{State: true, Dimming: 80, R: 51, G: 204, B: 255})

	value, ok := state.Get(wire.TraitColor)
	if !ok {
		t.Fatal("no colour reported")
	}

	text, _ := value.Text()
	if text != "#33ccff" {
		t.Errorf("colour = %q, want #33ccff", text)
	}
}
