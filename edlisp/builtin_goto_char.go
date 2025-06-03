package edlisp

import (
	"fmt"
)

// BuiltinGotoChar moves the point to the specified character position in the buffer.
// The position is 1-based, where position 1 is the beginning of the buffer.
// If the position is less than 1, the point moves to the beginning (position 1).
// If the position is greater than the buffer size plus 1, the point moves to the end.
// This function is commonly used for precise cursor positioning in text editing operations.
func BuiltinGotoChar(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("goto-char expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheNumberKind) {
		return nil, fmt.Errorf("goto-char expects a number argument")
	}

	num := args[0].(*Number)
	pos := int(num.Value)
	
	content := buffer.String()
	if pos < 1 {
		pos = 1
	} else if pos > len(content)+1 {
		pos = len(content) + 1
	}
	
	buffer.SetPoint(pos)
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "goto-char",
		Summary:     "Move point to specified character position",
		Description: "Moves the point to the specified character position in the buffer. The position is 1-based, where position 1 is the beginning of the buffer. If the position is less than 1, the point moves to the beginning. If the position is greater than the buffer size plus 1, the point moves to the end. This function provides precise cursor positioning for text editing operations.",
		Category:    "movement",
		Parameters: []ParameterDoc{
			{
				Name:        "position",
				Type:        "number",
				Description: "1-based character position to move to",
				Optional:    false,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Move to specific position in buffer",
				Input:       `goto-char 7`,
				Buffer:      "Hello world, this is a test buffer.",
				Output:      "Point moves to position 7 (after 'world')",
			},
			{
				Description: "Move to beginning of buffer",
				Input:       `goto-char 1`,
				Buffer:      "Hello world",
				Output:      "Point moves to position 1 (beginning)",
			},
			{
				Description: "Move beyond buffer end (clamped)",
				Input:       `goto-char 100`,
				Buffer:      "Hello world",
				Output:      "Point moves to position 12 (end of buffer)",
			},
			{
				Description: "Use with column calculation",
				Input:       `goto-char 25; current-column`,
				Buffer:      "First line\nSecond line with content\nThird line",
				Output:      "Returns column 14 on second line",
			},
		},
		SeeAlso: []string{"point", "goto-line", "beginning-of-buffer", "end-of-buffer", "current-column"},
	})
}
