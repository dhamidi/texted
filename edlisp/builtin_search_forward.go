package edlisp

import (
	"fmt"
	"strings"
)

// BuiltinSearchForward searches for the given string forward from the current point.
// If found, moves point to the end of the match and returns an empty string.
// If not found, returns an error and leaves point unchanged.
// The function stores information about the last search match for use with replace-match.
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

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "search-forward",
		Summary:     "Search for text forward from current position",
		Description: "Searches for the given string forward from the current point. If found, moves point to the end of the match and stores match information for use with replace-match. If not found, returns an error and leaves point unchanged.",
		Category:    "search",
		Parameters: []ParameterDoc{
			{
				Name:        "pattern",
				Type:        "string",
				Description: "Text pattern to search for",
				Optional:    false,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Basic text search",
				Input:       `search-forward "test"`,
				Buffer:      "Hello world, this is a test buffer.",
				Output:      "Point moves to position after 'test' (position 28)",
			},
			{
				Description: "Search that fails",
				Input:       `search-forward "missing"`,
				Buffer:      "Hello world",
				Output:      "Error: search failed",
			},
		},
		SeeAlso: []string{"search-backward", "re-search-forward", "replace-match"},
	})
}
