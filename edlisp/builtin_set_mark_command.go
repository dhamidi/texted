package edlisp

import (
	"fmt"
)

func BuiltinSetMarkCommand(args []Value, buffer *Buffer) (Value, error) {
	var pos int

	if len(args) > 1 {
		return nil, fmt.Errorf("set-mark-command expects at most 1 argument, got %d", len(args))
	}

	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("set-mark-command expects a number argument")
		}
		pos = int(args[0].(*Number).Value)
	} else {
		pos = buffer.Point()
	}

	buffer.SetMark(pos)
	return NewString(""), nil
}
