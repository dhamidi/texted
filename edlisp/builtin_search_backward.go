package edlisp

import (
	"fmt"
	"strings"
)

// BuiltinSearchBackward searches for the given string backward from the current point.
// If found, moves point to the end of the match and returns an empty string.
// If not found, returns an error and leaves point unchanged.
// The function stores information about the last search match for use with replace-match.
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

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "search-backward",
		Summary:     "Search for text backward from current position",
		Description: "Searches for the given string backward from the current point. If found, moves point to the end of the match and stores match information for use with replace-match. If not found, returns an error and leaves point unchanged. The search examines text before the current point position.",
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
				Description: "Basic backward text search",
				Input:       `end-of-buffer search-backward "test"`,
				Buffer:      "Hello world, this is a test buffer.\nThe word \"test\" appears twice here.\nAnother line with test content.",
				Output:      "Point moves to position after the last 'test' (position 95)",
			},
			{
				Description: "Search that fails",
				Input:       `beginning-of-buffer search-backward "missing"`,
				Buffer:      "Hello world",
				Output:      "Error: search failed",
			},
		},
		SeeAlso: []string{"search-forward", "re-search-backward", "replace-match"},
	})
}
