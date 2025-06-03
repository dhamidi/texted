package edlisp

import (
	"fmt"
	"strings"
)

func BuiltinDowncase(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("downcase expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("downcase expects a string argument")
	}

	str := args[0].(*String)
	return NewString(strings.ToLower(str.Value)), nil
}
