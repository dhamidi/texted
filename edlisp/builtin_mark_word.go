package edlisp

import (
	"fmt"
)

func BuiltinMarkWord(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("mark-word expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	if pos < 0 || pos >= len(content) {
		return NewString(""), nil
	}
	
	// Find the start of the word (move backward to find non-letter)
	start := pos
	for start > 0 && isLetter(content[start-1]) {
		start--
	}
	
	// Find the end of the word (move forward to find non-letter)
	end := pos
	for end < len(content) && isLetter(content[end]) {
		end++
	}
	
	// Set mark at beginning of word, point at end
	buffer.SetMark(start + 1)     // Convert back to 1-based
	buffer.SetPoint(end + 1)      // Convert back to 1-based
	
	return NewString(""), nil
}
