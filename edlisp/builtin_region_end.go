package edlisp

import (
	"fmt"
)

func BuiltinRegionEnd(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("region-end expects 0 arguments, got %d", len(args))
	}
	
	start := buffer.Mark()
	end := buffer.Point()
	
	if start > end {
		end = start
	}
	
	return NewNumber(float64(end)), nil
}
