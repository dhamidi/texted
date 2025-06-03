package edlisp

import (
	"fmt"
)

// BuiltinLength returns the length of a string.
//
// This function takes a single string argument and returns the number of
// characters (bytes) in the string as a number. For empty strings, returns 0.
//
// Parameters:
//   - string: The string whose length to calculate
//
// Returns:
//   - number: The length of the string in characters
//
// Examples:
//   length "Hello world" → 11
//   length "" → 0
//   length "test" → 4
//
// Related functions:
//   - substring: Extract a portion of a string using indices
//   - concat: Combine multiple strings together
//
// Category: string
func BuiltinLength(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("length expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("length expects a string argument")
	}

	str := args[0].(*String)
	return NewNumber(float64(len(str.Value))), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "length",
		Category:    "string",
		Summary:     "Get the length of a string",
		Description: "Returns the number of characters (bytes) in STRING. For empty strings, returns 0.",
		Parameters: []ParameterDoc{
			{Name: "string", Type: "string", Description: "The string whose length to calculate"},
		},
		Examples: []ExampleDoc{
			{Description: "Get length of text", Input: `length "Hello world"`, Output: `11`},
			{Description: "Get length of empty string", Input: `length ""`, Output: `0`},
		},
		SeeAlso: []string{"substring", "concat"},
	})
}
