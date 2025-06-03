package edlisp

import (
	"fmt"
)

// BuiltinMarkWord marks the word at or after the current point position.
// This function identifies word boundaries using letter characters and sets up a region
// that encompasses the entire word. The mark is positioned at the beginning of the word,
// and the point is moved to the end of the word, creating a region that selects the word.
// If the point is not currently on a letter, the function still attempts to find word boundaries.
func BuiltinMarkWord(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("mark-word expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based
	
	if pos < 0 || pos >= len(content) {
		return NewString(""), nil
	}
	
	// Find the start of the word (move backward to find non-letter)
	start := pos
	for start > 0 && isLetter(content[start-1]) {
		start--
	}
	
	// Find the end of the word (move forward to find non-letter)
	end := pos
	for end < len(content) && isLetter(content[end]) {
		end++
	}
	
	// Set mark at beginning of word, point at end
	buffer.SetMark(start + 1)     // Convert back to 1-based
	buffer.SetPoint(end + 1)      // Convert back to 1-based
	
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "mark-word",
		Summary:     "Mark the word at or after current position",
		Description: "Marks the word at or after the current point position. This function identifies word boundaries using letter characters and sets up a region that encompasses the entire word. The mark is positioned at the beginning of the word, and the point is moved to the end of the word.",
		Category:    "mark",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{
				Description: "Mark word from within the word",
				Input:       `goto-char 7; mark-word; region-beginning; region-end`,
				Buffer:      "Hello world, this is a test buffer.",
				Output:      "Marks 'world' - mark at position 7, point at position 12",
			},
		},
		SeeAlso: []string{"mark-line", "mark-whole-buffer", "forward-word", "backward-word"},
	})
}
