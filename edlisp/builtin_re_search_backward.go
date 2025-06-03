package edlisp

import (
	"fmt"
	"regexp"
)

// BuiltinReSearchBackward searches for the given regular expression pattern backward from the current point.
// If found, moves point to the end of the rightmost match before the current position and returns an empty string.
// If not found, returns an error and leaves point unchanged.
// The function stores information about the last search match for use with replace-match.
// The pattern is compiled as a regular expression using Go's regexp package syntax.
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

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "re-search-backward",
		Summary:     "Search for regular expression pattern backward from current position",
		Description: "Searches for the given regular expression pattern backward from the current point using Go's regexp package syntax. Finds all matches before the current position and selects the rightmost (closest to point) match. If found, moves point to the end of the match and stores match information for use with replace-match. If not found, returns an error and leaves point unchanged. If the pattern is invalid, returns a compilation error.",
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
				Description: "Search backward for literal text",
				Input:       `end-of-buffer; re-search-backward "two"`,
				Buffer:      "One two three",
				Output:      "Point moves to position 8 (after 'two')",
			},
			{
				Description: "Search backward for pattern",
				Input:       `end-of-buffer; re-search-backward "[a-z]+"`,
				Buffer:      "Hello 123 world",
				Output:      "Point moves to position after rightmost word match",
			},
		},
		SeeAlso: []string{"re-search-forward", "search-backward", "replace-match", "looking-back"},
	})
}
