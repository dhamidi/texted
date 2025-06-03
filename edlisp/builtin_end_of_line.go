package edlisp

import (
	"fmt"
)

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
