package edlisp

import (
	"fmt"
)

// BuiltinSetMark sets the mark at the current point position.
// The mark serves as a secondary position that, together with point, defines a region.
// This function takes no arguments and always sets the mark to the current point location.
// After calling this function, the region spans from the mark to the current point.
func BuiltinSetMark(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("set-mark expects 0 arguments, got %d", len(args))
	}

	buffer.SetMark(buffer.Point())
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "set-mark",
		Summary:     "Set mark at current point position",
		Description: "Sets the mark at the current point position. The mark serves as a secondary position that, together with point, defines a region. This function takes no arguments and always sets the mark to the current point location.",
		Category:    "mark",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{
				Description: "Set mark at current position",
				Input:       `goto-char 7; set-mark; mark`,
				Buffer:      "Hello world, this is a test buffer.",
				Output:      "Mark is set to position 7, same as current point",
			},
		},
		SeeAlso: []string{"set-mark-command", "mark", "region-beginning", "region-end"},
	})
}
