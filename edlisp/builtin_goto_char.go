package edlisp

import (
	"fmt"
)

func BuiltinGotoChar(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("goto-char expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheNumberKind) {
		return nil, fmt.Errorf("goto-char expects a number argument")
	}

	num := args[0].(*Number)
	pos := int(num.Value)

	content := buffer.String()
	if pos < 1 {
		pos = 1
	} else if pos > len(content)+1 {
		pos = len(content) + 1
	}

	buffer.SetPoint(pos)
	return NewString(""), nil
}
