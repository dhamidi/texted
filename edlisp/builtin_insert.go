package edlisp

import "fmt"

// BuiltinInsert implements the insert function that inserts text at the current point.
func BuiltinInsert(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("insert expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("insert expects a string argument")
	}

	str := args[0].(*String)
	buffer.Insert(str.Value)
	return NewString(""), nil
}