package edlisp

import (
	"fmt"
	"regexp"
)

// BuiltinReplaceRegexpInString performs regular expression replacement on a string.
// Takes three arguments: a regular expression pattern, a replacement string, and a target string.
// The pattern must be a valid regular expression using Go's regexp package syntax.
// Returns a new string with all matches of the pattern replaced by the replacement string.
// This is a pure string operation that does not modify the buffer or affect the point.
// If the pattern is invalid, returns a compilation error.
func BuiltinReplaceRegexpInString(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("replace-regexp-in-string expects 3 arguments, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) || !IsA(args[1], TheStringKind) || !IsA(args[2], TheStringKind) {
		return nil, fmt.Errorf("replace-regexp-in-string expects string arguments")
	}

	pattern := args[0].(*String)
	replacement := args[1].(*String)
	str := args[2].(*String)
	
	re, err := regexp.Compile(pattern.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid regexp: %v", err)
	}
	
	result := re.ReplaceAllString(str.Value, replacement.Value)
	return NewString(result), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "replace-regexp-in-string",
		Summary:     "Replace all matches of a regular expression in a string",
		Description: "Performs regular expression replacement on a string. Takes three arguments: a regular expression pattern, a replacement string, and a target string. The pattern must be a valid regular expression using Go's regexp package syntax. Returns a new string with all matches of the pattern replaced by the replacement string. This is a pure string operation that does not modify the buffer or affect the point. If the pattern is invalid, returns a compilation error.",
		Category:    "string",
		Parameters: []ParameterDoc{
			{
				Name:        "regexp",
				Type:        "string",
				Description: "Regular expression pattern to search for (Go regexp syntax)",
				Optional:    false,
			},
			{
				Name:        "replacement",
				Type:        "string",
				Description: "Text to replace each match with",
				Optional:    false,
			},
			{
				Name:        "string",
				Type:        "string",
				Description: "Target string to perform replacements on",
				Optional:    false,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Replace all occurrences of a character",
				Input:       `replace-regexp-in-string "o" "0" "Hello world"`,
				Buffer:      "Test buffer",
				Output:      "Returns 'Hell0 w0rld' (all 'o' characters replaced with '0')",
			},
			{
				Description: "Replace all number sequences with placeholder",
				Input:       `replace-regexp-in-string "[0-9]+" "NUM" "Hello 123 world 456"`,
				Buffer:      "Test buffer",
				Output:      "Returns 'Hello NUM world NUM' (all digit sequences replaced)",
			},
			{
				Description: "Replace multiple word patterns",
				Input:       `replace-regexp-in-string "\\b[a-z]+\\b" "WORD" "hello world test"`,
				Buffer:      "Test buffer",
				Output:      "Returns 'WORD WORD WORD' (all lowercase words replaced)",
			},
		},
		SeeAlso: []string{"string-match", "replace-match", "re-search-forward", "re-search-backward"},
	})
}
