package edlisp

import (
	"fmt"
	"strings"
)

func BuiltinUpcase(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("upcase expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("upcase expects a string argument")
	}

	str := args[0].(*String)
	return NewString(strings.ToUpper(str.Value)), nil
}
