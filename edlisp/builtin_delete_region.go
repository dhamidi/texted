package edlisp

import (
	"fmt"
)

// BuiltinDeleteRegion deletes the text between the mark and point.
// The region is defined by the current mark and point positions, with the smaller
// position used as the start and the larger as the end. After deletion, the point
// is positioned at the beginning of the deleted region.
func BuiltinDeleteRegion(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("delete-region expects 0 arguments, got %d", len(args))
	}
	
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
	
	newContent := content[:start] + content[end:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)
	
	// Set point to start of deleted region
	buffer.SetPoint(start + 1)
	
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "delete-region",
		Summary:     "Delete the text between the mark and point",
		Description: "Deletes the text between the current mark and point positions. The region is defined by the mark and point, with the smaller position used as the start and the larger as the end. After deletion, the point is positioned at the beginning of the deleted region. If no mark is set, the behavior is undefined.",
		Category:    "editing",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{
				Description: "Delete selected region",
				Input:       `goto-char 7; set-mark; goto-char 10; delete-region`,
				Buffer:      "Hello old world, this is a test.",
				Output:      "Hello  world, this is a test.",
			},
		},
		SeeAlso: []string{"set-mark", "mark", "point", "replace-region", "kill-line"},
	})
}
