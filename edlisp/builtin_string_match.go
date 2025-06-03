package edlisp

import (
	"fmt"
	"regexp"
	"strings"
)

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
