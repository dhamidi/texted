package edlisp

import (
	"fmt"
)

// BuiltinReplaceRegion replaces the text between mark and point with the given string.
// The region is automatically normalized so it doesn't matter which of mark or point
// comes first. After replacement, point is positioned at the end of the new text.
func BuiltinReplaceRegion(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("replace-region expects 1 argument, got %d", len(args))
	}

	if !IsA(args[0], TheStringKind) {
		return nil, fmt.Errorf("replace-region expects a string argument")
	}

	str := args[0].(*String)
	
	start := buffer.Mark()
	end := buffer.Point()
	
	if start > end {
		start, end = end, start
	}
	
	content := buffer.String()
	start-- // Convert to 0-based
	end--   // Convert to 0-based
	
	if start < 0 {
		start = 0
	}
	if end > len(content) {
		end = len(content)
	}
	if start >= end {
		return NewString(""), nil
	}
	
	newContent := content[:start] + str.Value + content[end:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	// Set point to end of replacement
	buffer.SetPoint(start + len(str.Value) + 1)
	
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "replace-region",
		Summary:     "Replace text between mark and point with new text",
		Description: "Replaces the text between mark and point with the given string. The region is automatically normalized so it doesn't matter which of mark or point comes first. After replacement, point is positioned at the end of the new text. This is useful for making targeted edits to specific parts of the buffer.",
		Category:    "editing",
		Parameters: []ParameterDoc{
			{
				Name:        "replacement",
				Type:        "string",
				Description: "Text to replace the selected region with",
				Optional:    false,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Replace selected text",
				Input:       `goto-char 7; set-mark; goto-char 10; replace-region "new"`,
				Buffer:      "Hello old world, this is a test.",
				Output:      "Hello new world, this is a test.",
			},
			{
				Description: "Replace with empty string (delete region)",
				Input:       `mark-word; replace-region ""`,
				Buffer:      "delete this word",
				Output:      " this word (first word deleted)",
			},
		},
		SeeAlso: []string{"set-mark", "mark-word", "mark-line", "delete-region", "insert"},
	})
}
