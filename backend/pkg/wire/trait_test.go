package wire_test

import (
	"encoding/json"
	"errors"
	"testing"

	"ava/pkg/wire"
)

func num(v float64) *float64 { return &v }

func TestAFanDeclaresItsOwnSpeedRange(t *testing.T) {
	fan := wire.Capability{
		Trait: wire.TraitFanSpeed, Kind: wire.KindNumber, Access: wire.AccessReadWrite,
		Min: num(1), Max: num(10), Step: num(1),
	}

	if err := fan.Verify(); err != nil {
		t.Fatalf("Verify: %v", err)
	}

	for _, speed := range []float64{1, 7, 10} {
		if err := fan.Validate(wire.Number(speed)); err != nil {
			t.Errorf("speed %g rejected: %v", speed, err)
		}
	}

	if err := fan.Validate(wire.Number(11)); !errors.Is(err, wire.ErrOutOfRange) {
		t.Errorf("speed 11 gave %v", err)
	}

	if err := fan.Validate(wire.Number(7.5)); !errors.Is(err, wire.ErrOffStep) {
		t.Errorf("speed 7.5 gave %v", err)
	}
}

func TestAnOvenModeIsAnEnum(t *testing.T) {
	mode := wire.Capability{
		Trait: wire.TraitMode, Kind: wire.KindEnum, Access: wire.AccessReadWrite,
		Values: []string{"bake", "grill", "fan", "defrost"},
	}

	if err := mode.Validate(wire.Text("grill")); err != nil {
		t.Errorf("grill rejected: %v", err)
	}

	if err := mode.Validate(wire.Text("broil")); !errors.Is(err, wire.ErrNotAMember) {
		t.Errorf("broil gave %v", err)
	}

	if err := mode.Validate(wire.Number(2)); !errors.Is(err, wire.ErrWrongRepr) {
		t.Errorf("a number gave %v", err)
	}
}

func TestAKettleTemperatureIsReadOnlyButItsTargetIsNot(t *testing.T) {
	kettle := wire.Capabilities{
		{Trait: wire.TraitPower, Kind: wire.KindBool, Access: wire.AccessReadWrite},
		{Trait: wire.TraitTargetTemp, Kind: wire.KindNumber, Access: wire.AccessReadWrite, Min: num(40), Max: num(100), Unit: "C"},
		{Trait: wire.TraitTemperature, Kind: wire.KindNumber, Access: wire.AccessRead, Unit: "C"},
	}

	if err := kettle.Verify(); err != nil {
		t.Fatalf("Verify: %v", err)
	}

	if err := kettle.ValidateWrite(wire.TraitTargetTemp, wire.Number(80)); err != nil {
		t.Errorf("setting the target rejected: %v", err)
	}

	if err := kettle.ValidateWrite(wire.TraitTemperature, wire.Number(80)); !errors.Is(err, wire.ErrReadOnly) {
		t.Errorf("writing a measurement gave %v", err)
	}

	if err := kettle.ValidateWrite(wire.TraitBrightness, wire.Number(50)); !errors.Is(err, wire.ErrUnknownTrait) {
		t.Errorf("an absent trait gave %v", err)
	}
}

func TestATraitTheSystemHasNeverHeardOfStillWorks(t *testing.T) {
	custom := wire.Capabilities{
		{Trait: "oven:pyrolytic_clean", Kind: wire.KindBool, Access: wire.AccessReadWrite},
	}

	if err := custom.Verify(); err != nil {
		t.Fatalf("Verify: %v", err)
	}

	if err := custom.ValidateWrite("oven:pyrolytic_clean", wire.Bool(true)); err != nil {
		t.Errorf("a vendor trait was rejected: %v", err)
	}
}

func TestACapabilityThatContradictsItselfIsRejected(t *testing.T) {
	cases := map[string]wire.Capability{
		"no trait":      {Kind: wire.KindBool, Access: wire.AccessRead},
		"min above max": {Trait: wire.TraitBrightness, Kind: wire.KindNumber, Access: wire.AccessReadWrite, Min: num(100), Max: num(10)},
		"zero step":     {Trait: wire.TraitBrightness, Kind: wire.KindNumber, Access: wire.AccessReadWrite, Step: num(0)},
		"empty enum":    {Trait: wire.TraitMode, Kind: wire.KindEnum, Access: wire.AccessReadWrite},
		"unknown kind":  {Trait: wire.TraitMode, Kind: "blob", Access: wire.AccessReadWrite},
		"no access":     {Trait: wire.TraitPower, Kind: wire.KindBool},
	}

	for name, capability := range cases {
		t.Run(name, func(t *testing.T) {
			if err := capability.Verify(); !errors.Is(err, wire.ErrBadCapability) {
				t.Errorf("got %v", err)
			}
		})
	}
}

