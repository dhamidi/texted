package edlisp

import (
	"fmt"
)

// BuiltinPointMin returns the minimum valid position for the point in the buffer.
//
// The point-min is always 1, representing the position just before the first
// character in the buffer. This is the earliest position the point can be
// moved to, regardless of buffer content or size.
//
// In texted's 1-based positioning system:
// - Position 1 (point-min) is before the first character
// - Position 2 is before the second character
// - Position N+1 (point-max) is after the last character in an N-character buffer
//
// This function is useful for:
// - Determining the start boundary of the buffer
// - Validating position arguments for other functions
// - Moving to the beginning of the buffer programmatically
//
// Returns:
//   - number: Always returns 1 (the minimum valid point position)
//
// Examples:
//   point-min → 1 (for any buffer, regardless of content)
//   point-min → 1 (for empty buffer)
//   point-min → 1 (for large buffer)
//
// Related functions:
//   - point-max: Get maximum valid point position
//   - buffer-size: Get total buffer size
//   - point: Get current point position
//   - beginning-of-buffer: Move point to beginning of buffer
//
// Category: position
func BuiltinPointMin(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("point-min expects 0 arguments, got %d", len(args))
	}
	
	return NewNumber(1), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "point-min",
		Category:    "position",
		Summary:     "Return the minimum valid position for the point in the buffer",
		Description: "Returns the minimum valid point position, which is always 1. This represents the position just before the first character in the buffer. The point-min is constant regardless of buffer content or size, providing the start boundary for all position-based operations.",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{Description: "Get point-min of any buffer", Input: `point-min`, Buffer: "Hello world", Output: "1"},
			{Description: "Get point-min of empty buffer", Input: `point-min`, Buffer: "", Output: "1"},
			{Description: "Move to beginning and verify", Input: `beginning-of-buffer; point; point-min`, Buffer: "Hello", Output: "1; 1"},
		},
		SeeAlso: []string{"point-max", "buffer-size", "point", "beginning-of-buffer"},
	})
}
