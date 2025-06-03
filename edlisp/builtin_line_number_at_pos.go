package edlisp

import (
	"fmt"
)

// BuiltinLineNumberAtPos returns the line number of the current point position.
//
// The line number represents the vertical position of the point within the buffer.
// It is 1-based, where line 1 is the first line of the buffer. Lines are separated
// by newline characters (\n).
//
// The function counts the number of newline characters from the beginning of the
// buffer up to (but not including) the current point position, then adds 1 to get
// the line number. This means:
// - Characters before the first \n are on line 1
// - Characters between the first and second \n are on line 2
// - And so on...
//
// This function is useful for:
// - Determining vertical cursor position in multi-line text
// - Implementing line-based navigation and editing
// - Error reporting and debugging with line references
// - Text processing that needs line-aware operations
//
// Returns:
//   - number: The line number (1-based) of the current point
//
// Examples:
//
//	line-number-at-pos → 1 (at beginning of buffer or first line)
//	line-number-at-pos → 2 (on second line of multi-line buffer)
//	line-number-at-pos → 3 (on third line after two newlines)
//
// Related functions:
//   - current-column: Get the column number of the current point
//   - point: Get current point position
//   - goto-line: Move to a specific line number
//   - beginning-of-line: Move to beginning of current line
//   - end-of-line: Move to end of current line
//
// Category: position
func BuiltinLineNumberAtPos(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("line-number-at-pos expects 0 arguments, got %d", len(args))
	}

	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based

	if pos < 0 {
		pos = 0
	}
	if pos > len(content) {
		pos = len(content)
	}

	lineNum := 1
	for i := 0; i < pos; i++ {
		if content[i] == '\n' {
			lineNum++
		}
	}

	return NewNumber(float64(lineNum)), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "line-number-at-pos",
		Category:    "position",
		Summary:     "Return the line number of the current point position",
		Description: "Returns the vertical position of the point within the buffer. Line numbers are 1-based, where line 1 is the first line. Lines are separated by newline characters. Calculated by counting newlines from the beginning of the buffer to the current point position.",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{Description: "Get line number at beginning", Input: `line-number-at-pos`, Buffer: "Hello world", Output: "1"},
			{Description: "Get line number on second line", Input: `goto-char 15; line-number-at-pos`, Buffer: "First line\nSecond line\nThird line", Output: "2"},
			{Description: "Get line number on third line", Input: `goto-char 25; line-number-at-pos`, Buffer: "First line\nSecond line\nThird line", Output: "3"},
		},
		SeeAlso: []string{"current-column", "point", "goto-line", "beginning-of-line", "end-of-line"},
	})
}
