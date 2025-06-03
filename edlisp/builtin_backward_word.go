package edlisp

import (
	"fmt"
)

// BuiltinBackwardWord moves the point backward by the specified number of words.
// A word is defined as a sequence of letter characters (as determined by isLetter).
// The function skips over non-word characters to find the end of each word,
// then moves to the beginning of that word. If no count is provided, moves backward by 1 word.
// The point cannot move before the beginning of the buffer.
func BuiltinBackwardWord(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("backward-word expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("backward-word expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	for i := 0; i < count && pos > 0; i++ {
		// Skip current non-word characters
		for pos > 0 && !isLetter(content[pos-1]) {
			pos--
		}
		// Skip word characters to get to beginning of word
		for pos > 0 && isLetter(content[pos-1]) {
			pos--
		}
	}
	
	buffer.SetPoint(pos + 1) // Convert back to 1-based
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "backward-word",
		Summary:     "Move point backward by a specified number of words",
		Description: "Moves the point backward by the specified number of words. A word is defined as a sequence of letter characters. The function skips over non-word characters to find the end of each word, then moves to the beginning of that word. If no count is provided, moves backward by 1 word. The point cannot move before the beginning of the buffer.",
		Category:    "movement",
		Parameters: []ParameterDoc{
			{
				Name:        "count",
				Type:        "number",
				Description: "Number of words to move backward (default: 1)",
				Optional:    true,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Move backward by default amount (1 word)",
				Input:       `goto-char 15; backward-word; point`,
				Buffer:      "Hello world test",
				Output:      "13",
			},
			{
				Description: "Move backward by specific count",
				Input:       `goto-char 20; backward-word 2; point`,
				Buffer:      "Hello world test buffer",
				Output:      "13",
			},
		},
		SeeAlso: []string{"forward-word", "backward-char", "backward-kill-word", "mark-word"},
	})
}
