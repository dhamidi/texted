package edlisp

import (
	"fmt"
)

func BuiltinSetMark(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("set-mark expects 0 arguments, got %d", len(args))
	}

	buffer.SetMark(buffer.Point())
	return NewString(""), nil
}
