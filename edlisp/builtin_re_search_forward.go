package edlisp

import (
	"fmt"
	"regexp"
)

func BuiltinReSearchForward(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("re-search-forward expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("re-search-forward expects a string argument")
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
	
	re, err := regexp.Compile(str.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid regexp: %v", err)
	}
	
	match := re.FindStringIndex(content[startPos:])
	if match == nil {
		return nil, fmt.Errorf("search failed")
	}
	
	// Set point to end of found text
	matchStart := startPos + match[0] + 1 // Convert back to 1-based
	matchEnd := startPos + match[1] + 1
	buffer.SetPoint(matchEnd)
	buffer.lastSearchMatch = content[startPos+match[0] : startPos+match[1]]
	buffer.lastSearchStart = matchStart
	buffer.lastSearchEnd = matchEnd
	
	return NewString(""), nil
}
