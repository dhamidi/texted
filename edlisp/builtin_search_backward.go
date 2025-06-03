package edlisp

import (
	"fmt"
	"strings"
)

func BuiltinSearchBackward(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("search-backward expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("search-backward expects a string argument")
	}

	str := args[0].(*String)
	content := buffer.String()
	endPos := buffer.Point() - 1 // Convert to 0-based

	if endPos > len(content) {
		endPos = len(content)
	}
	if endPos < 0 {
		return nil, fmt.Errorf("search failed")
	}

	searchArea := content[:endPos]
	index := strings.LastIndex(searchArea, str.Value)
	if index == -1 {
		return nil, fmt.Errorf("search failed")
	}

	// Set point to end of found text
	matchStart := index + 1 // Convert back to 1-based
	matchEnd := matchStart + len(str.Value)
	buffer.SetPoint(matchEnd)
	buffer.lastSearchMatch = str.Value
	buffer.lastSearchStart = matchStart
	buffer.lastSearchEnd = matchEnd

	return NewString(""), nil
}
