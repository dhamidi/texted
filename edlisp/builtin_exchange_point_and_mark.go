package edlisp

import (
	"fmt"
)

// BuiltinExchangePointAndMark swaps the positions of point and mark.
// This function exchanges the current point position with the mark position,
// effectively moving the cursor to where the mark was while setting the mark
// to where the cursor was. This is useful for quickly moving between the two
// ends of a region or for reversing the direction of a region selection.
func BuiltinExchangePointAndMark(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("exchange-point-and-mark expects 0 arguments, got %d", len(args))
	}
	
	point := buffer.Point()
	mark := buffer.Mark()
	
	buffer.SetPoint(mark)
	buffer.SetMark(point)
	
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "exchange-point-and-mark",
		Summary:     "Swap the positions of point and mark",
		Description: "Exchanges the current point position with the mark position, effectively moving the cursor to where the mark was while setting the mark to where the cursor was. This is useful for quickly moving between the two ends of a region or for reversing the direction of a region selection.",
		Category:    "region",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{
				Description: "Exchange point and mark positions",
				Input:       `goto-char 5; set-mark; goto-char 10; exchange-point-and-mark; point`,
				Buffer:      "Hello world test",
				Output:      "Point moves from 10 to 5, mark moves from 5 to 10",
			},
		},
		SeeAlso: []string{"set-mark", "mark", "point", "region-beginning", "region-end"},
	})
}
