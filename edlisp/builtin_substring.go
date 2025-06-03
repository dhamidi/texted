package edlisp

import (
	"fmt"
)

func BuiltinSubstring(args []Value, buffer *Buffer) (Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("substring expects 2 or 3 arguments, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("substring expects a string as first argument")
	}

	if !IsA(args[1], TheNumberKind) {
		return nil, fmt.Errorf("substring expects a number as second argument")
	}

	str := args[0].(*String)
	start := int(args[1].(*Number).Value)
	end := len(str.Value)

	if len(args) == 3 {
		if !IsA(args[2], TheNumberKind) {
			return nil, fmt.Errorf("substring expects a number as third argument")
		}
		end = int(args[2].(*Number).Value)
	}

	// Convert from 1-based to 0-based indexing
	start--

	// For two-argument form, end should be to the end of string
	if len(args) == 2 {
		end = len(str.Value)
	} else {
		// For three-argument form, end is 1-based and exclusive
		// Convert to 0-based exclusive by decrementing
		end--
	}

	// Bounds checking
	if start < 0 {
		start = 0
	}
	if end > len(str.Value) {
		end = len(str.Value)
	}
	if start > end {
		return NewString(""), nil
	}

	result := str.Value[start:end]
	return NewString(result), nil
}
