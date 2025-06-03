package edlisp

import (
	"fmt"
)

// BuiltinBufferSubstring extracts a portion of the buffer content between two positions.
//
// This function extracts text from the buffer between the specified START and END positions.
// Positions are 1-based, where 1 is the first character in the buffer. The extracted
// substring includes the character at START but excludes the character at END.
//
// Special handling:
// - If END is -1, it means extract to the end of the buffer
// - Positions are automatically adjusted to stay within buffer bounds
// - If START >= END after adjustment, returns an empty string
//
// This function is useful for extracting text regions, copying content, or analyzing
// specific portions of the buffer without modifying the buffer itself.
//
// Parameters:
//   - start: The starting position (1-based, inclusive)
//   - end: The ending position (1-based, exclusive), or -1 for end of buffer
//
// Returns:
//   - string: The extracted substring from the buffer
//
// Examples:
//
//	buffer-substring 1 6 → "Hello" (from buffer "Hello world")
//	buffer-substring 7 -1 → "world" (from buffer "Hello world")
//	buffer-substring 5 5 → "" (empty range)
//
// Related functions:
//   - buffer-size: Get total buffer size
//   - substring: Extract substring from a string value
//   - region-beginning: Get start of marked region
//   - region-end: Get end of marked region
//
// Category: buffer
func BuiltinBufferSubstring(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("buffer-substring expects 2 arguments, got %d", len(args))
	}

	if !IsA(args[0], TheNumberKind) || !IsA(args[1], TheNumberKind) {
		return nil, fmt.Errorf("buffer-substring expects number arguments")
	}

	start := int(args[0].(*Number).Value)
	end := int(args[1].(*Number).Value)

	content := buffer.String()

	// Handle special case: -1 means end of buffer
	if end == -1 {
		end = len(content) + 1
	}

	start-- // Convert to 0-based
	end--   // Convert to 0-based

	if start < 0 {
		start = 0
	}
	if end > len(content) {
		end = len(content)
	}
	if start >= end {
		return NewString(""), nil
	}

	return NewString(content[start:end]), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "buffer-substring",
		Category:    "buffer",
		Summary:     "Extract a portion of the buffer content between two positions",
		Description: "Extracts text from the buffer between the specified START and END positions. Positions are 1-based, where 1 is the first character. The extracted substring includes the character at START but excludes the character at END. If END is -1, extracts to the end of the buffer. Positions are automatically bounded to stay within the buffer.",
		Parameters: []ParameterDoc{
			{Name: "start", Type: "number", Description: "The starting position (1-based, inclusive)"},
			{Name: "end", Type: "number", Description: "The ending position (1-based, exclusive), or -1 for end of buffer"},
		},
		Examples: []ExampleDoc{
			{Description: "Extract first 5 characters", Input: `buffer-substring 1 6`, Buffer: "Hello world", Output: `"Hello"`},
			{Description: "Extract from position to end", Input: `buffer-substring 7 -1`, Buffer: "Hello world", Output: `"world"`},
			{Description: "Extract empty range", Input: `buffer-substring 5 5`, Buffer: "Hello world", Output: `""`},
		},
		SeeAlso: []string{"buffer-size", "substring", "region-beginning", "region-end"},
	})
}
