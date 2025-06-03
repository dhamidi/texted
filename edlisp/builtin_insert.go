package edlisp

import (
	"fmt"
)

// BuiltinInsert inserts the given string at the current point position.
// The point is moved to after the inserted text.
func BuiltinInsert(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("insert expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("insert expects a string argument")
	}

	str := args[0].(*String)
	buffer.Insert(str.Value)
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "insert",
		Summary:     "Insert text at the current point position",
		Description: "Inserts the given string at the current point position. The point is moved to after the inserted text. This is the basic function for adding text to the buffer.",
		Category:    "editing",
		Parameters: []ParameterDoc{
			{
				Name:        "text",
				Type:        "string",
				Description: "Text to insert into the buffer",
				Optional:    false,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Insert text into empty buffer",
				Input:       `insert "hello, world"`,
				Buffer:      "",
				Output:      "hello, world",
			},
			{
				Description: "Insert text at specific position",
				Input:       `goto-char 6; insert " beautiful"`,
				Buffer:      "hello world",
				Output:      "hello beautiful world",
			},
		},
		SeeAlso: []string{"delete-char", "replace-region", "kill-line"},
	})
}
