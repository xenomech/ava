package wire_test

import (
	"errors"
	"testing"

	"ava/pkg/wire"
)

func TestACommandCarriesAnyTraitAndValue(t *testing.T) {
	cases := []struct {
		name    string
		payload string
		trait   wire.Trait
	}{
		{"bool", `{"device_id":"a1","trait":"power","value":true}`, wire.TraitPower},
		{"number", `{"device_id":"a1","trait":"fan_speed","value":7}`, wire.TraitFanSpeed},
		{"enum", `{"device_id":"a1","trait":"mode","value":"grill"}`, wire.TraitMode},
		{"vendor", `{"device_id":"a1","trait":"oven:clean","value":true}`, "oven:clean"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd, err := wire.DecodeCommand([]byte(tc.payload))
			if err != nil {
				t.Fatalf("DecodeCommand: %v", err)
			}

			if cmd.Trait != tc.trait {
				t.Errorf("trait = %s, want %s", cmd.Trait, tc.trait)
			}
		})
	}
}

func TestACommandWithoutADeviceTraitOrValueIsRejected(t *testing.T) {
	cases := map[string]struct {
		payload string
		want    error
	}{
		"no device":  {`{"trait":"power","value":true}`, wire.ErrNoDevice},
		"no trait":   {`{"device_id":"a1","value":true}`, wire.ErrNoTrait},
		"no value":   {`{"device_id":"a1","trait":"power"}`, wire.ErrValueUnset},
		"null value": {`{"device_id":"a1","trait":"power","value":null}`, wire.ErrValueUnset},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := wire.DecodeCommand([]byte(tc.payload)); !errors.Is(err, tc.want) {
				t.Errorf("got %v, want %v", err, tc.want)
			}
		})
	}
}

func TestAStateEventCarriesOnlyTheTraitsTheDeviceHas(t *testing.T) {
	event, err := wire.DecodeStateEvent([]byte(`{"device_id":"a1","state":{"power":true,"brightness":60}}`))
	if err != nil {
		t.Fatalf("DecodeStateEvent: %v", err)
	}

	if brightness, ok := event.State.Get(wire.TraitBrightness); !ok {
		t.Error("brightness missing")
	} else if number, _ := brightness.Number(); number != 60 {
		t.Errorf("brightness = %g", number)
	}

	if _, ok := event.State.Get(wire.TraitColorTemp); ok {
		t.Error("color temp appeared from nowhere")
	}
}

func TestTopicsAreDerivedFromTheTenantAndHub(t *testing.T) {
	topics := wire.TopicsFor("acme", "hub-1")

	if topics.Command != "ava/acme/hub-1/cmd" {
		t.Errorf("command = %s", topics.Command)
	}

	if topics.State != "ava/acme/hub-1/state" {
		t.Errorf("state = %s", topics.State)
	}

	if topics.Status != "ava/acme/hub-1/status" {
		t.Errorf("status = %s", topics.Status)
	}
}

func TestAnApplyCarriesManyDevicesAtOnce(t *testing.T) {
	apply, err := wire.DecodeApply([]byte(`{"targets":[
		{"device_id":"a1","trait":"power","value":true},
		{"device_id":"b2","trait":"fan_speed","value":7},
		{"device_id":"c3","trait":"mode","value":"grill"}
	]}`))
	if err != nil {
		t.Fatalf("DecodeApply: %v", err)
	}

	if len(apply.Targets) != 3 {
		t.Fatalf("got %d targets, want 3", len(apply.Targets))
	}

	if apply.Targets[1].Trait != wire.TraitFanSpeed {
		t.Errorf("second target trait = %s", apply.Targets[1].Trait)
	}
}

func TestAnApplyWithNothingUsableIsRejected(t *testing.T) {
	cases := map[string]struct {
		payload string
		want    error
	}{
		"no targets":       {`{"targets":[]}`, wire.ErrNoTargets},
		"missing field":    {`{}`, wire.ErrNoTargets},
		"target no device": {`{"targets":[{"trait":"power","value":true}]}`, wire.ErrNoDevice},
		"target no trait":  {`{"targets":[{"device_id":"a1","value":true}]}`, wire.ErrNoTrait},
		"target no value":  {`{"targets":[{"device_id":"a1","trait":"power"}]}`, wire.ErrValueUnset},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := wire.DecodeApply([]byte(tc.payload)); !errors.Is(err, tc.want) {
				t.Errorf("got %v, want %v", err, tc.want)
			}
		})
	}
}

func TestTheApplyTopicSitsBesideCommand(t *testing.T) {
	topics := wire.TopicsFor("acme", "hub-1")

	if topics.Apply != "ava/acme/hub-1/apply" {
		t.Errorf("apply = %s", topics.Apply)
	}
}
