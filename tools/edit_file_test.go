package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestExecuteScript(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		script   string
		expected string
		wantErr  bool
	}{
		{
			name:     "insert text at beginning",
			input:    "world",
			script:   "beginning-of-buffer\ninsert \"hello \"",
			expected: "hello world",
			wantErr:  false,
		},
		{
			name:     "replace entire buffer",
			input:    "old text",
			script:   "mark-whole-buffer\nreplace-region \"new text\"",
			expected: "new text",
			wantErr:  false,
		},
		{
			name:     "insert at end",
			input:    "hello",
			script:   "end-of-buffer\ninsert \" world\"",
			expected: "hello world",
			wantErr:  false,
		},
		{
			name:    "invalid script",
			input:   "test",
			script:  "invalid-function",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExecuteScript(tt.input, tt.script)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ExecuteScript() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ExecuteScript() error = %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("ExecuteScript() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestEditFileHandler(t *testing.T) {
	// Create a temporary file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	// Write initial content
	err := os.WriteFile(testFile, []byte("hello world"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a mock request
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"script": "mark-whole-buffer\nreplace-region \"HELLO WORLD\"",
				"files":  []string{testFile},
			},
		},
	}

	// Call the handler
	result, err := EditFileHandler(ctx, request)
	if err != nil {
		t.Fatalf("EditFileHandler() error = %v", err)
	}

	// Check that result is not an error
	if result == nil {
		t.Fatal("EditFileHandler() returned nil result")
	}

	// Check file was modified
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read modified file: %v", err)
	}

	expected := "HELLO WORLD"
	if string(content) != expected {
		t.Errorf("File content = %q, want %q", string(content), expected)
	}

	// Check result message contains success
	textContent, ok := mcp.AsTextContent(result.Content[0])
	if !ok {
		t.Fatal("Result content is not text content")
	}
	if !strings.Contains(textContent.Text, "Successfully edited") {
		t.Errorf("Result should contain success message, got: %s", textContent.Text)
	}
}

func TestEditFileHandler_NonexistentFile(t *testing.T) {
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"script": `insert "test"`,
				"files":  []string{"/nonexistent/file.txt"},
			},
		},
	}

	result, err := EditFileHandler(ctx, request)
	if err != nil {
		t.Fatalf("EditFileHandler() should not return error for file issues: %v", err)
	}

	// Should return a result with error message
	textContent, ok := mcp.AsTextContent(result.Content[0])
	if !ok {
		t.Fatal("Result content is not text content")
	}
	if !strings.Contains(textContent.Text, "Failed to read") {
		t.Errorf("Result should contain error message for nonexistent file")
	}
}
