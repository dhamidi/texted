package parser

import (
	"strings"
	"testing"

	"github.com/dhamidi/texted/edlisp"
)

func TestParseJSONString_SingleCommand(t *testing.T) {
	input := `[["search-forward", "doIt"]]`
	result, err := ParseJSONString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 expression, got %d", len(result))
	}

	list, ok := result[0].(*edlisp.List)
	if !ok {
		t.Fatalf("expected List, got %T", result[0])
	}

	if len(list.Elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(list.Elements))
	}

	sym, ok := list.Elements[0].(*edlisp.Symbol)
	if !ok {
		t.Fatalf("expected Symbol, got %T", list.Elements[0])
	}
	if sym.Name != "search-forward" {
		t.Errorf("expected symbol 'search-forward', got '%s'", sym.Name)
	}

	str, ok := list.Elements[1].(*edlisp.String)
	if !ok {
		t.Fatalf("expected String, got %T", list.Elements[1])
	}
	if str.Value != "doIt" {
		t.Errorf("expected string 'doIt', got '%s'", str.Value)
	}
}

func TestParseJSONString_MultipleCommands(t *testing.T) {
	input := `[
		["search-forward", "doIt"],
		["set-mark"],
		["search-forward", "("],
		["replace-region", "helloWorld"]
	]`

	result, err := ParseJSONString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 4 {
		t.Fatalf("expected 4 expressions, got %d", len(result))
	}

	// Check first command
	list := result[0].(*edlisp.List)
	if list.Elements[0].(*edlisp.Symbol).Name != "search-forward" {
		t.Errorf("expected 'search-forward', got '%s'", list.Elements[0].(*edlisp.Symbol).Name)
	}
	if list.Elements[1].(*edlisp.String).Value != "doIt" {
		t.Errorf("expected 'doIt', got '%s'", list.Elements[1].(*edlisp.String).Value)
	}

	// Check second command (no arguments)
	list = result[1].(*edlisp.List)
	if len(list.Elements) != 1 {
		t.Errorf("expected 1 element for set-mark, got %d", len(list.Elements))
	}
	if list.Elements[0].(*edlisp.Symbol).Name != "set-mark" {
		t.Errorf("expected 'set-mark', got '%s'", list.Elements[0].(*edlisp.Symbol).Name)
	}
}

func TestParseJSONString_WithNumbers(t *testing.T) {
	input := `[["move-point", 42], ["set-value", 3.14]]`
	result, err := ParseJSONString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 expressions, got %d", len(result))
	}

	// Check first command with integer
	list := result[0].(*edlisp.List)
	num, ok := list.Elements[1].(*edlisp.Number)
	if !ok {
		t.Fatalf("expected Number, got %T", list.Elements[1])
	}
	if num.Value != 42 {
		t.Errorf("expected 42, got %f", num.Value)
	}

	// Check second command with float
	list = result[1].(*edlisp.List)
	num, ok = list.Elements[1].(*edlisp.Number)
	if !ok {
		t.Fatalf("expected Number, got %T", list.Elements[1])
	}
	if num.Value != 3.14 {
		t.Errorf("expected 3.14, got %f", num.Value)
	}
}

func TestParseJSONString_NestedArrays(t *testing.T) {
	input := `[["progn", ["search-forward", "test"], ["set-mark"]]]`
	result, err := ParseJSONString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 expression, got %d", len(result))
	}

	list := result[0].(*edlisp.List)
	if len(list.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(list.Elements))
	}

	// Check first element is symbol 'progn'
	sym := list.Elements[0].(*edlisp.Symbol)
	if sym.Name != "progn" {
		t.Errorf("expected 'progn', got '%s'", sym.Name)
	}

	// Check second element is nested list
	nestedList := list.Elements[1].(*edlisp.List)
	if nestedList.Elements[0].(*edlisp.Symbol).Name != "search-forward" {
		t.Errorf("expected 'search-forward' in nested list")
	}
	if nestedList.Elements[1].(*edlisp.String).Value != "test" {
		t.Errorf("expected 'test' in nested list")
	}

	// Check third element is another nested list
	nestedList2 := list.Elements[2].(*edlisp.List)
	if nestedList2.Elements[0].(*edlisp.Symbol).Name != "set-mark" {
		t.Errorf("expected 'set-mark' in second nested list")
	}
}

func TestParseJSONString_EmptyArray(t *testing.T) {
	input := `[[]]`
	result, err := ParseJSONString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 expression, got %d", len(result))
	}

	list := result[0].(*edlisp.List)
	if len(list.Elements) != 0 {
		t.Errorf("expected empty list, got %d elements", len(list.Elements))
	}
}

func TestParseJSONReader(t *testing.T) {
	input := `["search-forward", "test"]
["set-mark"]
["replace-match", "hello"]`

	reader := strings.NewReader(input)
	result, err := ParseJSONReader(reader)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 expressions, got %d", len(result))
	}

	// Check first command
	list := result[0].(*edlisp.List)
	if list.Elements[0].(*edlisp.Symbol).Name != "search-forward" {
		t.Errorf("expected 'search-forward'")
	}
}

