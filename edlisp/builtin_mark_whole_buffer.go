package edlisp

import (
	"fmt"
)

func BuiltinMarkWholeBuffer(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("mark-whole-buffer expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	buffer.SetMark(1)
	buffer.SetPoint(len(content) + 1)
	
	return NewString(""), nil
}
