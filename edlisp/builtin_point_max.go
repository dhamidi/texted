package edlisp

import (
	"fmt"
)

func BuiltinPointMax(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("point-max expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	return NewNumber(float64(len(content) + 1)), nil
}
