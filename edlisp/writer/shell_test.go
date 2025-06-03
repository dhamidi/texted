package writer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dhamidi/texted/edlisp"
)

func TestShellWriter_WriteValue_SimpleCommand(t *testing.T) {
	writer := &ShellWriter{}
	var buf bytes.Buffer

	value := edlisp.NewList(
		edlisp.NewSymbol("search-forward"),
		edlisp.NewString("text"),
	)

	err := writer.WriteValue(&buf, value)
	if err != nil {
		t.Fatalf("WriteValue failed: %v", err)
	}

	expected := "search-forward \"text\""
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestShellWriter_WriteValue_EmptyList(t *testing.T) {
	writer := &ShellWriter{}
	var buf bytes.Buffer

	value := edlisp.NewEmptyList()
	err := writer.WriteValue(&buf, value)
	if err != nil {
		t.Fatalf("WriteValue failed: %v", err)
	}

	expected := ""
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestShellWriter_WriteValue_SingleSymbol(t *testing.T) {
	writer := &ShellWriter{}
	var buf bytes.Buffer

	value := edlisp.NewList(edlisp.NewSymbol("set-mark"))
	err := writer.WriteValue(&buf, value)
	if err != nil {
		t.Fatalf("WriteValue failed: %v", err)
	}

	expected := "set-mark"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestShellWriter_WriteValue_CommandWithNumber(t *testing.T) {
	tests := []struct {
		name     string
		value    *edlisp.Number
		expected string
	}{
		{"integer", edlisp.NewIntNumber(42), "move 42"},
		{"float", edlisp.NewNumber(3.14), "move 3.14"},
		{"zero", edlisp.NewIntNumber(0), "move 0"},
		{"negative", edlisp.NewIntNumber(-5), "move -5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &ShellWriter{}
			var buf bytes.Buffer

			value := edlisp.NewList(
				edlisp.NewSymbol("move"),
				tt.value,
			)

			err := writer.WriteValue(&buf, value)
			if err != nil {
				t.Fatalf("WriteValue failed: %v", err)
			}

			if buf.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestShellWriter_WriteValue_StringQuoting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple string", "hello", "test \"hello\""},
		{"string with spaces", "hello world", "test \"hello world\""},
		{"string with quotes", "hello \"world\"", "test \"hello \\\"world\\\"\""},
		{"string with newlines", "hello\nworld", "test \"hello\\nworld\""},
		{"string with backslashes", "hello\\world", "test \"hello\\\\world\""},
		{"empty string", "", "test \"\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &ShellWriter{}
			var buf bytes.Buffer

			value := edlisp.NewList(
				edlisp.NewSymbol("test"),
				edlisp.NewString(tt.input),
			)

			err := writer.WriteValue(&buf, value)
			if err != nil {
				t.Fatalf("WriteValue failed: %v", err)
			}

			if buf.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestShellWriter_WriteValue_MultipleArguments(t *testing.T) {
	writer := &ShellWriter{}
	var buf bytes.Buffer

	value := edlisp.NewList(
		edlisp.NewSymbol("command"),
		edlisp.NewString("arg1"),
		edlisp.NewIntNumber(42),
		edlisp.NewString("arg3"),
	)

	err := writer.WriteValue(&buf, value)
	if err != nil {
		t.Fatalf("WriteValue failed: %v", err)
	}

	expected := "command \"arg1\" 42 \"arg3\""
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestShellWriter_Write_MultipleExpressions(t *testing.T) {
	writer := &ShellWriter{}
	var buf bytes.Buffer

	expressions := []edlisp.Value{
		edlisp.NewList(
			edlisp.NewSymbol("search-forward"),
			edlisp.NewString("doIt"),
		),
		edlisp.NewList(edlisp.NewSymbol("set-mark")),
		edlisp.NewList(
			edlisp.NewSymbol("search-forward"),
			edlisp.NewString("("),
		),
		edlisp.NewList(
			edlisp.NewSymbol("replace-region"),
			edlisp.NewString("helloWorld"),
		),
	}

	err := writer.Write(&buf, expressions)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 4 {
		t.Fatalf("Expected 4 lines, got %d", len(lines))
	}

	expectedLines := []string{
		"search-forward \"doIt\"",
		"set-mark",
		"search-forward \"(\"",
		"replace-region \"helloWorld\"",
	}

	for i, expected := range expectedLines {
		if lines[i] != expected {
			t.Errorf("Line %d: expected %q, got %q", i, expected, lines[i])
		}
	}
}

func TestShellWriter_WriteValue_NonListValue(t *testing.T) {
	writer := &ShellWriter{}
	var buf bytes.Buffer

	value := edlisp.NewSymbol("standalone-symbol")
	err := writer.WriteValue(&buf, value)
	if err == nil {
		t.Fatal("Expected error for non-list value, got none")
	}

	if !strings.Contains(err.Error(), "can only convert lists to shell format") {
		t.Errorf("Expected 'can only convert lists' error, got: %v", err)
	}
}

func TestShellWriter_WriteValue_NestedList(t *testing.T) {
	writer := &ShellWriter{}
	var buf bytes.Buffer

	value := edlisp.NewList(
		edlisp.NewSymbol("command"),
		edlisp.NewList(
			edlisp.NewSymbol("nested"),
			edlisp.NewString("arg"),
		),
	)

	err := writer.WriteValue(&buf, value)
	if err == nil {
		t.Fatal("Expected error for nested list, got none")
	}

	if !strings.Contains(err.Error(), "nested lists are not supported") {
		t.Errorf("Expected 'nested lists are not supported' error, got: %v", err)
	}
}

func TestShellWriter_valueToToken_UnsupportedType(t *testing.T) {
	writer := &ShellWriter{}

	_, err := writer.valueToToken(nil)
	if err == nil {
		t.Fatal("Expected error for nil value, got none")
	}

	if !strings.Contains(err.Error(), "unsupported value type") {
		t.Errorf("Expected 'unsupported value type' error, got: %v", err)
	}
}
