package edlisp

import (
	"fmt"
)

func BuiltinBackwardWord(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("backward-word expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("backward-word expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	for i := 0; i < count && pos > 0; i++ {
		// Skip current non-word characters
		for pos > 0 && !isLetter(content[pos-1]) {
			pos--
		}
		// Skip word characters to get to beginning of word
		for pos > 0 && isLetter(content[pos-1]) {
			pos--
		}
	}
	
	buffer.SetPoint(pos + 1) // Convert back to 1-based
	return NewString(""), nil
}
