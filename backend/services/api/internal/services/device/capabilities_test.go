package device

import (
	"encoding/json"
	"errors"
	"testing"

	"ava/pkg/wire"
)

func TestCapabilitiesAreReadBackFromTheColumn(t *testing.T) {
	stored := json.RawMessage(`[
		{"trait":"power","kind":"bool","access":"rw"},
		{"trait":"fan_speed","kind":"number","access":"rw","min":1,"max":10,"step":1}
	]`)

	capabilities, err := decodeCapabilities(stored)
	if err != nil {
		t.Fatalf("decodeCapabilities: %v", err)
	}

	if err := capabilities.ValidateWrite(wire.TraitFanSpeed, wire.Number(7)); err != nil {
		t.Errorf("speed 7 rejected: %v", err)
	}

	if err := capabilities.ValidateWrite(wire.TraitFanSpeed, wire.Number(11)); !errors.Is(err, wire.ErrOutOfRange) {
		t.Errorf("speed 11 gave %v", err)
	}

	if err := capabilities.ValidateWrite(wire.TraitBrightness, wire.Number(50)); !errors.Is(err, wire.ErrUnknownTrait) {
		t.Errorf("brightness on a fan gave %v", err)
	}
}

func TestAnEmptyCapabilityColumnRejectsEveryWrite(t *testing.T) {
	for _, stored := range []json.RawMessage{nil, json.RawMessage("[]")} {
		capabilities, err := decodeCapabilities(stored)
		if err != nil {
			t.Fatalf("decodeCapabilities: %v", err)
		}

		if err := capabilities.ValidateWrite(wire.TraitPower, wire.Bool(true)); !errors.Is(err, wire.ErrUnknownTrait) {
			t.Errorf("got %v", err)
		}
	}
}

func TestUnreadableCapabilitiesDoNotPanic(t *testing.T) {
	if _, err := decodeCapabilities(json.RawMessage(`{"not":"an array"}`)); err == nil {
		t.Error("expected an error for a malformed column")
	}
}

func TestReadOnlyTraitsCannotBeCommanded(t *testing.T) {
	stored := json.RawMessage(`[{"trait":"temperature","kind":"number","access":"r","unit":"C"}]`)

	capabilities, err := decodeCapabilities(stored)
	if err != nil {
		t.Fatalf("decodeCapabilities: %v", err)
	}

	if err := capabilities.ValidateWrite(wire.TraitTemperature, wire.Number(20)); !errors.Is(err, wire.ErrReadOnly) {
		t.Errorf("got %v", err)
	}
}
