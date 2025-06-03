package edlisp

import (
	"fmt"
)

// BuiltinSubstring extracts a portion of a string.
//
// This function takes a string and one or two numeric indices to extract a
// substring. Indices are 1-based (the first character is at position 1).
// With two arguments, extracts from START to the end of the string.
// With three arguments, extracts from START to END (exclusive).
//
// The function performs bounds checking and adjusts invalid indices:
// - Negative start indices are adjusted to 0
// - End indices beyond string length are adjusted to string length
// - If start is greater than end, returns an empty string
//
// Parameters:
//   - string: The source string to extract from
//   - start: The starting position (1-based, inclusive)
//   - end: The ending position (1-based, exclusive) - optional
//
// Returns:
//   - string: The extracted substring
//
// Examples:
//   substring "Hello world" 1 5 → "Hell"
//   substring "Hello world" 7 → "world"
//   substring "test" 2 3 → "e"
//   substring "test" 10 → ""
//
// Related functions:
//   - length: Get the length of a string
//   - concat: Combine multiple strings together
//
// Category: string
func BuiltinSubstring(args []Value, buffer *Buffer) (Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("substring expects 2 or 3 arguments, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("substring expects a string as first argument")
	}

	if !IsA(args[1], TheNumberKind) {
		return nil, fmt.Errorf("substring expects a number as second argument")
	}

	str := args[0].(*String)
	start := int(args[1].(*Number).Value)
	end := len(str.Value)

	if len(args) == 3 {
		if !IsA(args[2], TheNumberKind) {
			return nil, fmt.Errorf("substring expects a number as third argument")
		}
		end = int(args[2].(*Number).Value)
	}

	// Convert from 1-based to 0-based indexing
	start--

	// For two-argument form, end should be to the end of string
	if len(args) == 2 {
		end = len(str.Value)
	} else {
		// For three-argument form, end is 1-based and exclusive
		// Convert to 0-based exclusive by decrementing
		end--
	}

	// Bounds checking
	if start < 0 {
		start = 0
	}
	if end > len(str.Value) {
		end = len(str.Value)
	}
	if start > end {
		return NewString(""), nil
	}

	result := str.Value[start:end]
	return NewString(result), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "substring",
		Category:    "string",
		Summary:     "Extract a portion of a string",
		Description: "Extracts a substring from STRING starting at START (1-based, inclusive). If END is provided, extracts up to END (1-based, exclusive). If END is omitted, extracts to the end of the string. Performs bounds checking and adjusts invalid indices safely.",
		Parameters: []ParameterDoc{
			{Name: "string", Type: "string", Description: "The source string to extract from"},
			{Name: "start", Type: "number", Description: "The starting position (1-based, inclusive)"},
			{Name: "end", Type: "number", Description: "The ending position (1-based, exclusive) - optional", Optional: true},
		},
		Examples: []ExampleDoc{
			{Description: "Extract substring with start and end", Input: `substring "Hello world" 1 5`, Output: `"Hell"`},
			{Description: "Extract substring to end", Input: `substring "Hello world" 7`, Output: `"world"`},
		},
		SeeAlso: []string{"length", "concat"},
	})
}
