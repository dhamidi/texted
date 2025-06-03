package edlisp

import (
	"fmt"
)

// BuiltinSetMarkCommand sets the mark at a specified position or at the current point.
// This function provides more flexibility than set-mark by accepting an optional position argument.
// When called without arguments, it behaves identically to set-mark (sets mark at current point).
// When called with a position argument, it sets the mark at that specific position.
func BuiltinSetMarkCommand(args []Value, buffer *Buffer) (Value, error) {
	var pos int

	if len(args) > 1 {
		return nil, fmt.Errorf("set-mark-command expects at most 1 argument, got %d", len(args))
	}

	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("set-mark-command expects a number argument")
		}
		pos = int(args[0].(*Number).Value)
	} else {
		pos = buffer.Point()
	}

	buffer.SetMark(pos)
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "set-mark-command",
		Summary:     "Set mark at specified position or current point",
		Description: "Sets the mark at a specified position or at the current point. When called without arguments, it behaves identically to set-mark (sets mark at current point). When called with a position argument, it sets the mark at that specific position.",
		Category:    "mark",
		Parameters: []ParameterDoc{
			{
				Name:        "position",
				Type:        "number",
				Description: "Buffer position where to set the mark (1-based). If omitted, uses current point",
				Optional:    true,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Set mark at specific position",
				Input:       `set-mark-command 5; mark`,
				Buffer:      "Hello world test",
				Output:      "Mark is set to position 5",
			},
			{
				Description: "Set mark at current point (no argument)",
				Input:       `goto-char 8; set-mark-command; mark`,
				Buffer:      "Hello world test",
				Output:      "Mark is set to position 8 (current point)",
			},
		},
		SeeAlso: []string{"set-mark", "mark", "goto-char", "region-beginning", "region-end"},
	})
}
