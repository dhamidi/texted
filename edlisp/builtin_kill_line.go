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
	
	var startPos, lineEnd int
	
	if count == 1 {
		// For single line kill, preserve cursor character and kill from after cursor
		startPos = pos + 1
		lineEnd = startPos
		// Find end of current line
		for lineEnd < len(content) && content[lineEnd] != '\n' {
			lineEnd++
		}
	} else {
		// For multi-line kill, kill entire lines starting from cursor
		startPos = pos
		lineEnd = startPos
		for i := 0; i < count; i++ {
			// Find end of current line
			for lineEnd < len(content) && content[lineEnd] != '\n' {
				lineEnd++
			}
			// Include the newline character if present
			if lineEnd < len(content) && content[lineEnd] == '\n' {
				lineEnd++
			}
		}
	}
	
	newContent := content[:startPos] + content[lineEnd:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	return NewString(""), nil
}
