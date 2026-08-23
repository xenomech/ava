package device_test

import (
	"errors"
	"testing"

	"ava/hub/internal/device"
	"ava/pkg/wire"
)

func TestUnsupportedIsMatchable(t *testing.T) {
	err := device.Unsupported("wiz", wire.TraitColor)

	if !errors.Is(err, device.ErrUnsupported) {
		t.Fatal("expected errors.Is to match ErrUnsupported")
	}

	if got := err.Error(); got != "wiz: color: device: trait not supported" {
		t.Errorf("got %q", got)
	}
}

func TestTheBuildersProduceWellFormedCapabilities(t *testing.T) {
	built := wire.Capabilities{
		device.Switch(wire.TraitPower),
		device.Bounded(wire.TraitBrightness, 10, 100, "%"),
		device.Reading(wire.TraitTemperature, "C"),
	}

	if err := built.Verify(); err != nil {
		t.Fatalf("Verify: %v", err)
	}

	if !built[0].Writable() || !built[1].Writable() {
		t.Error("a switch and a bounded number must be writable")
	}

	if built[2].Writable() {
		t.Error("a reading must not be writable")
	}
}