func TestParseJSONString_ErrorCases(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"non-string first element", `[[42, "test"]]`},
		{"top-level string", `["not-array"]`},
		{"top-level number", `[42]`},
		{"boolean value", `[["command", true]]`},
		{"null value", `[["command", null]]`},
		{"object value", `[["command", {"key": "value"}]]`},
		{"invalid JSON", `[["command", "test"`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseJSONString(tc.input)
			if err == nil {
				t.Errorf("expected error for input: %q", tc.input)
			}
		})
	}
}

func TestConvertJSONValue(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected interface{}
		hasError bool
	}{
		{
			name:     "string array",
			input:    []interface{}{"search-forward", "test"},
			expected: "List with Symbol and String",
			hasError: false,
		},
		{
			name:     "number array",
			input:    []interface{}{"move", 42.0},
			expected: "List with Symbol and Number",
			hasError: false,
		},
		{
			name:     "empty array",
			input:    []interface{}{},
			expected: "Empty List",
			hasError: false,
		},
		{
			name:     "standalone string",
			input:    "symbol",
			expected: "Symbol",
			hasError: false,
		},
		{
			name:     "standalone number",
			input:    42.0,
			expected: "Number",
			hasError: false,
		},
		{
			name:     "boolean",
			input:    true,
			expected: nil,
			hasError: true,
		},
		{
			name:     "null",
			input:    nil,
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := convertJSONValue(tc.input)
			if tc.hasError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected result but got nil")
				}
			}
		})
	}
}

func TestConvertJSONArray(t *testing.T) {
	testCases := []struct {
		name     string
		input    []interface{}
		hasError bool
	}{
		{
			name:     "valid command",
			input:    []interface{}{"search-forward", "test"},
			hasError: false,
		},
		{
			name:     "command with number",
			input:    []interface{}{"move", 42.0},
			hasError: false,
		},
		{
			name:     "empty array",
			input:    []interface{}{},
			hasError: false,
		},
		{
			name:     "non-string first element",
			input:    []interface{}{42, "test"},
			hasError: true,
		},
		{
			name:     "boolean argument",
			input:    []interface{}{"command", true},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := convertJSONArray(tc.input)
			if tc.hasError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected result but got nil")
				}
			}
		})
	}
}

func TestConvertJSONItem(t *testing.T) {
	testCases := []struct {
		name         string
		input        interface{}
		expectedType string
		hasError     bool
	}{
		{
			name:         "string",
			input:        "hello",
			expectedType: "*edlisp.String",
			hasError:     false,
		},
		{
			name:         "number",
			input:        42.0,
			expectedType: "*edlisp.Number",
			hasError:     false,
		},
		{
			name:         "nested array",
			input:        []interface{}{"command", "arg"},
			expectedType: "*edlisp.List",
			hasError:     false,
		},
		{
			name:         "boolean",
			input:        true,
			expectedType: "",
			hasError:     true,
		},
		{
			name:         "null",
			input:        nil,
			expectedType: "",
			hasError:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := convertJSONItem(tc.input)
			if tc.hasError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected result but got nil")
				}

				actualType := ""
				switch result.(type) {
				case *edlisp.String:
					actualType = "*edlisp.String"
				case *edlisp.Number:
					actualType = "*edlisp.Number"
				case *edlisp.List:
					actualType = "*edlisp.List"
				case *edlisp.Symbol:
					actualType = "*edlisp.Symbol"
				}

				if actualType != tc.expectedType {
					t.Errorf("expected type %s, got %s", tc.expectedType, actualType)
				}
			}
		})
	}
}

func TestValidateJSONFormat(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		hasError bool
	}{
		{
			name:     "valid array",
			input:    []interface{}{"command", "arg"},
			hasError: false,
		},
		{
			name:     "empty array",
			input:    []interface{}{},
			hasError: false,
		},
		{
			name:     "nested valid arrays",
			input:    []interface{}{"progn", []interface{}{"search", "text"}},
			hasError: false,
		},
		{
			name:     "top-level string",
			input:    "not-allowed",
			hasError: true,
		},
		{
			name:     "top-level number",
			input:    42.0,
			hasError: true,
		},
		{
			name:     "non-string first element",
			input:    []interface{}{42, "arg"},
			hasError: true,
		},
		{
			name:     "boolean in array",
			input:    []interface{}{"command", true},
			hasError: true,
		},
		{
			name:     "null in array",
			input:    []interface{}{"command", nil},
			hasError: true,
		},
		{
			name:     "object",
			input:    map[string]interface{}{"key": "value"},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateJSONFormat(tc.input)
			if tc.hasError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestParseJSONReader_EmptyInput(t *testing.T) {
	reader := strings.NewReader("")
	result, err := ParseJSONReader(reader)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d expressions", len(result))
	}
}

func TestParseJSONReader_InvalidJSON(t *testing.T) {
	reader := strings.NewReader(`["command", "arg"`)
	_, err := ParseJSONReader(reader)
	if err == nil {
		t.Errorf("expected error for invalid JSON")
	}
}
