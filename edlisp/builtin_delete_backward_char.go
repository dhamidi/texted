package edlisp

import (
	"fmt"
)

// BuiltinDeleteBackwardChar deletes characters backward from the current point position.
// By default, deletes one character backward from the point. The point is moved to the
// beginning of the deleted region after deletion.
func BuiltinDeleteBackwardChar(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1

	if len(args) > 1 {
		return nil, fmt.Errorf("delete-backward-char expects at most 1 argument, got %d", len(args))
	}

	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("delete-backward-char expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}

	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based

	// Delete count characters before the current position
	endPos := pos
	startPos := pos - count
	if startPos < 0 {
		startPos = 0
	}

	if startPos >= endPos {
		return NewString(""), nil
	}

	newContent := content[:startPos] + content[endPos:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)

	// Update point to the new position
	buffer.SetPoint(startPos + 1)

	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "delete-backward-char",
		Summary:     "Delete characters backward from the current point position",
		Description: "Deletes the specified number of characters backward from the current point position. By default, deletes one character backward from the point. The point is moved to the beginning of the deleted region after deletion. If the count exceeds the available characters before the point, deletes up to the beginning of the buffer.",
		Category:    "editing",
		Parameters: []ParameterDoc{
			{
				Name:        "count",
				Type:        "number",
				Description: "Number of characters to delete backward (default: 1)",
				Optional:    true,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Delete one character backward",
				Input:       `forward-char; delete-backward-char`,
				Buffer:      "123",
				Output:      "23",
			},
			{
				Description: "Delete multiple characters backward",
				Input:       `forward-char 2; delete-backward-char 2`,
				Buffer:      "Hello",
				Output:      "llo",
			},
		},
		SeeAlso: []string{"delete-char", "backward-char", "delete-region", "kill-line"},
	})
}