func TestTheSameTraitCannotBeDeclaredTwice(t *testing.T) {
	twice := wire.Capabilities{
		{Trait: wire.TraitPower, Kind: wire.KindBool, Access: wire.AccessReadWrite},
		{Trait: wire.TraitPower, Kind: wire.KindBool, Access: wire.AccessRead},
	}

	if err := twice.Verify(); !errors.Is(err, wire.ErrBadCapability) {
		t.Errorf("got %v", err)
	}
}

func TestClampPullsANumberInsideTheDeclaredRange(t *testing.T) {
	brightness := wire.Capability{
		Trait: wire.TraitBrightness, Kind: wire.KindNumber, Access: wire.AccessReadWrite,
		Min: num(10), Max: num(100),
	}

	if got, _ := brightness.Clamp(wire.Number(3)).Number(); got != 10 {
		t.Errorf("low clamp = %g, want 10", got)
	}

	if got, _ := brightness.Clamp(wire.Number(140)).Number(); got != 100 {
		t.Errorf("high clamp = %g, want 100", got)
	}
}

func TestStateRoundTripsAsPlainJSON(t *testing.T) {
	state := wire.State{
		wire.TraitPower:    wire.Bool(true),
		wire.TraitFanSpeed: wire.Number(7),
		wire.TraitMode:     wire.Text("grill"),
	}

	encoded, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var decoded wire.State
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if on, ok := decoded[wire.TraitPower].Bool(); !ok || !on {
		t.Errorf("power did not survive: %v", decoded[wire.TraitPower])
	}

	if speed, ok := decoded[wire.TraitFanSpeed].Number(); !ok || speed != 7 {
		t.Errorf("fan speed did not survive: %v", decoded[wire.TraitFanSpeed])
	}

	if mode, ok := decoded[wire.TraitMode].Text(); !ok || mode != "grill" {
		t.Errorf("mode did not survive: %v", decoded[wire.TraitMode])
	}
}

func TestAnUnsetValueIsNotZero(t *testing.T) {
	var missing wire.Value

	if missing.IsSet() {
		t.Error("the zero value claims to be set")
	}

	if _, ok := missing.Number(); ok {
		t.Error("an unset value produced a number")
	}

	encoded, err := json.Marshal(missing)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if string(encoded) != "null" {
		t.Errorf("unset marshalled to %s", encoded)
	}
}

func TestAPlugReportsNoBrightnessAtAll(t *testing.T) {
	plug := wire.State{wire.TraitPower: wire.Bool(false)}

	encoded, err := json.Marshal(plug)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if string(encoded) != `{"power":false}` {
		t.Errorf("plug state = %s, want only power", encoded)
	}

	if _, ok := plug.Get(wire.TraitBrightness); ok {
		t.Error("a plug reported brightness")
	}
}

func TestANullTraitIsARetractionRatherThanAValue(t *testing.T) {
	var reported wire.State
	if err := json.Unmarshal([]byte(`{"power":true,"color":null,"brightness":40}`), &reported); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	set, cleared := reported.Settled()

	if len(set) != 2 {
		t.Errorf("kept %d values, want power and brightness only: %v", len(set), set)
	}

	if _, held := set[wire.TraitColor]; held {
		t.Error("color was kept as a value")
	}

	if len(cleared) != 1 || cleared[0] != wire.TraitColor {
		t.Errorf("cleared = %v, want [color]", cleared)
	}
}

func TestAStateWithNothingToRetireReportsNoClears(t *testing.T) {
	set, cleared := wire.State{wire.TraitPower: wire.Bool(false)}.Settled()

	if len(set) != 1 || len(cleared) != 0 {
		t.Errorf("set = %v, cleared = %v", set, cleared)
	}
}
