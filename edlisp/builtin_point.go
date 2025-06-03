package edlisp

import (
	"fmt"
)

// BuiltinPoint returns the current position of the point (cursor) in the buffer.
//
// The point represents the current cursor position in the buffer. It is 1-based,
// meaning the first character in the buffer is at position 1. The point can be
// anywhere from point-min (1) to point-max (buffer-size + 1).
//
// When the point is at position N, it means the cursor is positioned just before
// the Nth character. When the point equals point-max, it means the cursor is
// positioned after the last character (at the end of the buffer).
//
// This function is fundamental for buffer navigation and is used by many other
// functions to determine the current editing position.
//
// Returns:
//   - number: The current point position (1-based)
//
// Examples:
//
//	point → 1 (at beginning of buffer)
//	point → 12 (somewhere in middle of buffer)
//	point → 17 (at end of 16-character buffer)
//
// Related functions:
//   - mark: Get the mark position
//   - goto-char: Move point to a specific position
//   - point-min: Get minimum valid point position
//   - point-max: Get maximum valid point position
//
// Category: position
func BuiltinPoint(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("point expects 0 arguments, got %d", len(args))
	}

	return NewNumber(float64(buffer.Point())), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "point",
		Category:    "position",
		Summary:     "Return the current position of the point (cursor) in the buffer",
		Description: "Returns the current cursor position in the buffer. The point is 1-based, where 1 is before the first character. When the point equals point-max, it is positioned after the last character. The point is fundamental for buffer navigation and editing operations.",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{Description: "Get point at beginning", Input: `point`, Buffer: "Hello world", Output: "1"},
			{Description: "Get point after moving", Input: `goto-char 5; point`, Buffer: "Hello world", Output: "5"},
			{Description: "Get point at end", Input: `goto-char 12; point`, Buffer: "Hello world", Output: "12"},
		},
		SeeAlso: []string{"mark", "goto-char", "point-min", "point-max"},
	})
}
