package edlisp

import (
	"fmt"
)

func BuiltinPoint(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("point expects 0 arguments, got %d", len(args))
	}
	
	return NewNumber(float64(buffer.Point())), nil
}
