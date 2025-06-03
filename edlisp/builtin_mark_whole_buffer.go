package edlisp

import (
	"fmt"
)

// BuiltinMarkWholeBuffer marks the entire buffer contents.
// This function creates a region that encompasses all text in the buffer by setting
// the mark at the beginning of the buffer (position 1) and moving the point to the
// end of the buffer. This is equivalent to manually setting mark at position 1 and
// then moving to the end of the buffer.
func BuiltinMarkWholeBuffer(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("mark-whole-buffer expects 0 arguments, got %d", len(args))
	}

	content := buffer.String()
	buffer.SetMark(1)
	buffer.SetPoint(len(content) + 1)

	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "mark-whole-buffer",
		Summary:     "Mark the entire buffer contents",
		Description: "Marks the entire buffer contents by setting the mark at the beginning of the buffer (position 1) and moving the point to the end of the buffer. This creates a region that encompasses all text in the buffer.",
		Category:    "mark",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{
				Description: "Mark entire buffer",
				Input:       `mark-whole-buffer; mark; point`,
				Buffer:      "Hello world test buffer content",
				Output:      "Mark at position 1, point at end of buffer",
			},
		},
		SeeAlso: []string{"mark-line", "mark-word", "beginning-of-buffer", "end-of-buffer"},
	})
}
