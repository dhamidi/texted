package edlisp

import (
	"fmt"
)

func BuiltinReplaceMatch(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("replace-match expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("replace-match expects a string argument")
	}

	str := args[0].(*String)
	
	if buffer.lastSearchMatch == "" {
		return nil, fmt.Errorf("no previous search")
	}
	
	content := buffer.String()
	start := buffer.lastSearchStart - 1 // Convert to 0-based
	end := buffer.lastSearchEnd - 1     // Convert to 0-based
	
	if start < 0 || end > len(content) || start >= end {
		return nil, fmt.Errorf("invalid search match positions")
	}
	
	newContent := content[:start] + str.Value + content[end:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	// Update point to end of replacement
	buffer.SetPoint(start + len(str.Value) + 1)
	
	return NewString(""), nil
}
