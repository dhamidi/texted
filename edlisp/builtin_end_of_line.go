package edlisp

import (
	"fmt"
)

// BuiltinEndOfLine moves the point to the end of the current line.
// The end of a line is defined as the position just before a newline character,
// or the end of the buffer if on the last line.
func BuiltinEndOfLine(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("end-of-line expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	if pos < 0 {
		pos = 0
	}
	if pos >= len(content) {
		buffer.SetPoint(len(content) + 1)
		return NewString(""), nil
	}
	
	// Move forward to find end of line (newline character or end of content)
	for pos < len(content) && content[pos] != '\n' {
		pos++
	}
	
	// pos now points to newline or beyond end of content
	// We want to be at the last character of the line, not the newline
	buffer.SetPoint(pos + 1) // Convert back to 1-based
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "end-of-line",
		Summary:     "Move point to the end of the current line",
		Description: "Moves the point to the end of the current line. The end of a line is defined as the position just before a newline character, or the end of the buffer if on the last line. This function takes no arguments.",
		Category:    "movement",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{
				Description: "Move to end of current line",
				Input:       `search-forward "T"; end-of-line; point`,
				Buffer:      "One\nTwo\nThree",
				Output:      "8",
			},
			{
				Description: "Move to end from beginning of line",
				Input:       `beginning-of-line; end-of-line; point`,
				Buffer:      "Hello world",
				Output:      "12",
			},
		},
		SeeAlso: []string{"beginning-of-line", "end-of-buffer", "goto-line"},
	})
}
