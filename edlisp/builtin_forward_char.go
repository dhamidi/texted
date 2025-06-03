package edlisp

import (
	"fmt"
)

// BuiltinForwardChar moves the point forward by the specified number of characters.
// If no count is provided, moves forward by 1 character. The point cannot move
// beyond the end of the buffer or before the beginning.
func BuiltinForwardChar(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1

	if len(args) > 1 {
		return nil, fmt.Errorf("forward-char expects at most 1 argument, got %d", len(args))
	}

	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("forward-char expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}

	content := buffer.String()
	newPos := buffer.Point() + count

	if newPos < 1 {
		newPos = 1
	} else if newPos > len(content)+1 {
		newPos = len(content) + 1
	}

	buffer.SetPoint(newPos)
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "forward-char",
		Summary:     "Move point forward by a specified number of characters",
		Description: "Moves the point forward by the specified number of characters. If no count is provided, moves forward by 1 character. The point is constrained to stay within buffer bounds - it cannot move beyond the end of the buffer or before the beginning.",
		Category:    "movement",
		Parameters: []ParameterDoc{
			{
				Name:        "count",
				Type:        "number",
				Description: "Number of characters to move forward (default: 1)",
				Optional:    true,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Move forward by default amount (1 character)",
				Input:       `goto-char 5; forward-char; point`,
				Buffer:      "Hello world",
				Output:      "6",
			},
			{
				Description: "Move forward by specific count",
				Input:       `goto-char 1; forward-char 3; point`,
				Buffer:      "Hello world",
				Output:      "4",
			},
		},
		SeeAlso: []string{"backward-char", "forward-word", "goto-char"},
	})
}
