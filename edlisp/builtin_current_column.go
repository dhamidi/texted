package edlisp

import (
	"fmt"
)

func BuiltinCurrentColumn(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("current-column expects 0 arguments, got %d", len(args))
	}

	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based

	if pos < 0 {
		return NewNumber(0), nil
	}
	if pos >= len(content) {
		pos = len(content) - 1
	}

	column := 0
	// Count backward to find beginning of line
	for i := pos; i >= 0 && content[i] != '\n'; i-- {
		column++
	}

	return NewNumber(float64(column)), nil
}
