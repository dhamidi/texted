package edlisp

import (
	"fmt"
	"regexp"
)

func BuiltinReSearchBackward(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("re-search-backward expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("re-search-backward expects a string argument")
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

	re, err := regexp.Compile(str.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid regexp: %v", err)
	}

	searchArea := content[:endPos]
	matches := re.FindAllStringIndex(searchArea, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("search failed")
	}

	// Get the last match (rightmost before point)
	match := matches[len(matches)-1]

	// Set point to end of found text
	matchStart := match[0] + 1 // Convert back to 1-based
	matchEnd := match[1] + 1
	buffer.SetPoint(matchEnd)
	buffer.lastSearchMatch = content[match[0]:match[1]]
	buffer.lastSearchStart = matchStart
	buffer.lastSearchEnd = matchEnd

	return NewString(""), nil
}
