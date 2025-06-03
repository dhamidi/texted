package edlisp

import (
	"fmt"
	"regexp"
	"strings"
)

// BuiltinLookingBack checks if the text before the current point matches the given pattern.
// Returns the symbol 't' if the pattern matches ending at the current position, 'nil' otherwise.
// The pattern can be either a literal string or a regular expression.
// If the pattern is a valid regular expression, it uses regexp matching on text before the point.
// If the pattern is not a valid regexp, it falls back to literal string matching.
// For literal strings, checks if the text before point ends with the given string.
// For regexps, checks if there's a match that ends exactly at the current point.
// Does not move the point or modify the buffer in any way.
func BuiltinLookingBack(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("looking-back expects 1 argument, got %d", len(args))
	}
	
	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("looking-back expects a string argument")
	}
	
	pattern := args[0].(*String)
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	if pos <= 0 {
		return NewSymbol("nil"), nil
	}
	
	// Try to compile as regular expression
	re, err := regexp.Compile(pattern.Value)
	if err != nil {
		// If not a valid regex, treat as literal string
		if pos >= len(pattern.Value) && strings.HasSuffix(content[:pos], pattern.Value) {
			return NewSymbol("t"), nil
		}
		return NewSymbol("nil"), nil
	}
	
	// Use regular expression matching on text before point
	beforeText := content[:pos]
	match := re.FindStringIndex(beforeText)
	if match != nil && match[1] == len(beforeText) {
		return NewSymbol("t"), nil
	}
	
	return NewSymbol("nil"), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "looking-back",
		Summary:     "Check if text before current position matches a pattern",
		Description: "Checks if the text before the current point matches the given pattern ending at the current position. Returns the symbol 't' if the pattern matches, 'nil' otherwise. The pattern can be either a literal string or a regular expression. For literal strings, checks if the text before point ends with the given string. For regular expressions, checks if there's a match that ends exactly at the current point. Does not move the point or modify the buffer in any way.",
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
				Description: "Check literal string match before point",
				Input:       `goto-char 12; looking-back "world"`,
				Buffer:      "Hello world test buffer",
				Output:      "Returns 't' (text before position 12 ends with 'world')",
			},
			{
				Description: "Check regex pattern match before point",
				Input:       `goto-char 10; looking-back "[0-9]+"`,
				Buffer:      "Hello 123 world",
				Output:      "Returns 't' (digits end at position 10)",
			},
			{
				Description: "Pattern does not match before point",
				Input:       `goto-char 8; looking-back "test"`,
				Buffer:      "Hello world test",
				Output:      "Returns 'nil' (text before position 8 doesn't end with 'test')",
			},
		},
		SeeAlso: []string{"looking-at", "re-search-backward", "string-match"},
	})
}
