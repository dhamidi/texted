package edlisp

import (
	"fmt"
)

// BuiltinPointMax returns the maximum valid position for the point in the buffer.
//
// The point-max represents the position just after the last character in the buffer.
// It is calculated as buffer-size + 1. This is the furthest position the point can
// be moved to, representing the end of the buffer.
//
// For example, in a buffer containing "Hello" (5 characters), point-max would be 6,
// meaning the point can be positioned after the last character 'o'.
//
// This function is useful for:
// - Determining buffer boundaries
// - Validating position arguments for other functions
// - Moving to the end of the buffer programmatically
//
// Returns:
//   - number: The maximum valid point position (buffer-size + 1)
//
// Examples:
//   point-max → 17 (for 16-character buffer "Hello world test")
//   point-max → 1 (for empty buffer)
//   point-max → 12 (for 11-character buffer "Hello world")
//
// Related functions:
//   - point-min: Get minimum valid point position
//   - buffer-size: Get total buffer size
//   - point: Get current point position
//   - end-of-buffer: Move point to end of buffer
//
// Category: position
func BuiltinPointMax(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("point-max expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	return NewNumber(float64(len(content) + 1)), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "point-max",
		Category:    "position",
		Summary:     "Return the maximum valid position for the point in the buffer",
		Description: "Returns the position just after the last character in the buffer (buffer-size + 1). This is the furthest position the point can be moved to, representing the end of the buffer. Useful for determining buffer boundaries and validating positions.",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{Description: "Get point-max of buffer", Input: `point-max`, Buffer: "Hello world test", Output: "17"},
			{Description: "Get point-max of empty buffer", Input: `point-max`, Buffer: "", Output: "1"},
			{Description: "Move to end and verify", Input: `end-of-buffer; point; point-max`, Buffer: "Hello", Output: "6; 6"},
		},
		SeeAlso: []string{"point-min", "buffer-size", "point", "end-of-buffer"},
	})
}
