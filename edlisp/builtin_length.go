package edlisp

import (
	"fmt"
)

func BuiltinLength(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("length expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("length expects a string argument")
	}

	str := args[0].(*String)
	return NewNumber(float64(len(str.Value))), nil
}
