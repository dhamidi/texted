package edlisp

import (
	"fmt"
)

func BuiltinBufferSubstring(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("buffer-substring expects 2 arguments, got %d", len(args))
	}

	if !IsA(args[0], TheNumberKind) || !IsA(args[1], TheNumberKind) {
		return nil, fmt.Errorf("buffer-substring expects number arguments")
	}

	start := int(args[0].(*Number).Value)
	end := int(args[1].(*Number).Value)
	
	content := buffer.String()
	
	// Handle special case: -1 means end of buffer
	if end == -1 {
		end = len(content) + 1
	}
	
	start-- // Convert to 0-based
	end--   // Convert to 0-based
	
	if start < 0 {
		start = 0
	}
	if end > len(content) {
		end = len(content)
	}
	if start >= end {
		return NewString(""), nil
	}
	
	return NewString(content[start:end]), nil
}
