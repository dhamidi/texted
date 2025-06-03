package edlisp

import (
	"fmt"
	"regexp"
	"strings"
)

// BuiltinStringMatch searches for a pattern within a string and returns the index of the first match.
// Takes two arguments: a pattern and a target string to search within.
// The pattern can be either a literal string or a regular expression.
// If the pattern is a valid regular expression, it uses regexp matching.
// If the pattern is not a valid regexp, it falls back to literal string search.
// Returns the 0-based index of the first match as a number, or the symbol 'nil' if no match is found.
// This function operates on string arguments and does not modify the buffer.
func BuiltinStringMatch(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("string-match expects 2 arguments, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) || !IsA(args[1], TheStringKind) {
		return nil, fmt.Errorf("string-match expects string arguments")
	}

	pattern := args[0].(*String)
	str := args[1].(*String)

	// Try to compile as regular expression
	re, err := regexp.Compile(pattern.Value)
	if err != nil {
		// If not a valid regex, treat as literal string
		index := strings.Index(str.Value, pattern.Value)
		if index == -1 {
			return NewSymbol("nil"), nil
		}
		return NewNumber(float64(index)), nil
	}

	// Use regular expression matching
	match := re.FindStringIndex(str.Value)
	if match == nil {
		return NewSymbol("nil"), nil
	}

	return NewNumber(float64(match[0])), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "string-match",
		Summary:     "Search for pattern within a string and return match index",
		Description: "Searches for a pattern within a string and returns the index of the first match. Takes two arguments: a pattern and a target string to search within. The pattern can be either a literal string or a regular expression. If the pattern is a valid regular expression, it uses regexp matching. If the pattern is not a valid regexp, it falls back to literal string search. Returns the 0-based index of the first match as a number, or the symbol 'nil' if no match is found. This function operates on string arguments and does not modify the buffer.",
		Category:    "string",
		Parameters: []ParameterDoc{
			{
				Name:        "pattern",
				Type:        "string",
				Description: "Pattern to search for (literal string or regular expression)",
				Optional:    false,
			},
			{
				Name:        "string",
				Type:        "string",
				Description: "Target string to search within",
				Optional:    false,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Search for literal text in string",
				Input:       `string-match "wor" "Hello world"`,
				Buffer:      "Test buffer",
				Output:      "Returns 6 (index of 'wor' in 'Hello world')",
			},
			{
				Description: "Search for regex pattern in string",
				Input:       `string-match "[0-9]+" "Hello 123 world"`,
				Buffer:      "Test buffer",
				Output:      "Returns 6 (index of first digit sequence)",
			},
			{
				Description: "Pattern not found in string",
				Input:       `string-match "xyz" "Hello world"`,
				Buffer:      "Test buffer",
				Output:      "Returns 'nil' (pattern not found)",
			},
		},
		SeeAlso: []string{"looking-at", "re-search-forward", "replace-regexp-in-string"},
	})
}
