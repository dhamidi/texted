package edlisp

import (
	"fmt"
)

func BuiltinKillLine(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("kill-line expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("kill-line expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	if pos < 0 || pos >= len(content) {
		return NewString(""), nil
	}
	
	// Kill multiple lines if count > 1
	lineEnd := pos
	for i := 0; i < count; i++ {
		// Find end of current line
		for lineEnd < len(content) && content[lineEnd] != '\n' {
			lineEnd++
		}
		// Include the newline character if present and if we're killing multiple lines
		if lineEnd < len(content) && content[lineEnd] == '\n' && (i < count-1 || count > 1) {
			lineEnd++
		}
	}
	
	newContent := content[:pos] + content[lineEnd:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	return NewString(""), nil
}
