package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type State struct {
	HubID        string `json:"device_id"`
	HubName      string `json:"device_name"`
	TenantSlug   string `json:"tenant_slug"`
	RefreshToken string `json:"refresh_token"`
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
