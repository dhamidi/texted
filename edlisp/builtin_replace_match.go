package edlisp

import (
	"fmt"
)

// BuiltinReplaceMatch replaces the text of the last successful search match with new text.
// Takes one argument: the replacement string.
// This function requires that a search operation (search-forward, search-backward, re-search-forward,
// or re-search-backward) has been performed previously to establish match boundaries.
// Replaces the matched text with the provided replacement string and moves point to the end
// of the replacement text. Returns an empty string on success.
// If no previous search has been performed, returns an error.
func BuiltinReplaceMatch(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("replace-match expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("replace-match expects a string argument")
	}

	str := args[0].(*String)
	
	if buffer.lastSearchMatch == "" {
		return nil, fmt.Errorf("no previous search")
	}
	
	content := buffer.String()
	start := buffer.lastSearchStart - 1 // Convert to 0-based
	end := buffer.lastSearchEnd - 1     // Convert to 0-based
	
	if start < 0 || end > len(content) || start >= end {
		return nil, fmt.Errorf("invalid search match positions")
	}
	
	newContent := content[:start] + str.Value + content[end:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	// Update point to end of replacement
	buffer.SetPoint(start + len(str.Value) + 1)
	
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "replace-match",
		Summary:     "Replace text of last search match with new text",
		Description: "Replaces the text of the last successful search match with new text. Takes one argument: the replacement string. This function requires that a search operation (search-forward, search-backward, re-search-forward, or re-search-backward) has been performed previously to establish match boundaries. Replaces the matched text with the provided replacement string and moves point to the end of the replacement text. Returns an empty string on success. If no previous search has been performed, returns an error.",
		Category:    "search",
		Parameters: []ParameterDoc{
			{
				Name:        "replacement",
				Type:        "string",
				Description: "Text to replace the last search match with",
				Optional:    false,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Replace matched text after search",
				Input:       `search-forward "old"; replace-match "new"`,
				Buffer:      "Hello old world, this is a test.",
				Output:      "Buffer becomes 'Hello new world, this is a test.' and point moves to end of replacement",
			},
			{
				Description: "Replace regex match",
				Input:       `re-search-forward "[0-9]+"; replace-match "NUM"`,
				Buffer:      "Version 123 released",
				Output:      "Buffer becomes 'Version NUM released' and point moves after 'NUM'",
			},
		},
		SeeAlso: []string{"search-forward", "search-backward", "re-search-forward", "re-search-backward", "replace-regexp-in-string"},
	})
}
