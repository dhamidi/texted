package edlisp

import (
	"fmt"
	"strings"
)

func BuiltinConcat(args []Value, buffer *Buffer) (Value, error) {
	var result strings.Builder

	for i, arg := range args {
		if !IsA(arg, TheStringKind) {
			return nil, fmt.Errorf("concat expects string arguments, got non-string at position %d", i+1)
		}
		str := arg.(*String)
		result.WriteString(str.Value)
	}

	return NewString(result.String()), nil
}
