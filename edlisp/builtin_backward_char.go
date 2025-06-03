package edlisp

import (
	"fmt"
)

func BuiltinBackwardChar(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1

	if len(args) > 1 {
		return nil, fmt.Errorf("backward-char expects at most 1 argument, got %d", len(args))
	}

	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("backward-char expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}

	newPos := buffer.Point() - count

	if newPos < 1 {
		newPos = 1
	}

	buffer.SetPoint(newPos)
	return NewString(""), nil
}
