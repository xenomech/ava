package wire

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type Repr string

const (
	ReprNone   Repr = ""
	ReprBool   Repr = "bool"
	ReprNumber Repr = "number"
	ReprText   Repr = "text"
)

var ErrValueUnset = errors.New("wire: value is unset")

type Value struct {
	repr   Repr
	flag   bool
	number float64
	text   string
}

func Bool(v bool) Value {
	return Value{repr: ReprBool, flag: v}
}

func Number(v float64) Value {
	return Value{repr: ReprNumber, number: v}
}

func Text(v string) Value {
	return Value{repr: ReprText, text: v}
}

func (v Value) Repr() Repr {
	return v.repr
}

func (v Value) IsSet() bool {
	return v.repr != ReprNone
}

func (v Value) Bool() (flag, ok bool) {
	return v.flag, v.repr == ReprBool
}

func (v Value) Number() (number float64, ok bool) {
	return v.number, v.repr == ReprNumber
}

func (v Value) Text() (text string, ok bool) {
	return v.text, v.repr == ReprText
}

func (v Value) String() string {
	switch v.repr {
	case ReprBool:
		return fmt.Sprintf("%t", v.flag)
	case ReprNumber:
		return fmt.Sprintf("%g", v.number)
	case ReprText:
		return v.text
	case ReprNone:
		return "unset"
	default:
		return "unset"
	}
}

func (v Value) MarshalJSON() ([]byte, error) {
	switch v.repr {
	case ReprBool:
		return json.Marshal(v.flag)
	case ReprNumber:
		return json.Marshal(v.number)
	case ReprText:
		return json.Marshal(v.text)
	case ReprNone:
		return []byte("null"), nil
	default:
		return []byte("null"), nil
	}
}

func (v *Value) UnmarshalJSON(payload []byte) error {
	trimmed := bytes.TrimSpace(payload)

	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		*v = Value{}

		return nil
	}

	var flag bool
	if err := json.Unmarshal(trimmed, &flag); err == nil {
		*v = Bool(flag)

		return nil
	}

	var number float64
	if err := json.Unmarshal(trimmed, &number); err == nil {
		*v = Number(number)

		return nil
	}

	var text string
	if err := json.Unmarshal(trimmed, &text); err == nil {
		*v = Text(text)

		return nil
	}

	return fmt.Errorf("wire: value must be a boolean, number or string, got %s", trimmed)
}
