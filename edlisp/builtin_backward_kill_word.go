package edlisp

import (
	"fmt"
)

func BuiltinBackwardKillWord(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1

	if len(args) > 1 {
		return nil, fmt.Errorf("backward-kill-word expects at most 1 argument, got %d", len(args))
	}

	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("backward-kill-word expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}

	content := buffer.String()
	startPos := buffer.Point() - 1 // Convert to 0-based
	pos := startPos

	// Use the same logic as backward-word to find where to move backward to
	for i := 0; i < count && pos > 0; i++ {
		// Skip current non-word characters
		for pos > 0 && !isLetter(content[pos-1]) {
			pos--
		}
		// Skip word characters to get to beginning of word
		for pos > 0 && isLetter(content[pos-1]) {
			pos--
		}
	}

	// Delete from pos to startPos+1 (to include the character at startPos)
	// Handle case where startPos is at or beyond end of buffer
	endIndex := startPos + 1
	if endIndex > len(content) {
		endIndex = len(content)
	}
	newContent := content[:pos] + content[endIndex:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)

	buffer.SetPoint(pos + 1) // Convert back to 1-based

	return NewString(""), nil
}
