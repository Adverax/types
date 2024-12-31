package json

import (
	"encoding/json"
	"errors"
	"github.com/adverax/types"
)

type Number = json.Number

type RawMessage = json.RawMessage

// Boolean is a custom type for JSON booleans.
type Boolean bool

func (b *Boolean) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"true"`, `true`, `"1"`, `1`:
		*b = true
		return nil
	case `"false"`, `false`, `"0"`, `0`, `""`:
		*b = false
		return nil
	default:
		return errors.New("CustomBool: parsing \"" + string(data) + "\": unknown value")
	}
}

func (b *Boolean) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(*b)
	return data, err
}

// String is a custom type for JSON strings.
type String string

func (s *String) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `true`:
		*s = "true"
		return nil
	case `false`:
		*s = "false"
		return nil
	default:
		if val, ok := types.Type.String.TryCast(json.RawMessage(data)); ok {
			*s = String(val)
			return nil
		}
		return errors.New("CustomBool: parsing \"" + string(data) + "\": unknown value")
	}
}

func (s *String) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(*s)
	return data, err
}

// Logical is a custom type for JSON.
type Logical int

func (b *Logical) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"true"`, `true`, `"1"`, `1`:
		*b = 1
		return nil
	case `"false"`, `false`, `"0"`, `0`:
		*b = -1
		return nil
	default:
		*b = 0
		return nil
	}
}

func (b *Logical) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(*b)
	return data, err
}

// Strings is a custom type for JSON strings.
type Strings struct {
	value []string
}

func (s *Strings) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		if err := json.Unmarshal(data, &s.value); err != nil {
			return errors.New("Strings: parsing \"" + string(data) + "\": unknown value")
		}
		return nil
	}
	s.value = []string{value}
	return nil
}

func (s *Strings) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(s.value)
	return data, err
}

func (s *Strings) Values() []string {
	return s.value
}

func (s *Strings) First() string {
	if len(s.value) > 0 {
		return s.value[0]
	}
	return ""
}
