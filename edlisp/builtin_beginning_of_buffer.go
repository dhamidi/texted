package edlisp

import (
	"fmt"
)

// BuiltinBeginningOfBuffer moves the point to the very beginning of the buffer.
// This is always position 1, regardless of buffer content.
func BuiltinBeginningOfBuffer(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("beginning-of-buffer expects 0 arguments, got %d", len(args))
	}
	
	buffer.SetPoint(1)
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "beginning-of-buffer",
		Summary:     "Move point to the very beginning of the buffer",
		Description: "Moves the point to the very beginning of the buffer, which is always position 1. This is a simple navigation command that provides a quick way to jump to the start of any buffer content.",
		Category:    "movement",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{
				Description: "Move to beginning from any position",
				Input:       `goto-char 20; beginning-of-buffer; point`,
				Buffer:      "Hello world\nSecond line\nThird line",
				Output:      "1",
			},
			{
				Description: "Move to beginning of empty buffer",
				Input:       `beginning-of-buffer; point`,
				Buffer:      "",
				Output:      "1",
			},
		},
		SeeAlso: []string{"end-of-buffer", "beginning-of-line", "goto-char"},
	})
}
