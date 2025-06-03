package edlisp

import (
	"fmt"
)

// BuiltinCurrentColumn returns the column number of the current point position.
//
// The column number represents the horizontal position of the point within the
// current line. It is 0-based, where column 0 is the first character of the line.
// The column is calculated by counting characters from the beginning of the current
// line (after the last newline character) to the point position.
//
// For multi-line buffers, this function finds the most recent newline character
// before the point and counts the characters from there. If there is no newline
// before the point, it counts from the beginning of the buffer.
//
// This function is useful for:
// - Determining horizontal cursor position
// - Aligning text or implementing indentation logic
// - Line-based text processing and formatting
//
// Returns:
//   - number: The column number (0-based) of the current point
//
// Examples:
//   current-column → 0 (at beginning of line)
//   current-column → 14 (at position 25 in "First line\nSecond line with content")
//   current-column → 5 (at position 6 in "Hello world")
//
// Related functions:
//   - line-number-at-pos: Get the line number of the current point
//   - point: Get current point position
//   - beginning-of-line: Move to beginning of current line
//   - end-of-line: Move to end of current line
//
// Category: position
func BuiltinCurrentColumn(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("current-column expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	if pos < 0 {
		return NewNumber(0), nil
	}
	if pos >= len(content) {
		pos = len(content) - 1
	}
	
	column := 0
	// Count backward to find beginning of line
	for i := pos; i >= 0 && content[i] != '\n'; i-- {
		column++
	}
	
	return NewNumber(float64(column)), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "current-column",
		Category:    "position",
		Summary:     "Return the column number of the current point position",
		Description: "Returns the horizontal position of the point within the current line. The column is 0-based, where column 0 is the first character of the line. Calculated by counting characters from the beginning of the current line (after the last newline) to the point position.",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{Description: "Get column at beginning of line", Input: `beginning-of-line; current-column`, Buffer: "Hello world", Output: "0"},
			{Description: "Get column in middle of multi-line", Input: `goto-char 25; current-column`, Buffer: "First line\nSecond line with content\nThird line", Output: "14"},
			{Description: "Get column in single line", Input: `goto-char 6; current-column`, Buffer: "Hello world", Output: "5"},
		},
		SeeAlso: []string{"line-number-at-pos", "point", "beginning-of-line", "end-of-line"},
	})
}
