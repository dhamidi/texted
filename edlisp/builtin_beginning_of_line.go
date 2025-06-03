package edlisp

import (
	"fmt"
)

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
