package edlisp

import (
	"fmt"
)

// BuiltinMarkLine marks one or more lines starting from the current line.
// This function creates a region that encompasses complete lines, including their newline characters.
// The mark is positioned at the beginning of the current line, and the point is moved to the
// end of the specified number of lines. When count is 1 (default), it marks just the current line.
// When count is greater than 1, it marks multiple consecutive lines starting from the current line.
func BuiltinMarkLine(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1
	
	if len(args) > 1 {
		return nil, fmt.Errorf("mark-line expects at most 1 argument, got %d", len(args))
	}
	
	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("mark-line expects a number argument")
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
	
	buffer.SetMark(lineStart + 1) // Convert back to 1-based
	buffer.SetPoint(lineEnd + 1)  // Convert back to 1-based
	
	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "mark-line",
		Summary:     "Mark one or more complete lines",
		Description: "Marks one or more lines starting from the current line. This function creates a region that encompasses complete lines, including their newline characters. The mark is positioned at the beginning of the current line, and the point is moved to the end of the specified number of lines.",
		Category:    "mark",
		Parameters: []ParameterDoc{
			{
				Name:        "count",
				Type:        "number",
				Description: "Number of lines to mark, starting from current line. Defaults to 1",
				Optional:    true,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Mark current line (default behavior)",
				Input:       `search-forward "Second"; mark-line; buffer-substring (region-beginning) (region-end)`,
				Buffer:      "First line\nSecond line\nThird line",
				Output:      "Marks 'Second line\\n'",
			},
			{
				Description: "Mark multiple lines",
				Input:       `goto-char 15; mark-line 2`,
				Buffer:      "First line\nSecond line\nThird line\nFourth line",
				Output:      "Marks 'Second line\\nThird line\\n'",
			},
		},
		SeeAlso: []string{"mark-word", "mark-whole-buffer", "beginning-of-line", "end-of-line"},
	})
}
