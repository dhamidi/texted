package edlisp

import (
	"fmt"
)

// BuiltinBackwardKillWord deletes text from the current point backward by the specified number of words.
// A word is defined as a sequence of letter characters (as determined by isLetter).
// The function follows the same word boundary logic as backward-word: it skips over non-word
// characters to find the end of each word, then deletes from the beginning of that word to the current point.
// If no count is provided, deletes backward by 1 word. The point moves to the beginning of the deleted region.
func BuiltinBackwardKillWord(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("backward-kill-word expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("backward-kill-word expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}
	
	content := buffer.String()
	startPos := buffer.Point() - 1 // Convert to 0-based
	pos := startPos
	
	// Use the same logic as backward-word to find where to move backward to
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
	
	// Delete from pos to startPos+1 (to include the character at startPos)
	// Handle case where startPos is at or beyond end of buffer
	endIndex := startPos + 1
	if endIndex > len(content) {
		endIndex = len(content)
	}
	newContent := content[:pos] + content[endIndex:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	buffer.SetPoint(pos + 1) // Convert back to 1-based
	
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "backward-kill-word",
		Summary:     "Delete text backward by a specified number of words",
		Description: "Deletes text from the current point backward by the specified number of words. A word is defined as a sequence of letter characters. The function follows the same word boundary logic as backward-word: it skips over non-word characters to find the end of each word, then deletes from the beginning of that word to the current point. If no count is provided, deletes backward by 1 word. The point moves to the beginning of the deleted region.",
		Category:    "editing",
		Parameters: []ParameterDoc{
			{
				Name:        "count",
				Type:        "number",
				Description: "Number of words to delete backward (default: 1)",
				Optional:    true,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Delete one word backward",
				Input:       `goto-char 17; backward-kill-word; buffer-substring 1 -1`,
				Buffer:      "Hello world test",
				Output:      "Hello world ",
			},
			{
				Description: "Delete multiple words backward",
				Input:       `goto-char 25; backward-kill-word 2; buffer-substring 1 -1`,
				Buffer:      "Hello world test buffer content",
				Output:      "Hello world ontent",
			},
		},
		SeeAlso: []string{"kill-word", "backward-word", "delete-region", "kill-line"},
	})
}
