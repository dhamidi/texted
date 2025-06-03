package edlisp

import (
	"fmt"
)

func BuiltinForwardWord(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("forward-word expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("forward-word expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	for i := 0; i < count && pos < len(content); i++ {
		// Skip non-word characters to get to a word
		for pos < len(content) && !isLetter(content[pos]) {
			pos++
		}
		// Skip current word characters to get to end of word
		for pos < len(content) && isLetter(content[pos]) {
			pos++
		}
	}
	
	buffer.SetPoint(pos + 1) // Convert back to 1-based
	return NewString(""), nil
}
