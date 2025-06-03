package edlisp

import (
	"fmt"
)

// BuiltinBeginningOfLine moves the point to the beginning of the current line.
// The beginning of a line is defined as the position immediately after a newline
// character, or the start of the buffer if on the first line.
func BuiltinBeginningOfLine(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("beginning-of-line expects 0 arguments, got %d", len(args))
	}

	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based

	if pos < 0 {
		buffer.SetPoint(1)
		return NewString(""), nil
	}
	if pos >= len(content) {
		pos = len(content) - 1
	}

	// Move backward to find beginning of line
	for pos > 0 && content[pos-1] != '\n' {
		pos--
	}

	buffer.SetPoint(pos + 1) // Convert back to 1-based
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "beginning-of-line",
		Summary:     "Move point to the beginning of the current line",
		Description: "Moves the point to the beginning of the current line. The beginning of a line is defined as the position immediately after a newline character, or the start of the buffer if on the first line. This function takes no arguments.",
		Category:    "movement",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{
				Description: "Move to beginning of second line",
				Input:       `search-forward "5"; beginning-of-line; point`,
				Buffer:      "1\n3 5\n7",
				Output:      "3",
			},
			{
				Description: "Move to beginning from middle of line",
				Input:       `goto-char 5; beginning-of-line; point`,
				Buffer:      "Hello world",
				Output:      "1",
			},
		},
		SeeAlso: []string{"end-of-line", "beginning-of-buffer", "goto-line"},
	})
}
