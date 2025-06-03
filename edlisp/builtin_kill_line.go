package edlisp

import (
	"fmt"
)

// BuiltinKillLine deletes text from the point to the end of line(s).
// For a single line (count=1), deletes from after the current point to the end
// of the line, preserving the character at the point. For multiple lines, deletes
// entire lines starting from the current point. The point position remains unchanged.
func BuiltinKillLine(args []Value, buffer *Buffer) (Value, error) {
	var count int = 1

	if len(args) > 1 {
		return nil, fmt.Errorf("kill-line expects at most 1 argument, got %d", len(args))
	}

	if len(args) == 1 {
		if !IsA(args[0], TheNumberKind) {
			return nil, fmt.Errorf("kill-line expects a number argument")
		}
		count = int(args[0].(*Number).Value)
	}

	content := buffer.String()
	pos := buffer.Point() - 1 // Convert to 0-based

	if pos < 0 || pos >= len(content) {
		return NewString(""), nil
	}

	var startPos, lineEnd int

	if count == 1 {
		// For single line kill, preserve cursor character and kill from after cursor
		startPos = pos + 1
		lineEnd = startPos
		// Find end of current line
		for lineEnd < len(content) && content[lineEnd] != '\n' {
			lineEnd++
		}
	} else {
		// For multi-line kill, kill entire lines starting from cursor
		startPos = pos
		lineEnd = startPos
		for i := 0; i < count; i++ {
			// Find end of current line
			for lineEnd < len(content) && content[lineEnd] != '\n' {
				lineEnd++
			}
			// Include the newline character if present
			if lineEnd < len(content) && content[lineEnd] == '\n' {
				lineEnd++
			}
		}
	}

	newContent := content[:startPos] + content[lineEnd:]
	buffer.content.Reset()
	buffer.content.WriteString(newContent)

	return NewString(""), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "kill-line",
		Summary:     "Delete text from the point to the end of line(s)",
		Description: "Deletes text from the current point to the end of the specified number of lines. For a single line (count=1), deletes from after the current point to the end of the line, preserving the character at the point. For multiple lines, deletes entire lines starting from the current point. The point position remains unchanged after the operation.",
		Category:    "editing",
		Parameters: []ParameterDoc{
			{
				Name:        "count",
				Type:        "number",
				Description: "Number of lines to kill (default: 1)",
				Optional:    true,
			},
		},
		Examples: []ExampleDoc{
			{
				Description: "Kill to end of current line",
				Input:       `goto-char 8; kill-line`,
				Buffer:      "First line content\nSecond line content\nThird line",
				Output:      "First li\nSecond line content\nThird line",
			},
			{
				Description: "Kill multiple lines",
				Input:       `goto-char 1; kill-line 2`,
				Buffer:      "First line\nSecond line\nThird line\nFourth line",
				Output:      "Third line\nFourth line",
			},
		},
		SeeAlso: []string{"delete-line", "delete-region", "end-of-line", "kill-word"},
	})
}
