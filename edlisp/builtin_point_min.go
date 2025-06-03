package edlisp

import (
	"fmt"
)

func BuiltinPointMin(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("point-min expects 0 arguments, got %d", len(args))
	}

	return NewNumber(1), nil
}
