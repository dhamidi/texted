package edlisp

import (
	"fmt"
	"strings"
)

// BuiltinConcat concatenates multiple strings into a single string.
//
// This function takes zero or more string arguments and returns a new string
// that is the concatenation of all input strings in the order they were provided.
// If no arguments are provided, returns an empty string.
//
// Parameters:
//   - strings...: Zero or more strings to concatenate
//
// Returns:
//   - string: A new string containing all input strings concatenated together
//
// Examples:
//   concat "Hello" " world" → "Hello world"
//   concat "Hello" " " "beautiful" " " "world" → "Hello beautiful world"
//   concat → ""
//   concat "single" → "single"
//
// Related functions:
//   - substring: Extract a portion of a string
//   - length: Get the length of a string
//
// Category: string
func BuiltinConcat(args []Value, buffer *Buffer) (Value, error) {
	var result strings.Builder
	
	for i, arg := range args {
		if !IsA(arg, TheStringKind) {
			return nil, fmt.Errorf("concat expects string arguments, got non-string at position %d", i+1)
		}
		str := arg.(*String)
		result.WriteString(str.Value)
	}
	
	return NewString(result.String()), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "concat",
		Category:    "string",
		Summary:     "Concatenate multiple strings",
		Description: "Concatenates zero or more strings into a single string. All arguments must be strings. Returns an empty string if no arguments are provided.",
		Parameters: []ParameterDoc{
			{Name: "strings", Type: "string", Description: "Zero or more strings to concatenate"},
		},
		Examples: []ExampleDoc{
			{Description: "Concatenate two strings", Input: `concat "Hello" " world"`, Output: `"Hello world"`},
			{Description: "Concatenate multiple strings", Input: `concat "Hello" " " "beautiful" " " "world"`, Output: `"Hello beautiful world"`},
		},
		SeeAlso: []string{"substring", "length"},
	})
}
