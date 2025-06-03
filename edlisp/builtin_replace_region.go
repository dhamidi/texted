package edlisp

import (
	"fmt"
)

func BuiltinReplaceRegion(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("replace-region expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("replace-region expects a string argument")
	}

	str := args[0].(*String)

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

	newContent := content[:start] + str.Value + content[end:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)

	// Set point to end of replacement
	buffer.SetPoint(start + len(str.Value) + 1)

	return NewString(""), nil
}
