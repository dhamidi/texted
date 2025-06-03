package edlisp

import (
	"fmt"
	"strings"
)

// BuiltinGotoLine moves the point to the beginning of the specified line number.
// Line numbers are 1-based. If the line number is less than 1, moves to line 1.
// If the line number is greater than the total number of lines, moves to the last line.
// The point is positioned at the beginning of the target line.
func BuiltinGotoLine(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("goto-line expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheNumberKind) {
		return nil, fmt.Errorf("goto-line expects a number argument")
	}

	num := args[0].(*Number)
	lineNum := int(num.Value)
	
	content := buffer.String()
	lines := strings.Split(content, "\n")
	
	if lineNum < 1 {
		lineNum = 1
	} else if lineNum > len(lines) {
		lineNum = len(lines)
	}
	
	// Calculate position at beginning of target line
	pos := 1
	for i := 0; i < lineNum-1; i++ {
		pos += len(lines[i]) + 1 // +1 for newline
	}
	
	buffer.SetPoint(pos)
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "goto-line",
		Summary:     "Move point to the beginning of a specific line",
		Description: "Moves the point to the beginning of the specified line number. Line numbers are 1-based. If the line number is less than 1, moves to line 1. If the line number is greater than the total number of lines, moves to the last line. The point is positioned at the beginning of the target line.",
		Category:    "movement",
		Parameters: []ParameterDoc{
			{
				Name:        "line-number",
				Type:        "number",
				Description: "The line number to move to (1-based)",
				Optional:    false,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Move to line 3 in a multi-line buffer",
				Input:       `goto-line 3`,
				Buffer:      "Line 1: First line of text\nLine 2: Second line of text\nLine 3: Third line of text\nLine 4: Fourth line of text",
				Output:      "Point moves to beginning of line 3 (position 51)",
			},
			{
				Description: "Move to line beyond buffer end",
				Input:       `goto-line 10`,
				Buffer:      "Line 1\nLine 2\nLine 3",
				Output:      "Point moves to beginning of last line (line 3)",
			},
			{
				Description: "Move to line number less than 1",
				Input:       `goto-line 0`,
				Buffer:      "Line 1\nLine 2\nLine 3",
				Output:      "Point moves to beginning of first line (line 1)",
			},
		},
		SeeAlso: []string{"goto-char", "beginning-of-line", "end-of-line", "line-number-at-pos"},
	})
}
