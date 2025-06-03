package edlisp

import (
	"fmt"
)

// BuiltinMark returns the current position of the mark in the buffer.
//
// The mark is a secondary position in the buffer that works together with the point
// to define a region (selection). Like the point, the mark is 1-based and can be
// anywhere from point-min (1) to point-max (buffer-size + 1).
//
// The mark is typically set using set-mark or set-mark-command, and together with
// the point, it defines a text region that can be operated on by functions like
// delete-region, replace-region, or copied with buffer-substring.
//
// The region spans from the smaller of point and mark to the larger of point and mark.
// This allows for flexible text selection regardless of the direction of marking.
//
// Returns:
//   - number: The current mark position (1-based)
//
// Examples:
//
//	mark → 1 (mark at beginning of buffer)
//	mark → 5 (mark at position 5)
//	mark → 17 (mark at end of 16-character buffer)
//
// Related functions:
//   - point: Get the point position
//   - set-mark: Set mark at a specific position
//   - set-mark-command: Set mark at current point
//   - region-beginning: Get start of marked region
//   - region-end: Get end of marked region
//   - exchange-point-and-mark: Swap point and mark positions
//
// Category: position
func BuiltinMark(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("mark expects 0 arguments, got %d", len(args))
	}

	return NewNumber(float64(buffer.Mark())), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "mark",
		Category:    "position",
		Summary:     "Return the current position of the mark in the buffer",
		Description: "Returns the current mark position in the buffer. The mark is a secondary position that works with the point to define text regions. Like the point, it is 1-based and can be anywhere from point-min to point-max. The mark is typically set using set-mark or set-mark-command.",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{Description: "Get mark position", Input: `set-mark 5; mark`, Buffer: "Hello world", Output: "5"},
			{Description: "Get mark after set-mark-command", Input: `goto-char 3; set-mark-command; mark`, Buffer: "Hello world", Output: "3"},
		},
		SeeAlso: []string{"point", "set-mark", "set-mark-command", "region-beginning", "region-end", "exchange-point-and-mark"},
	})
}
