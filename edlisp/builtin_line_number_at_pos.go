package edlisp

import (
	"fmt"
)

func BuiltinLineNumberAtPos(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("line-number-at-pos expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	if pos < 0 {
		pos = 0
	}
	if pos > len(content) {
		pos = len(content)
	}
	
	lineNum := 1
	for i := 0; i < pos; i++ {
		if content[i] == '\n' {
			lineNum++
		}
	}
	
	return NewNumber(float64(lineNum)), nil
}
