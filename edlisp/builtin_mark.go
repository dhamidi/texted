package edlisp

import (
	"fmt"
)

func BuiltinMark(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("mark expects 0 arguments, got %d", len(args))
	}

	return NewNumber(float64(buffer.Mark())), nil
}
