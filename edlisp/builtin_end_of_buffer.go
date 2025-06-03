package edlisp

import (
	"fmt"
)

// BuiltinEndOfBuffer moves the point to the very end of the buffer.
// This is always one position past the last character in the buffer.
func BuiltinEndOfBuffer(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("end-of-buffer expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	buffer.SetPoint(len(content) + 1)
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "end-of-buffer",
		Summary:     "Move point to the very end of the buffer",
		Description: "Moves the point to the very end of the buffer, which is one position past the last character. This position allows for inserting text at the end of the buffer content.",
		Category:    "movement",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{
				Description: "Move to end from any position",
				Input:       `goto-char 1; end-of-buffer; point`,
				Buffer:      "Hello world",
				Output:      "12",
			},
			{
				Description: "Move to end of empty buffer",
				Input:       `end-of-buffer; point`,
				Buffer:      "",
				Output:      "1",
			},
		},
		SeeAlso: []string{"beginning-of-buffer", "end-of-line", "goto-char"},
	})
}
