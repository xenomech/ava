package wire

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrNoDevice  = errors.New("wire: device_id is required")
	ErrNoTargets = errors.New("wire: at least one target is required")
)

type Topics struct {
	Command string
	Apply   string
	Status  string
	State   string
}

func TopicsFor(tenantSlug, hubID string) Topics {
	base := fmt.Sprintf("ava/%s/%s", tenantSlug, hubID)

	return Topics{
		Command: base + "/cmd",
		Apply:   base + "/apply",
		Status:  base + "/status",
		State:   base + "/state",
	}
}

type Presence struct {
	Online bool   `json:"online"`
	HubID  string `json:"hub_id"`
}

type Command struct {
	DeviceID string `json:"device_id"`
	Trait    Trait  `json:"trait"`
	Value    Value  `json:"value"`
}

func DecodeCommand(payload []byte) (Command, error) {
	var cmd Command

	if err := json.Unmarshal(payload, &cmd); err != nil {
		return Command{}, fmt.Errorf("wire: decode command: %w", err)
	}

	if cmd.DeviceID == "" {
		return Command{}, ErrNoDevice
	}

	if cmd.Trait == "" {
		return Command{}, ErrNoTrait
	}

	if !cmd.Value.IsSet() {
		return Command{}, fmt.Errorf("%w: %s", ErrValueUnset, cmd.Trait)
	}

	return cmd, nil
}

type ApplyTarget struct {
	DeviceID string `json:"device_id"`
	Trait    Trait  `json:"trait"`
	Value    Value  `json:"value"`
}

type Apply struct {
	Targets []ApplyTarget `json:"targets"`
}

func DecodeApply(payload []byte) (Apply, error) {
	var apply Apply

	if err := json.Unmarshal(payload, &apply); err != nil {
		return Apply{}, fmt.Errorf("wire: decode apply: %w", err)
	}

	if len(apply.Targets) == 0 {
		return Apply{}, ErrNoTargets
	}

	for at := range apply.Targets {
		target := &apply.Targets[at]

		if target.DeviceID == "" {
			return Apply{}, ErrNoDevice
		}

		if target.Trait == "" {
			return Apply{}, ErrNoTrait
		}

		if !target.Value.IsSet() {
			return Apply{}, fmt.Errorf("%w: %s", ErrValueUnset, target.Trait)
		}
	}

	return apply, nil
}

type StateEvent struct {
	DeviceID string `json:"device_id"`
	State    State  `json:"state"`
}

func DecodeStateEvent(payload []byte) (StateEvent, error) {
	var event StateEvent

	if err := json.Unmarshal(payload, &event); err != nil {
		return StateEvent{}, fmt.Errorf("wire: decode state event: %w", err)
	}

	if event.DeviceID == "" {
		return StateEvent{}, ErrNoDevice
	}

	return event, nil
}
