package edlisp

import (
	"fmt"
	"regexp"
	"strings"
)

// BuiltinLookingAt checks if the text at the current point matches the given pattern.
// Returns the symbol 't' if the pattern matches at the current position, 'nil' otherwise.
// The pattern can be either a literal string or a regular expression.
// If the pattern is a valid regular expression, it uses regexp matching.
// If the pattern is not a valid regexp, it falls back to literal string matching.
// Does not move the point or modify the buffer in any way.
func BuiltinLookingAt(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("looking-at expects 1 argument, got %d", len(args))
	}
	
	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("looking-at expects a string argument")
	}
	
	pattern := args[0].(*String)
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	if pos < 0 || pos >= len(content) {
		return NewSymbol("nil"), nil
	}
	
	// Try to compile as regular expression
	re, err := regexp.Compile(pattern.Value)
	if err != nil {
		// If not a valid regex, treat as literal string
		if strings.HasPrefix(content[pos:], pattern.Value) {
			return NewSymbol("t"), nil
		}
		return NewSymbol("nil"), nil
	}
	
	// Use regular expression matching
	match := re.FindStringIndex(content[pos:])
	if match != nil && match[0] == 0 {
		return NewSymbol("t"), nil
	}
	
	return NewSymbol("nil"), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "looking-at",
		Summary:     "Check if text at current position matches a pattern",
		Description: "Checks if the text at the current point matches the given pattern. Returns the symbol 't' if the pattern matches at the current position, 'nil' otherwise. The pattern can be either a literal string or a regular expression. If the pattern is a valid regular expression, it uses regexp matching starting at the current position. If the pattern is not a valid regexp, it falls back to literal string matching. Does not move the point or modify the buffer in any way.",
		Category:    "search",
		Parameters: []ParameterDoc{
			{
				Name:        "pattern",
				Type:        "string",
				Description: "Pattern to match (literal string or regular expression)",
				Optional:    false,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Check literal string match",
				Input:       `goto-char 7; looking-at "world"`,
				Buffer:      "Hello world test buffer",
				Output:      "Returns 't' (pattern matches at position 7)",
			},
			{
				Description: "Check regex pattern match",
				Input:       `goto-char 7; looking-at "[0-9]+"`,
				Buffer:      "Hello 123 world",
				Output:      "Returns 't' (digits match at position 7)",
			},
			{
				Description: "Pattern does not match",
				Input:       `goto-char 7; looking-at "test"`,
				Buffer:      "Hello world test",
				Output:      "Returns 'nil' (pattern not at position 7)",
			},
		},
		SeeAlso: []string{"looking-back", "re-search-forward", "string-match"},
	})
}
