package edlisp

import (
	"fmt"
)

func BuiltinKillWord(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("kill-word expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("kill-word expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	startPos := buffer.Point() - 1 // Convert to 0-based
	pos := startPos
	
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
	
	newContent := content[:startPos] + content[pos:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	return NewString(""), nil
}
