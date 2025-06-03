package edlisp

import (
	"fmt"
)

// BuiltinDeleteLine deletes entire lines starting from the line containing the point.
// By default, deletes one line. When deleting multiple lines, includes the newline
// characters. After deletion, the point is positioned at the beginning of the line
// where deletion started.
func BuiltinDeleteLine(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1

	if len(args) > 1 {
		return nil, fmt.Errorf("delete-line expects at most 1 argument, got %d", len(args))
	}

	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("delete-line expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}

	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based

	// Find beginning of current line
	lineStart := pos
	for lineStart > 0 && content[lineStart-1] != '\n' {
		lineStart--
	}

	// Find end of line(s) based on count
	lineEnd := pos
	for i := 0; i < count; i++ {
		for lineEnd < len(content) && content[lineEnd] != '\n' {
			lineEnd++
		}
		if lineEnd < len(content) && content[lineEnd] == '\n' {
			lineEnd++ // Include the newline
		}
	}

	newContent := content[:lineStart] + content[lineEnd:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)

	buffer.SetPoint(lineStart + 1) // Convert back to 1-based

	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "delete-line",
		Summary:     "Delete entire lines starting from the line containing the point",
		Description: "Deletes the specified number of complete lines starting from the line containing the current point. By default, deletes one line. When deleting multiple lines, includes the newline characters. After deletion, the point is positioned at the beginning of the line where deletion started.",
		Category:    "editing",
		Parameters: []ParameterDoc{
			{
				Name:        "count",
				Type:        "number",
				Description: "Number of lines to delete (default: 1)",
				Optional:    true,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Delete one line",
				Input:       `goto-char 15; delete-line`,
				Buffer:      "First line\nSecond line\nThird line",
				Output:      "First line\nThird line",
			},
			{
				Description: "Delete multiple lines",
				Input:       `goto-char 15; delete-line 2`,
				Buffer:      "First line\nSecond line\nThird line\nFourth line",
				Output:      "First line\nFourth line",
			},
		},
		SeeAlso: []string{"kill-line", "delete-region", "beginning-of-line", "end-of-line"},
	})
}
