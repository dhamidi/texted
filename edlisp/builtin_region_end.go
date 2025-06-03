package edlisp

import (
	"fmt"
)

// BuiltinRegionEnd returns the position of the end of the current region.
// The region is defined by the point and mark positions. This function returns the
// larger of the two positions, ensuring that the end is always the rightmost
// position regardless of whether point is before or after mark.
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

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "region-end",
		Summary:     "Return the position of the end of the current region",
		Description: "Returns the position of the end of the current region. The region is defined by the point and mark positions. This function returns the larger of the two positions, ensuring that the end is always the rightmost position regardless of whether point is before or after mark.",
		Category:    "region",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{
				Description: "Get region end when point is after mark",
				Input:       `goto-char 5; set-mark; goto-char 10; region-end`,
				Buffer:      "Hello world test",
				Output:      "Returns 10 (point position, which is larger)",
			},
			{
				Description: "Get region end when mark is after point",
				Input:       `goto-char 10; set-mark; goto-char 5; region-end`,
				Buffer:      "Hello world test",
				Output:      "Returns 10 (mark position, which is larger)",
			},
		},
		SeeAlso: []string{"region-beginning", "mark", "point", "buffer-substring"},
	})
}
