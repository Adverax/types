package json

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adverax/core"
	"github.com/adverax/types"
	"time"
)

// Update is a helper function to update a JSON document.
// It takes a document and a function that takes a Map and returns an error.
// Example: Update(doc, func(m Map) error { m["key"] = "value"; return nil })
func Update(
	doc []byte,
	actions ...func(Map) error,
) ([]byte, error) {
	m, err := NewMap(doc)
	if err != nil {
		return nil, fmt.Errorf("NewMapFromJson: %w", err)
	}

	for _, action := range actions {
		err = action(m)
		if err != nil {
			return nil, fmt.Errorf("action: %w", err)
		}
	}

	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("Marshal: %w", err)
	}
	return data, nil
}

// Extract is a helper function to extract a value from a JSON document.
func Extract(
	doc []byte,
	path string,
) (RawMessage, error) {
	m, err := NewMap(doc)
	if err != nil {
		return nil, fmt.Errorf("NewMap: %w", err)
	}

	return m.GetJson(context.Background(), path, Empty)
}

// UpdateAll is a helper function to update a list of JSON documents.
func UpdateAll(
	doc []byte,
	actions ...func(Map) error,
) ([]byte, error) {
	var rows []RawMessage
	err := Unmarshal(doc, &rows)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %w", err)
	}

	for i, row := range rows {
		rows[i], err = Update(row, actions...)
		if err != nil {
			return nil, fmt.Errorf("Update: %w", err)
		}
	}

	return Marshal(rows)
}

// Get is a helper function to extract a value from a JSON document.
// Example: Get(doc, GetInteger("key", 0))
func Get(
	doc []byte,
	getters ...func(Map) error,
) (err error) {
	m, err := NewMap(doc)
	if err != nil {
		return fmt.Errorf("NewMap: %w", err)
	}

	for _, getter := range getters {
		err := getter(m)
		if err != nil {
			return fmt.Errorf("getter: %w", err)
		}
	}

	return nil
}

// WithValue is a helper function to override a value in a JSON document.
// See Update for more information.
func WithValue(key string, val interface{}) func(Map) error {
	return func(doc Map) error {
		return doc.SetProperty(context.Background(), key, val)
	}
}

// WithValues is a helper function to update a JSON document with a map of values.
// It takes a document and a map of fields to update.
// Example: Update(doc, Values(map[string]interface{}{"key": "value"}))
func WithValues(
	fields map[string]interface{},
) func(Map) error {
	return func(doc Map) error {
		doc.ExpandBy(NewMapFromStruct(fields))
		return nil
	}
}

// WithDefaultValue is a helper function to set a default value in a JSON document.
// See Update for more information.
func WithDefaultValue(key string, val interface{}) func(Map) error {
	return func(doc Map) error {
		value, err := doc.GetProperty(context.Background(), key)
		if err != nil {
			if errors.Is(err, types.GetErrNoMatch()) {
				doc[key] = val
				return nil
			}
			return fmt.Errorf("GetProperty: %w", err)
		}

		if core.IsZeroValue(value) {
			return doc.SetProperty(context.Background(), key, val)
		}
		return nil
	}
}

// WithRemove is a helper function to remove a value from a JSON document.
func WithRemove(key string) func(Map) error {
	return func(doc Map) error {
		delete(doc, key)
		return nil
	}
}

// WithIf is a helper function to conditionally execute an action.
// See Update for more information.
func WithIf(cond bool, action func(Map) error) func(Map) error {
	if cond {
		return action
	}
	return dummyAction
}

// WithInteger is a helper function to extract an integer from a JSON document.
// See Get for more information.
func WithInteger(key string, ref *int64) func(Map) error {
	return func(doc Map) (err error) {
		*ref, err = doc.GetInteger(context.Background(), key, *ref)
		return err
	}
}

// WithFloat is a helper function to extract a float from a JSON document.
// See Get for more information.
func WithFloat(key string, ref *float64) func(Map) error {
	return func(doc Map) (err error) {
		*ref, err = doc.GetFloat(context.Background(), key, *ref)
		return err
	}
}

// WithString is a helper function to extract a string from a JSON document.
// See Get for more information.
func WithString(key string, ref *string) func(Map) error {
	return func(doc Map) (err error) {
		*ref, err = doc.GetString(context.Background(), key, *ref)
		return err
	}
}

// WithBoolean is a helper function to extract a boolean from a JSON document.
// See Get for more information.
func WithBoolean(key string, ref *bool) func(Map) error {
	return func(doc Map) (err error) {
		*ref, err = doc.GetBoolean(context.Background(), key, *ref)
		return err
	}
}

// WithDuration is a helper function to extract a time from a JSON document.
// See Get for more information.
func WithDuration(key string, ref *time.Duration) func(Map) error {
	return func(doc Map) (err error) {
		*ref, err = doc.GetDuration(context.Background(), key, *ref)
		return err
	}
}

// WihtAny is a helper function to extract a any value from a JSON document.
// See Get for more information.
func WithAny(key string, ref *any) func(Map) error {
	return func(doc Map) (err error) {
		*ref, err = doc.GetProperty(context.Background(), key)
		return err
	}
}

// WithJson is a helper function to extract a raw message from a JSON document.
// See Get for more information.
func WithJson(key string, ref *RawMessage) func(Map) error {
	return func(doc Map) (err error) {
		val, err := doc.GetProperty(context.Background(), key)
		if err != nil {
			return err
		}
		*ref, err = Marshal(val)
		return err
	}
}
