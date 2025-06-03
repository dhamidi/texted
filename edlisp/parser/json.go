// Package parser implements JSON parsing for texted scripts.
package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/dhamidi/texted/edlisp"
)

// ParseJSONReader parses texted scripts from a JSON-encoded io.Reader.
// Each JSON array represents a command list where:
// - The first element must be a string (the command/symbol)
// - Subsequent elements can be strings, numbers, or nested arrays
// - Arrays are converted to edlisp Lists
// - Numbers are converted to edlisp Numbers
// - Strings are converted to edlisp Strings (except the first element which becomes a Symbol)
func ParseJSONReader(r io.Reader) ([]edlisp.Value, error) {
	decoder := json.NewDecoder(r)
	var expressions []edlisp.Value

	for {
		var rawValue interface{}
		err := decoder.Decode(&rawValue)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("JSON decode error: %w", err)
		}

		expr, err := convertJSONValue(rawValue)
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, expr)
	}

	return expressions, nil
}

// ParseJSONString parses a JSON string containing texted script commands.
func ParseJSONString(s string) ([]edlisp.Value, error) {
	var rawValues []interface{}
	err := json.Unmarshal([]byte(s), &rawValues)
	if err != nil {
		return nil, fmt.Errorf("JSON unmarshal error: %w", err)
	}

	var expressions []edlisp.Value
	for i, rawValue := range rawValues {
		// Validate that each top-level value is an array
		if err := ValidateJSONFormat(rawValue); err != nil {
			return nil, fmt.Errorf("invalid format at index %d: %w", i, err)
		}

		expr, err := convertJSONValue(rawValue)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}

	return expressions, nil
}

// convertJSONValue converts a Go interface{} value from JSON into an edlisp.Value.
func convertJSONValue(value interface{}) (edlisp.Value, error) {
	switch v := value.(type) {
	case []interface{}:
		return convertJSONArray(v)
	case string:
		// Standalone strings become symbols (this shouldn't happen in valid texted JSON)
		return &edlisp.Symbol{Name: v}, nil
	case float64:
		// JSON numbers are always float64
		return &edlisp.Number{Value: v}, nil
	case bool, nil:
		return nil, fmt.Errorf("unsupported JSON type: %T", value)
	default:
		return nil, fmt.Errorf("unexpected JSON type: %T", value)
	}
}

// convertJSONArray converts a JSON array to an edlisp List.
// The first element must be a string (becomes a Symbol).
// Subsequent elements are converted according to their types.
func convertJSONArray(arr []interface{}) (edlisp.Value, error) {
	if len(arr) == 0 {
		return &edlisp.List{Elements: []edlisp.Value{}}, nil
	}

	var elements []edlisp.Value

	// First element must be a string (the command/symbol)
	firstElement, ok := arr[0].(string)
	if !ok {
		return nil, fmt.Errorf("first element of JSON array must be a string (symbol), got %T", arr[0])
	}
	elements = append(elements, &edlisp.Symbol{Name: firstElement})

	// Convert remaining elements
	for i, item := range arr[1:] {
		converted, err := convertJSONItem(item)
		if err != nil {
			return nil, fmt.Errorf("error converting array element %d: %w", i+1, err)
		}
		elements = append(elements, converted)
	}

	return &edlisp.List{Elements: elements}, nil
}

// convertJSONItem converts individual JSON items to edlisp values.
// This handles the conversion rules for non-first elements in arrays.
func convertJSONItem(item interface{}) (edlisp.Value, error) {
	switch v := item.(type) {
	case string:
		// Non-first strings become String values
		return &edlisp.String{Value: v}, nil
	case float64:
		// JSON numbers become Number values
		return &edlisp.Number{Value: v}, nil
	case []interface{}:
		// Nested arrays become Lists
		return convertJSONArray(v)
	case bool:
		return nil, fmt.Errorf("boolean values are not supported in texted JSON")
	case nil:
		return nil, fmt.Errorf("null values are not supported in texted JSON")
	default:
		return nil, fmt.Errorf("unsupported JSON type: %T (value: %v)", item, item)
	}
}

// ValidateJSONFormat checks if a JSON value conforms to texted JSON format rules.
func ValidateJSONFormat(value interface{}) error {
	return validateJSONValue(value, true)
}

// validateJSONValue recursively validates JSON structure.
func validateJSONValue(value interface{}, isTopLevel bool) error {
	switch v := value.(type) {
	case []interface{}:
		if len(v) == 0 {
			return nil // Empty arrays are allowed
		}

		// First element must be string
		if _, ok := v[0].(string); !ok {
			return fmt.Errorf("first element of array must be string, got %T", v[0])
		}

		// Validate remaining elements
		for i, item := range v[1:] {
			if err := validateJSONValue(item, false); err != nil {
				return fmt.Errorf("invalid element at index %d: %w", i+1, err)
			}
		}
		return nil

	case string:
		if isTopLevel {
			return fmt.Errorf("top-level strings are not allowed in texted JSON")
		}
		return nil

	case float64:
		if isTopLevel {
			return fmt.Errorf("top-level numbers are not allowed in texted JSON")
		}
		return nil

	case bool, nil:
		return fmt.Errorf("boolean and null values are not supported in texted JSON")

	default:
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Map {
			return fmt.Errorf("objects are not supported in texted JSON")
		}
		return fmt.Errorf("unsupported type: %T", value)
	}
}
