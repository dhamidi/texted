package edlisp

import (
	"fmt"
	"strings"
)

func BuiltinGotoLine(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("goto-line expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheNumberKind) {
		return nil, fmt.Errorf("goto-line expects a number argument")
	}

	num := args[0].(*Number)
	lineNum := int(num.Value)
	
	content := buffer.String()
	lines := strings.Split(content, "\n")
	
	if lineNum < 1 {
		lineNum = 1
	} else if lineNum > len(lines) {
		lineNum = len(lines)
	}
	
	// Calculate position at beginning of target line
	pos := 1
	for i := 0; i < lineNum-1; i++ {
		pos += len(lines[i]) + 1 // +1 for newline
	}
	
	buffer.SetPoint(pos)
	return NewString(""), nil
}
