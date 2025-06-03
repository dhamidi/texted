package edlisp

import (
	"fmt"
)

// BuiltinForwardWord moves the point forward by the specified number of words.
// A word is defined as a sequence of letter characters (as determined by isLetter).
// The function skips over non-word characters to find the start of each word,
// then moves to the end of that word. If no count is provided, moves forward by 1 word.
// The point cannot move beyond the end of the buffer.
func BuiltinForwardWord(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("forward-word expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("forward-word expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	for i := 0; i < count && pos < len(content); i++ {
		// Skip non-word characters to get to a word
		for pos < len(content) && !isLetter(content[pos]) {
			pos++
		}
		// Skip current word characters to get to end of word
		for pos < len(content) && isLetter(content[pos]) {
			pos++
		}
	}
	
	buffer.SetPoint(pos + 1) // Convert back to 1-based
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "forward-word",
		Summary:     "Move point forward by a specified number of words",
		Description: "Moves the point forward by the specified number of words. A word is defined as a sequence of letter characters. The function skips over non-word characters to find the start of each word, then moves to the end of that word. If no count is provided, moves forward by 1 word. The point cannot move beyond the end of the buffer.",
		Category:    "movement",
		Parameters: []ParameterDoc{
			{
				Name:        "count",
				Type:        "number",
				Description: "Number of words to move forward (default: 1)",
				Optional:    true,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Move forward by default amount (1 word)",
				Input:       `goto-char 3; forward-word; point`,
				Buffer:      "Hello world test",
				Output:      "6",
			},
			{
				Description: "Move forward by specific count",
				Input:       `goto-char 1; forward-word 2; point`,
				Buffer:      "Hello world test buffer",
				Output:      "12",
			},
		},
		SeeAlso: []string{"backward-word", "forward-char", "kill-word", "mark-word"},
	})
}
