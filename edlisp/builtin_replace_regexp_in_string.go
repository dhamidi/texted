package edlisp

import (
	"fmt"
	"regexp"
)

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
