package edlisp

import (
	"fmt"
)

func BuiltinDeleteRegion(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("delete-region expects 0 arguments, got %d", len(args))
	}
	
	start := buffer.Mark()
	end := buffer.Point()
	
	if start > end {
		start, end = end, start
	}
	
	content := buffer.String()
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
	
	newContent := content[:start] + content[end:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	// Set point to start of deleted region
	buffer.SetPoint(start + 1)
	
	return NewString(""), nil
}
