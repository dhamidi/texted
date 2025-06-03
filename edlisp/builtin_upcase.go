package edlisp

import (
	"fmt"
	"strings"
)

// BuiltinUpcase converts a string to uppercase.
//
// This function takes a single string argument and returns a new string with all
// alphabetic characters converted to their uppercase equivalents. Non-alphabetic
// characters remain unchanged.
//
// Parameters:
//   - string: The string to convert to uppercase
//
// Returns:
//   - string: A new string with all alphabetic characters in uppercase
//
// Examples:
//
//	upcase "Hello World" → "HELLO WORLD"
//	upcase "test123" → "TEST123"
//	upcase "" → ""
//
// Related functions:
//   - downcase: Converts string to lowercase
//   - capitalize: Capitalizes the first letter of a string
//
// Category: string
func BuiltinUpcase(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("upcase expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("upcase expects a string argument")
	}

	str := args[0].(*String)
	return NewString(strings.ToUpper(str.Value)), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "upcase",
		Category:    "string",
		Summary:     "Convert string to uppercase",
		Description: "Converts all alphabetic characters in STRING to uppercase. Non-alphabetic characters remain unchanged.",
		Parameters: []ParameterDoc{
			{Name: "string", Type: "string", Description: "The string to convert to uppercase"},
		},
		Examples: []ExampleDoc{
			{Description: "Convert mixed case to uppercase", Input: `upcase "Hello World"`, Output: `"HELLO WORLD"`},
			{Description: "Convert with numbers", Input: `upcase "test123"`, Output: `"TEST123"`},
		},
		SeeAlso: []string{"downcase", "capitalize"},
	})
}
