package edlisp

import (
	"fmt"
)

func BuiltinEndOfBuffer(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("end-of-buffer expects 0 arguments, got %d", len(args))
	}

	content := buffer.String()
	buffer.SetPoint(len(content) + 1)
	return NewString(""), nil
}
