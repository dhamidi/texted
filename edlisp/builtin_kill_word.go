package edlisp

import (
	"fmt"
)

// BuiltinKillWord deletes text from the current point forward by the specified number of words.
// A word is defined as a sequence of letter characters (as determined by isLetter).
// The function follows the same word boundary logic as forward-word: it skips over non-word
// characters to find the start of each word, then deletes to the end of that word.
// If no count is provided, deletes forward by 1 word. The point remains at its original position.
func BuiltinKillWord(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1

	if len(args) > 1 {
		return nil, fmt.Errorf("kill-word expects at most 1 argument, got %d", len(args))
	}

	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("kill-word expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}

	content := buffer.String()
	startPos := buffer.Point() - 1 // Convert to 0-based
	pos := startPos

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

	newContent := content[:startPos] + content[pos:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)

	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "kill-word",
		Summary:     "Delete text forward by a specified number of words",
		Description: "Deletes text from the current point forward by the specified number of words. A word is defined as a sequence of letter characters. The function follows the same word boundary logic as forward-word: it skips over non-word characters to find the start of each word, then deletes to the end of that word. If no count is provided, deletes forward by 1 word. The point remains at its original position after deletion.",
		Category:    "editing",
		Parameters: []ParameterDoc{
			{
				Name:        "count",
				Type:        "number",
				Description: "Number of words to delete forward (default: 1)",
				Optional:    true,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Delete one word forward from middle of word",
				Input:       `goto-char 3; kill-word; buffer-substring 1 -1`,
				Buffer:      "Hello world test",
				Output:      "He world test",
			},
			{
				Description: "Delete multiple words forward",
				Input:       `goto-char 1; kill-word 2; buffer-substring 1 -1`,
				Buffer:      "Hello world test buffer content",
				Output:      " test buffer content",
			},
		},
		SeeAlso: []string{"backward-kill-word", "forward-word", "delete-region", "kill-line"},
	})
}
