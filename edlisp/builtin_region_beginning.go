package edlisp

import (
	"fmt"
)

func BuiltinRegionBeginning(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("region-beginning expects 0 arguments, got %d", len(args))
	}
	
	start := buffer.Mark()
	end := buffer.Point()
	
	if start > end {
		start = end
	}
	
	return NewNumber(float64(start)), nil
}
