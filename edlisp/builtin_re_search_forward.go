package edlisp

import (
	"fmt"
	"regexp"
)

// BuiltinReSearchForward searches for the given regular expression pattern forward from the current point.
// If found, moves point to the end of the match and returns an empty string.
// If not found, returns an error and leaves point unchanged.
// The function stores information about the last search match for use with replace-match.
// The pattern is compiled as a regular expression using Go's regexp package syntax.
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

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "re-search-forward",
		Summary:     "Search for regular expression pattern forward from current position",
		Description: "Searches for the given regular expression pattern forward from the current point using Go's regexp package syntax. If found, moves point to the end of the match and stores match information for use with replace-match. If not found, returns an error and leaves point unchanged. If the pattern is invalid, returns a compilation error.",
		Category:    "search",
		Parameters: []ParameterDoc{
			{
				Name:        "regexp",
				Type:        "string",
				Description: "Regular expression pattern to search for (Go regexp syntax)",
				Optional:    false,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Search for word followed by digits",
				Input:       `re-search-forward "[a-z]+[0-9]+"`,
				Buffer:      "The function foo123 is defined here.",
				Output:      "Point moves to position 20 (after 'foo123')",
			},
			{
				Description: "Search for digit pattern",
				Input:       `re-search-forward "[0-9]+"`,
				Buffer:      "Hello 123 world",
				Output:      "Point moves to position after first number match",
			},
		},
		SeeAlso: []string{"re-search-backward", "search-forward", "replace-match", "looking-at"},
	})
}
