package edlisp

import (
	"fmt"
)

func BuiltinForwardChar(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("forward-char expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("forward-char expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	newPos := buffer.Point() + count
	
	if newPos < 1 {
		newPos = 1
	} else if newPos > len(content)+1 {
		newPos = len(content) + 1
	}
	
	buffer.SetPoint(newPos)
	return NewString(""), nil
}
