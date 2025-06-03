package edlisp

import (
	"fmt"
)

// BuiltinRegionBeginning returns the position of the beginning of the current region.
// The region is defined by the point and mark positions. This function returns the
// smaller of the two positions, ensuring that the beginning is always the leftmost
// position regardless of whether point is before or after mark.
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

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "region-beginning",
		Summary:     "Return the position of the beginning of the current region",
		Description: "Returns the position of the beginning of the current region. The region is defined by the point and mark positions. This function returns the smaller of the two positions, ensuring that the beginning is always the leftmost position regardless of whether point is before or after mark.",
		Category:    "region",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{
				Description: "Get region beginning when mark is before point",
				Input:       `goto-char 5; set-mark; goto-char 10; region-beginning`,
				Buffer:      "Hello world test",
				Output:      "Returns 5 (mark position, which is smaller)",
			},
			{
				Description: "Get region beginning when point is before mark",
				Input:       `goto-char 10; set-mark; goto-char 5; region-beginning`,
				Buffer:      "Hello world test",
				Output:      "Returns 5 (point position, which is smaller)",
			},
		},
		SeeAlso: []string{"region-end", "mark", "point", "buffer-substring"},
	})
}
