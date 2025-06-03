package edlisp

import (
	"fmt"
)

func BuiltinBeginningOfBuffer(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("beginning-of-buffer expects 0 arguments, got %d", len(args))
	}

	buffer.SetPoint(1)
	return NewString(""), nil
}
