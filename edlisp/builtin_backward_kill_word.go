package edlisp

import (
	"fmt"
)

func BuiltinBackwardKillWord(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("backward-kill-word expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("backward-kill-word expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	endPos := buffer.Point() - 1 // Convert to 0-based
	pos := endPos
	
	for i := 0; i < count && pos > 0; i++ {
		// Move backward to start of previous word
		// First skip any non-word characters before current position
		for pos > 0 && !isLetter(content[pos-1]) {
			pos--
		}
		// Then skip the word characters to get to beginning of word
		for pos > 0 && isLetter(content[pos-1]) {
			pos--
		}
	}
	
	newContent := content[:pos] + content[endPos:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	buffer.SetPoint(pos + 1) // Convert back to 1-based
	
	return NewString(""), nil
}
