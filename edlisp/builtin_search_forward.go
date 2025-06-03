package edlisp

import (
	"fmt"
	"strings"
)

func BuiltinSearchForward(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("search-forward expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("search-forward expects a string argument")
	}

	str := args[0].(*String)
	content := buffer.String()
	startPos := buffer.Point() - 1 // Convert to 0-based

	if startPos < 0 {
		startPos = 0
	}
	if startPos >= len(content) {
		return nil, fmt.Errorf("search failed")
	}

	index := strings.Index(content[startPos:], str.Value)
	if index == -1 {
		return nil, fmt.Errorf("search failed")
	}

	// Set point to end of found text
	matchStart := startPos + index + 1 // Convert back to 1-based
	matchEnd := matchStart + len(str.Value)
	buffer.SetPoint(matchEnd)
	buffer.lastSearchMatch = str.Value
	buffer.lastSearchStart = matchStart
	buffer.lastSearchEnd = matchEnd

	return NewString(""), nil
}
