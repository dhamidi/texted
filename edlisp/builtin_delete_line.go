package edlisp

import (
	"fmt"
)

func BuiltinDeleteLine(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1

	if len(args) > 1 {
		return nil, fmt.Errorf("delete-line expects at most 1 argument, got %d", len(args))
	}

	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("delete-line expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}

	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based

	// Find beginning of current line
	lineStart := pos
	for lineStart > 0 && content[lineStart-1] != '\n' {
		lineStart--
	}

	// Find end of line(s) based on count
	lineEnd := pos
	for i := 0; i < count; i++ {
		for lineEnd < len(content) && content[lineEnd] != '\n' {
			lineEnd++
		}
		if lineEnd < len(content) && content[lineEnd] == '\n' {
			lineEnd++ // Include the newline
		}
	}

	newContent := content[:lineStart] + content[lineEnd:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)

	buffer.SetPoint(lineStart + 1) // Convert back to 1-based

	return NewString(""), nil
}
