package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"ava/pkg/wire"
)

type KnownDevice struct {
	Vendor       string            `json:"vendor"`
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	IP           string            `json:"ip"`
	Kind         string            `json:"kind"`
	Capabilities wire.Capabilities `json:"capabilities"`
}

type State struct {
	HubID          string        `json:"hub_id"`
	HubName        string        `json:"hub_name"`
	TenantSlug     string        `json:"tenant_slug"`
	RefreshToken   string        `json:"refresh_token"`
	BrokerUsername string        `json:"broker_username,omitempty"`
	BrokerPassword string        `json:"broker_password,omitempty"`
	Devices        []KnownDevice `json:"devices,omitempty"`
}

func SameDevices(a, b []KnownDevice) bool {
	if len(a) != len(b) {
		return false
	}

	for at := range a {
		if a[at].ID != b[at].ID || a[at].IP != b[at].IP || a[at].Kind != b[at].Kind {
			return false
		}

		if len(a[at].Capabilities) != len(b[at].Capabilities) {
			return false
		}
	}

	return true
}

func Load(path string) (*State, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return &State{}, nil
		}

		return nil, fmt.Errorf("read state: %w", err)
	}

	var s State
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, fmt.Errorf("decode state: %w", err)
	}

	return &s, nil
}

func Save(path string, s *State) error {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return fmt.Errorf("create state dir: %w", err)
		}
	}

	raw, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("encode state: %w", err)
	}

	if err := os.WriteFile(path, raw, 0o600); err != nil {
		return fmt.Errorf("write state: %w", err)
	}

	return nil
}

func Clear(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("remove state: %w", err)
	}

	return nil
}

func (s *State) IsPaired() bool {
	return s.RefreshToken != ""
}
