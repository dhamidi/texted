package edlisp

import (
	"fmt"
	"strings"
)

func BuiltinCapitalize(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("capitalize expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("capitalize expects a string argument")
	}

	str := args[0].(*String)
	if len(str.Value) == 0 {
		return NewString(""), nil
	}
	
	result := strings.ToUpper(string(str.Value[0])) + strings.ToLower(str.Value[1:])
	return NewString(result), nil
}
