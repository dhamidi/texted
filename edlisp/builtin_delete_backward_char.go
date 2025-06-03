package edlisp

import (
	"fmt"
)

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
