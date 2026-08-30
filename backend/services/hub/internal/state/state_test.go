package state_test

import (
	"os"
	"path/filepath"
	"testing"

	"ava/hub/internal/state"
	"ava/pkg/wire"
)

func bulb(id, ip string) state.KnownDevice {
	return state.KnownDevice{
		Vendor: "wiz", ID: id, IP: ip, Kind: "bulb",
		Capabilities: wire.Capabilities{{Trait: wire.TraitPower, Kind: wire.KindBool, Access: wire.AccessReadWrite}},
	}
}

func TestAnUnchangedSweepIsNotWrittenBackToDisk(t *testing.T) {
	before := []state.KnownDevice{bulb("a", "192.168.1.10"), bulb("b", "192.168.1.11")}
	after := []state.KnownDevice{bulb("a", "192.168.1.10"), bulb("b", "192.168.1.11")}

	if !state.SameDevices(before, after) {
		t.Error("an identical sweep was treated as a change, which would rewrite the state file every minute")
	}
}

func TestAChangeWorthWritingIsDetected(t *testing.T) {
	before := []state.KnownDevice{bulb("a", "192.168.1.10")}

	cases := map[string][]state.KnownDevice{
		"a device appeared":     {bulb("a", "192.168.1.10"), bulb("b", "192.168.1.11")},
		"a device vanished":     {},
		"a device moved ip":     {bulb("a", "192.168.1.99")},
		"a device changed kind": {{Vendor: "wiz", ID: "a", IP: "192.168.1.10", Kind: "plug"}},
		"capabilities changed":  {{Vendor: "wiz", ID: "a", IP: "192.168.1.10", Kind: "bulb"}},
	}

	for name, after := range cases {
		t.Run(name, func(t *testing.T) {
			if state.SameDevices(before, after) {
				t.Error("the change was missed, so the hub would forget it after a restart")
			}
		})
	}
}

func TestKnownDevicesSurviveASaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "avahub-state.json")

	saved := &state.State{
		HubID:      "hub-1",
		TenantSlug: "acme",
		Devices:    []state.KnownDevice{bulb("a", "192.168.1.10")},
	}

	if err := state.Save(path, saved); err != nil {
		t.Fatalf("Save: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}

	if mode := info.Mode().Perm(); mode != 0o600 {
		t.Errorf("mode = %o, want 600 — the file holds broker credentials", mode)
	}

	loaded, err := state.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded.Devices) != 1 || loaded.Devices[0].ID != "a" || loaded.Devices[0].IP != "192.168.1.10" {
		t.Fatalf("devices did not survive: %+v", loaded.Devices)
	}

	if !loaded.Devices[0].Capabilities.Has(wire.TraitPower) {
		t.Error("capabilities did not survive, so a recovered device would be opened with the wrong ones")
	}
}
