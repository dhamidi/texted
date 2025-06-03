package edlisp

import (
	"fmt"
)

// BuiltinBufferSize returns the total number of characters in the buffer.
//
// This function calculates the size of the buffer content in characters (bytes).
// The size includes all characters including newlines, spaces, and special characters.
// This is useful for determining buffer boundaries or calculating buffer statistics.
//
// The function takes no parameters and always succeeds unless called with arguments.
//
// Returns:
//   - number: The total character count in the buffer
//
// Examples:
//   buffer-size → 23 (for buffer "Hello world test buffer")
//   buffer-size → 0 (for empty buffer)
//
// Related functions:
//   - point-max: Maximum valid position in buffer (buffer-size + 1)
//   - point-min: Minimum valid position in buffer (always 1)
//
// Category: buffer
func BuiltinBufferSize(args []Value, buffer *Buffer) (Value, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("buffer-size expects 0 arguments, got %d", len(args))
	}
	
	content := buffer.String()
	return NewNumber(float64(len(content))), nil
}

func init() {
	RegisterDocumentation(FunctionDoc{
		Name:        "buffer-size",
		Category:    "buffer",
		Summary:     "Return the total number of characters in the buffer",
		Description: "Calculates the size of the buffer content in characters (bytes). The size includes all characters including newlines, spaces, and special characters. This is useful for determining buffer boundaries or calculating buffer statistics.",
		Parameters:  []ParameterDoc{},
		Examples: []ExampleDoc{
			{Description: "Get size of buffer with content", Input: `buffer-size`, Buffer: "Hello world test buffer", Output: "23"},
			{Description: "Get size of empty buffer", Input: `buffer-size`, Buffer: "", Output: "0"},
		},
		SeeAlso: []string{"point-max", "point-min", "point"},
	})
}
