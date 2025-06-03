package edlisp

import (
	"fmt"
)

func BuiltinExchangePointAndMark(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("exchange-point-and-mark expects 0 arguments, got %d", len(args))
	}

	point := buffer.Point()
	mark := buffer.Mark()

	buffer.SetPoint(mark)
	buffer.SetMark(point)

	return NewString(""), nil
}
