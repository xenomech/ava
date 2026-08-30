package wire

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
)

type Trait string

const (
	TraitPower      Trait = "power"
	TraitBrightness Trait = "brightness"
	TraitColorTemp  Trait = "color_temp"
	TraitColor      Trait = "color"
	TraitFanSpeed   Trait = "fan_speed"
	TraitPosition   Trait = "position"
	TraitTargetTemp Trait = "target_temperature"
	TraitMode       Trait = "mode"

	TraitTemperature Trait = "temperature"
	TraitHumidity    Trait = "humidity"
	TraitBattery     Trait = "battery"
	TraitOccupancy   Trait = "occupancy"
	TraitPowerDraw   Trait = "power_draw"
)

type Kind string

const (
	KindBool   Kind = "bool"
	KindNumber Kind = "number"
	KindEnum   Kind = "enum"
	KindColor  Kind = "color"
)

type Access string

const (
	AccessRead      Access = "r"
	AccessReadWrite Access = "rw"
)

var (
	ErrUnknownTrait  = errors.New("wire: device does not have this trait")
	ErrReadOnly      = errors.New("wire: trait is read only")
	ErrWrongRepr     = errors.New("wire: value has the wrong type for this trait")
	ErrOutOfRange    = errors.New("wire: value is out of range")
	ErrOffStep       = errors.New("wire: value does not fall on a step")
	ErrNotAMember    = errors.New("wire: value is not one of the allowed values")
	ErrBadColor      = errors.New("wire: color must look like #rrggbb")
	ErrNoTrait       = errors.New("wire: trait is required")
	ErrBadCapability = errors.New("wire: capability is not well formed")
)

var hexColor = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

const stepEpsilon = 1e-9

type Capability struct {
	Trait  Trait    `json:"trait"`
	Kind   Kind     `json:"kind"`
	Access Access   `json:"access"`
	Min    *float64 `json:"min,omitempty"`
	Max    *float64 `json:"max,omitempty"`
	Step   *float64 `json:"step,omitempty"`
	Unit   string   `json:"unit,omitempty"`
	Values []string `json:"values,omitempty"`
}

func (c *Capability) Writable() bool {
	return c.Access == AccessReadWrite
}

func (c *Capability) Verify() error {
	if strings.TrimSpace(string(c.Trait)) == "" {
		return fmt.Errorf("%w: trait is empty", ErrBadCapability)
	}

	switch c.Kind {
	case KindBool, KindColor:
	case KindNumber:
		if c.Min != nil && c.Max != nil && *c.Min > *c.Max {
			return fmt.Errorf("%w: %s min is above max", ErrBadCapability, c.Trait)
		}

		if c.Step != nil && *c.Step <= 0 {
			return fmt.Errorf("%w: %s step must be positive", ErrBadCapability, c.Trait)
		}
	case KindEnum:
		if len(c.Values) == 0 {
			return fmt.Errorf("%w: %s declares no values", ErrBadCapability, c.Trait)
		}
	default:
		return fmt.Errorf("%w: %s has kind %q", ErrBadCapability, c.Trait, c.Kind)
	}

	if c.Access != AccessRead && c.Access != AccessReadWrite {
		return fmt.Errorf("%w: %s has access %q", ErrBadCapability, c.Trait, c.Access)
	}

	return nil
}

func (c *Capability) Validate(value Value) error {
	switch c.Kind {
	case KindBool:
		if _, ok := value.Bool(); !ok {
			return fmt.Errorf("%w: %s expects true or false", ErrWrongRepr, c.Trait)
		}

		return nil
	case KindNumber:
		number, ok := value.Number()
		if !ok {
			return fmt.Errorf("%w: %s expects a number", ErrWrongRepr, c.Trait)
		}

		return c.validateNumber(number)
	case KindEnum:
		text, ok := value.Text()
		if !ok {
			return fmt.Errorf("%w: %s expects one of %s", ErrWrongRepr, c.Trait, strings.Join(c.Values, ", "))
		}

		for _, allowed := range c.Values {
			if allowed == text {
				return nil
			}
		}

		return fmt.Errorf("%w: %s is not one of %s", ErrNotAMember, text, strings.Join(c.Values, ", "))
	case KindColor:
		text, ok := value.Text()
		if !ok || !hexColor.MatchString(text) {
			return fmt.Errorf("%w: %s", ErrBadColor, c.Trait)
		}

		return nil
	default:
		return fmt.Errorf("%w: %s has kind %q", ErrBadCapability, c.Trait, c.Kind)
	}
}

func (c *Capability) validateNumber(number float64) error {
	if c.Min != nil && number < *c.Min {
		return fmt.Errorf("%w: %s is below %g", ErrOutOfRange, c.Trait, *c.Min)
	}

	if c.Max != nil && number > *c.Max {
		return fmt.Errorf("%w: %s is above %g", ErrOutOfRange, c.Trait, *c.Max)
	}

	if c.Step == nil || *c.Step <= 0 {
		return nil
	}

	origin := 0.0
	if c.Min != nil {
		origin = *c.Min
	}

	offset := math.Abs(math.Mod(number-origin, *c.Step))
	if offset > stepEpsilon && math.Abs(offset-*c.Step) > stepEpsilon {
		return fmt.Errorf("%w: %s moves in steps of %g", ErrOffStep, c.Trait, *c.Step)
	}

	return nil
}

func (c *Capability) Clamp(value Value) Value {
	number, ok := value.Number()
	if !ok || c.Kind != KindNumber {
		return value
	}

	if c.Min != nil && number < *c.Min {
		number = *c.Min
	}

	if c.Max != nil && number > *c.Max {
		number = *c.Max
	}

	return Number(number)
}

type Capabilities []Capability

func (cs Capabilities) Find(trait Trait) (Capability, bool) {
	for at := range cs {
		if cs[at].Trait == trait {
			return cs[at], true
		}
	}

	return Capability{}, false
}

func (cs Capabilities) Has(trait Trait) bool {
	_, ok := cs.Find(trait)

	return ok
}

func (cs Capabilities) Verify() error {
	seen := make(map[Trait]struct{}, len(cs))

	for at := range cs {
		capability := &cs[at]

		if err := capability.Verify(); err != nil {
			return err
		}

		if _, duplicate := seen[capability.Trait]; duplicate {
			return fmt.Errorf("%w: %s is declared twice", ErrBadCapability, capability.Trait)
		}

		seen[capability.Trait] = struct{}{}
	}

	return nil
}

func (cs Capabilities) ValidateWrite(trait Trait, value Value) error {
	capability, ok := cs.Find(trait)
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownTrait, trait)
	}

	if !capability.Writable() {
		return fmt.Errorf("%w: %s", ErrReadOnly, trait)
	}

	return capability.Validate(value)
}

type State map[Trait]Value

func (s State) Get(trait Trait) (Value, bool) {
	value, ok := s[trait]

	return value, ok && value.IsSet()
}

// Settled splits traits the device reports a value for from ones it reports as no longer applicable.
func (s State) Settled() (set State, cleared []Trait) {
	set = make(State, len(s))

	for trait, value := range s {
		if value.IsSet() {
			set[trait] = value

			continue
		}

		cleared = append(cleared, trait)
	}

	return set, cleared
}

func (s State) With(trait Trait, value Value) State {
	if s == nil {
		s = State{}
	}

	s[trait] = value

	return s
}
