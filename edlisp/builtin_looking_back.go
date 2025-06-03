package edlisp

import (
	"fmt"
	"regexp"
	"strings"
)

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
