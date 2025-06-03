package edlisp

import (
	"fmt"
	"strings"
)

// BuiltinDowncase converts a string to lowercase.
//
// This function takes a single string argument and returns a new string with all
// alphabetic characters converted to their lowercase equivalents. Non-alphabetic
// characters remain unchanged.
//
// Parameters:
//   - string: The string to convert to lowercase
//
// Returns:
//   - string: A new string with all alphabetic characters in lowercase
//
// Examples:
//
//	downcase "Hello World" → "hello world"
//	downcase "TEST123" → "test123"
//	downcase "" → ""
//
// Related functions:
//   - upcase: Converts string to uppercase
//   - capitalize: Capitalizes the first letter of a string
//
// Category: string
func BuiltinDowncase(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("downcase expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("downcase expects a string argument")
	}

	str := args[0].(*String)
	return NewString(strings.ToLower(str.Value)), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "downcase",
		Category:    "string",
		Summary:     "Convert string to lowercase",
		Description: "Converts all alphabetic characters in STRING to lowercase. Non-alphabetic characters remain unchanged.",
		Parameters: []ParameterDoc{
			{Name: "string", Type: "string", Description: "The string to convert to lowercase"},
		},
		Examples: []ExampleDoc{
			{Description: "Convert mixed case to lowercase", Input: `downcase "Hello World"`, Output: `"hello world"`},
			{Description: "Convert with numbers", Input: `downcase "TEST123"`, Output: `"test123"`},
		},
		SeeAlso: []string{"upcase", "capitalize"},
	})
}
