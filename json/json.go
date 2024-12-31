package json

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
)

// Config is a custom JSON configuration.
var Config = jsoniter.Config{
	EscapeHTML:             false,
	SortMapKeys:            false,
	ValidateJsonRawMessage: false,
	UseNumber:              true,
}.Froze()

// ConfigSorted is a custom JSON configuration.
var ConfigSorted = jsoniter.Config{
	EscapeHTML:             false,
	SortMapKeys:            true,
	ValidateJsonRawMessage: false,
	UseNumber:              true,
}.Froze()

// Unmarshal is a helper function to unmarshal a JSON document.
func Unmarshal(raw []byte, value interface{}) error {
	err := Config.Unmarshal(raw, &value)
	if err != nil && raw != nil {
		return fmt.Errorf("Config.Unmarshal: %w", err)
	}
	return nil
}

// Marshal is a helper function to marshal a JSON document.
func Marshal(value interface{}) ([]byte, error) {
	return Config.Marshal(value)
}

// MarshalIndent is a helper function to marshal a JSON document with indentation.
func MarshalIndent(value interface{}) ([]byte, error) {
	return json.MarshalIndent(value, "", "  ")
}

// Ensure is a helper function to ensure a JSON document is not empty.
// It takes a document and returns an empty object if the document is empty.
func Ensure(doc []byte) []byte {
	if IsEmpty(doc) {
		return Empty
	}
	return doc
}

// TypeOf is a helper function to determine the type of a JSON document.
// It takes a document and returns the type (array, object).
// Example: TypeOf(doc)
func TypeOf(in io.Reader) (string, error) {
	dec := json.NewDecoder(in)
	// Get just the first valid JSON token from input
	t, err := dec.Token()
	if err != nil {
		return "", err
	}
	if d, ok := t.(json.Delim); ok {
		// The first token is a delimiter, so this is an array or an object
		switch d {
		case '[':
			return "array", nil
		case '{':
			return "object", nil
		default: // ] or }
			return "", errors.New("Unexpected delimiter")
		}
	}
	return "", errors.New("Input does not represent a JSON object or array")
}

// AsArray is a helper function to ensure a JSON document is an array.
func AsArray(data []byte) ([]byte, error) {
	isArray, err := IsArray(data)
	if err != nil {
		return nil, fmt.Errorf("JsonIsArray: %w", err)
	}
	if isArray {
		return data, nil
	}

	return []byte(fmt.Sprintf("[%s]", string(data))), nil
}

// IsArray is a helper function to determine if a JSON document is an array.
func IsArray(data []byte) (bool, error) {
	if len(data) == 0 {
		return false, nil
	}
	tp, err := TypeOf(bytes.NewBuffer(data))
	if err != nil {
		return false, fmt.Errorf("JsonType: %w", err)
	}
	return tp == "array", nil
}

// IsObject is a helper function to determine if a JSON document is an object.
func IsObject(data []byte) (bool, error) {
	if len(data) == 0 {
		return false, nil
	}
	tp, err := TypeOf(bytes.NewBuffer(data))
	if err != nil {
		return false, fmt.Errorf("JsonType: %w", err)
	}
	return tp == "object", nil
}

// Merge is a helper function to merge JSON documents.
func Merge(docs ...[]byte) ([]byte, error) {
	res := Map{}
	for i, doc := range docs {
		m, err := NewMap(doc)
		if err != nil {
			return nil, fmt.Errorf("NewMapFromJson[%d]: %w", i, err)
		}
		if i == 0 {
			res = m
			continue
		}
		res.ExpandBy(m)
	}
	return res.Json()
}

// Empty is empty document
var Empty = RawMessage("{}")

func IsEmpty(doc []byte) bool {
	return len(doc) <= 2
}

// CoalesceString is a helper function to coalesce a string from multiple documents.
func CoalesceString(ctx context.Context, path string, defVal string, docs ...Map) (string, error) {
	for _, doc := range docs {
		if doc == nil {
			continue
		}
		val, err := doc.GetString(ctx, path, "")
		if err != nil {
			return val, fmt.Errorf("GetString: %w", err)
		}
		if val != "" {
			return val, nil
		}
	}

	return defVal, nil
}

func RestoreString(ctx context.Context, path string, docs ...Map) error {
	val, err := CoalesceString(ctx, path, "", docs...)
	if err != nil {
		return fmt.Errorf("JsonCoalesceString: %w", err)
	}

	if val != "" {
		err = docs[0].SetString(ctx, path, val)
		if err != nil {
			return fmt.Errorf("SetString: %w", err)
		}
	}

	return nil
}

func Normalize(data RawMessage) (RawMessage, error) {
	m, err := NewMap(data)
	if err != nil {
		return nil, err
	}

	return MarshalIndent(m)
}

func IsEqual(a, b RawMessage) bool {
	aa, err := Normalize(a)
	if err != nil {
		return false
	}

	bb, err := Normalize(b)
	if err != nil {
		return false
	}

	return string(aa) == string(bb)
}

var dummyAction = func(Map) error { return nil }
