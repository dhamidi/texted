package edlisp

import (
	"fmt"
)

func BuiltinDeleteChar(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("delete-char expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("delete-char expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	pos := buffer.Point() // 1-based position
	
	if pos < 1 || pos > len(content) {
		return NewString(""), nil
	}
	
	// Convert 1-based position to 0-based index
	startPos := pos - 1
	endPos := startPos + count
	if endPos > len(content) {
		endPos = len(content)
	}
	
	// Delete characters starting at current position
	newContent := content[:startPos] + content[endPos:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	return NewString(""), nil
}
