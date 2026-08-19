package wire

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	ActionPower      = "power"
	ActionBrightness = "brightness"
	ActionColorTemp  = "color_temp"
)

var Actions = []string{ActionPower, ActionBrightness, ActionColorTemp}

var (
	ErrUnknownAction = errors.New("wire: unknown action")
	ErrBadValue      = errors.New("wire: value does not match the action")
	ErrNoDevice      = errors.New("wire: device_id is required")
)

type Topics struct {
	Command string
	Status  string
	State   string
}

func TopicsFor(tenantSlug, hubID string) Topics {
	base := fmt.Sprintf("ava/%s/%s", tenantSlug, hubID)

	return Topics{
		Command: base + "/cmd",
		Status:  base + "/status",
		State:   base + "/state",
	}
}

type Command struct {
	DeviceID string          `json:"device_id"`
	Action   string          `json:"action"`
	Value    json.RawMessage `json:"value"`
}

func (c Command) Bool() (bool, error) {
	var on bool

	if err := json.Unmarshal(c.Value, &on); err != nil {
		return false, fmt.Errorf("%w: %s expects true or false", ErrBadValue, c.Action)
	}

	return on, nil
}

func (c Command) Int() (int, error) {
	var value int

	if err := json.Unmarshal(c.Value, &value); err != nil {
		return 0, fmt.Errorf("%w: %s expects a number", ErrBadValue, c.Action)
	}

	return value, nil
}

func DecodeCommand(payload []byte) (Command, error) {
	var cmd Command

	if err := json.Unmarshal(payload, &cmd); err != nil {
		return Command{}, fmt.Errorf("wire: decode command: %w", err)
	}

	if cmd.DeviceID == "" {
		return Command{}, ErrNoDevice
	}

	if !IsAction(cmd.Action) {
		return Command{}, fmt.Errorf("%w: %q", ErrUnknownAction, cmd.Action)
	}

	return cmd, nil
}

func IsAction(action string) bool {
	for _, known := range Actions {
		if known == action {
			return true
		}
	}

	return false
}

type Presence struct {
	Online bool   `json:"online"`
	HubID  string `json:"hub_id"`
}

type StateEvent struct {
	DeviceID   string `json:"device_id"`
	Power      bool   `json:"power"`
	Brightness int    `json:"brightness,omitempty"`
	ColorTemp  int    `json:"color_temp,omitempty"`
}
