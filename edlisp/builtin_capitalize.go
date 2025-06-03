package edlisp

import (
	"fmt"
	"strings"
)

// BuiltinCapitalize capitalizes the first character of a string.
//
// This function takes a single string argument and returns a new string with the
// first character converted to uppercase and all remaining characters converted
// to lowercase. If the string is empty, returns an empty string.
//
// Parameters:
//   - string: The string to capitalize
//
// Returns:
//   - string: A new string with the first character uppercase and remaining lowercase
//
// Examples:
//
//	capitalize "hello world" → "Hello world"
//	capitalize "HELLO WORLD" → "Hello world"
//	capitalize "test" → "Test"
//	capitalize "" → ""
//
// Related functions:
//   - upcase: Converts entire string to uppercase
//   - downcase: Converts entire string to lowercase
//
// Category: string
func BuiltinCapitalize(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("capitalize expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("capitalize expects a string argument")
	}

	str := args[0].(*String)
	if len(str.Value) == 0 {
		return NewString(""), nil
	}

	result := strings.ToUpper(string(str.Value[0])) + strings.ToLower(str.Value[1:])
	return NewString(result), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "capitalize",
		Category:    "string",
		Summary:     "Capitalize the first character of a string",
		Description: "Converts the first character of STRING to uppercase and all remaining characters to lowercase. Returns an empty string if STRING is empty.",
		Parameters: []ParameterDoc{
			{Name: "string", Type: "string", Description: "The string to capitalize"},
		},
		Examples: []ExampleDoc{
			{Description: "Capitalize lowercase text", Input: `capitalize "hello world"`, Output: `"Hello world"`},
			{Description: "Capitalize uppercase text", Input: `capitalize "HELLO WORLD"`, Output: `"Hello world"`},
		},
		SeeAlso: []string{"upcase", "downcase"},
	})
}
