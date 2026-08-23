package device

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ava/pkg/wire"
)

type Info struct {
	ID       string    `json:"id"`
	Vendor   string    `json:"vendor"`
	Name     string    `json:"name"`
	Model    string    `json:"model,omitempty"`
	IP       string    `json:"ip,omitempty"`
	MAC      string    `json:"mac,omitempty"`
	Parent   string    `json:"parent,omitempty"`
	LastSeen time.Time `json:"last_seen,omitempty"`
}

type Device interface {
	Info() Info
	Capabilities() wire.Capabilities
	State(ctx context.Context) (wire.State, error)
	Apply(ctx context.Context, trait wire.Trait, value wire.Value) error
}

var ErrUnsupported = errors.New("device: trait not supported")

func Unsupported(vendor string, trait wire.Trait) error {
	return fmt.Errorf("%s: %s: %w", vendor, trait, ErrUnsupported)
}

func Bounded(trait wire.Trait, lowest, highest float64, unit string) wire.Capability {
	low, high := lowest, highest

	return wire.Capability{
		Trait:  trait,
		Kind:   wire.KindNumber,
		Access: wire.AccessReadWrite,
		Min:    &low,
		Max:    &high,
		Unit:   unit,
	}
}

func Switch(trait wire.Trait) wire.Capability {
	return wire.Capability{Trait: trait, Kind: wire.KindBool, Access: wire.AccessReadWrite}
}

func Reading(trait wire.Trait, unit string) wire.Capability {
	return wire.Capability{Trait: trait, Kind: wire.KindNumber, Access: wire.AccessRead, Unit: unit}
}
