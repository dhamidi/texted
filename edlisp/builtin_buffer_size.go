package edlisp

import (
	"fmt"
)

func BuiltinBufferSize(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("buffer-size expects 0 arguments, got %d", len(args))
	}

	content := buffer.String()
	return NewNumber(float64(len(content))), nil
}
