package edlisp

import (
	"fmt"
)

// BuiltinBackwardChar moves the point backward by the specified number of characters.
// If no count is provided, moves backward by 1 character. The point cannot move
// before the beginning of the buffer.
func BuiltinBackwardChar(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("backward-char expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("backward-char expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	newPos := buffer.Point() - count
	
	if newPos < 1 {
		newPos = 1
	}
	
	buffer.SetPoint(newPos)
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "backward-char",
		Summary:     "Move point backward by a specified number of characters",
		Description: "Moves the point backward by the specified number of characters. If no count is provided, moves backward by 1 character. The point is constrained to stay within buffer bounds - it cannot move before the beginning of the buffer.",
		Category:    "movement",
		Parameters: []ParameterDoc{
			{
				Name:        "count",
				Type:        "number",
				Description: "Number of characters to move backward (default: 1)",
				Optional:    true,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Move backward by default amount (1 character)",
				Input:       `goto-char 6; backward-char; point`,
				Buffer:      "Hello world",
				Output:      "5",
			},
			{
				Description: "Move backward by specific count",
				Input:       `goto-char 8; backward-char 3; point`,
				Buffer:      "Hello world",
				Output:      "5",
			},
		},
		SeeAlso: []string{"forward-char", "backward-word", "goto-char"},
	})
}
