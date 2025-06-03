package edlisp

import (
	"fmt"
)

// BuiltinDeleteChar deletes characters starting at the current point position.
// By default, deletes one character forward from the point. The point position
// remains unchanged after deletion.
func BuiltinDeleteChar(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1

	if len(args) > 1 {
		return nil, fmt.Errorf("delete-char expects at most 1 argument, got %d", len(args))
	}

	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("delete-char expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}

	content := buffer.String()
	pos := buffer.Point() // 1-based position

	if pos < 1 || pos > len(content) {
		return NewString(""), nil
	}

	// Convert 1-based position to 0-based index
	startPos := pos - 1
	endPos := startPos + count
	if endPos > len(content) {
		endPos = len(content)
	}

	// Delete characters starting at current position
	newContent := content[:startPos] + content[endPos:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)

	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "delete-char",
		Summary:     "Delete characters starting at the current point position",
		Description: "Deletes the specified number of characters starting at the current point position. By default, deletes one character forward from the point. The point position remains unchanged after deletion. If the count exceeds the available characters, deletes up to the end of the buffer.",
		Category:    "editing",
		Parameters: []ParameterDoc{
			{
				Name:        "count",
				Type:        "number",
				Description: "Number of characters to delete (default: 1)",
				Optional:    true,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Delete one character at point",
				Input:       `goto-char 6; delete-char`,
				Buffer:      "Hello world",
				Output:      "Helloworld",
			},
			{
				Description: "Delete multiple characters",
				Input:       `delete-char 3`,
				Buffer:      "Hello",
				Output:      "lo",
			},
		},
		SeeAlso: []string{"delete-backward-char", "delete-region", "kill-line", "insert"},
	})
}
