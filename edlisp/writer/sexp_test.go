package writer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dhamidi/texted/edlisp"
)

func TestSExpWriter_WriteValue_Symbol(t *testing.T) {
	writer := &SExpWriter{}
	var buf bytes.Buffer

	value := edlisp.NewSymbol("search-forward")
	err := writer.WriteValue(&buf, value)
	if err != nil {
		t.Fatalf("WriteValue failed: %v", err)
	}

	expected := "search-forward"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestSExpWriter_WriteValue_String(t *testing.T) {
	writer := &SExpWriter{}
	var buf bytes.Buffer

	value := edlisp.NewString("hello world")
	err := writer.WriteValue(&buf, value)
	if err != nil {
		t.Fatalf("WriteValue failed: %v", err)
	}

	expected := "\"hello world\""
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestSExpWriter_WriteValue_Number(t *testing.T) {
	tests := []struct {
		name     string
		value    *edlisp.Number
		expected string
	}{
		{"integer", edlisp.NewIntNumber(42), "42"},
		{"float", edlisp.NewNumber(3.14), "3.14"},
		{"zero", edlisp.NewIntNumber(0), "0"},
		{"negative", edlisp.NewIntNumber(-123), "-123"},
		{"large float", edlisp.NewNumber(1234.5678), "1234.5678"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &SExpWriter{}
			var buf bytes.Buffer

			err := writer.WriteValue(&buf, tt.value)
			if err != nil {
				t.Fatalf("WriteValue failed: %v", err)
			}

			if buf.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestSExpWriter_WriteValue_List(t *testing.T) {
	tests := []struct {
		name     string
		value    *edlisp.List
		expected string
	}{
		{
			"empty list",
			edlisp.NewEmptyList(),
			"()",
		},
		{
			"simple command",
			edlisp.NewList(
				edlisp.NewSymbol("search-forward"),
				edlisp.NewString("text"),
			),
			"(search-forward \"text\")",
		},
		{
			"command with number",
			edlisp.NewList(
				edlisp.NewSymbol("move"),
				edlisp.NewIntNumber(5),
			),
			"(move 5)",
		},
		{
			"nested list",
			edlisp.NewList(
				edlisp.NewSymbol("progn"),
				edlisp.NewList(
					edlisp.NewSymbol("search-forward"),
					edlisp.NewString("text"),
				),
				edlisp.NewList(
					edlisp.NewSymbol("replace-match"),
					edlisp.NewString("replacement"),
				),
			),
			"(progn (search-forward \"text\") (replace-match \"replacement\"))",
		},
		{
			"single symbol list",
			edlisp.NewList(edlisp.NewSymbol("set-mark")),
			"(set-mark)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &SExpWriter{}
			var buf bytes.Buffer

			err := writer.WriteValue(&buf, tt.value)
			if err != nil {
				t.Fatalf("WriteValue failed: %v", err)
			}

			if buf.String() != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, buf.String())
			}
		})
	}
}

func TestSExpWriter_Write_MultipleExpressions(t *testing.T) {
	writer := &SExpWriter{}
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
		"(search-forward \"doIt\")",
		"(set-mark)",
		"(search-forward \"(\")",
		"(replace-region \"helloWorld\")",
	}

	for i, expected := range expectedLines {
		if lines[i] != expected {
			t.Errorf("Line %d: expected %q, got %q", i, expected, lines[i])
		}
	}
}

func TestSExpWriter_WriteValue_StringWithSpecialCharacters(t *testing.T) {
	writer := &SExpWriter{}
	var buf bytes.Buffer

	value := edlisp.NewString("text with \"quotes\" and \n newlines")
	err := writer.WriteValue(&buf, value)
	if err != nil {
		t.Fatalf("WriteValue failed: %v", err)
	}

	expected := "\"text with \\\"quotes\\\" and \\n newlines\""
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestSExpWriter_WriteValue_ComplexNesting(t *testing.T) {
	writer := &SExpWriter{}
	var buf bytes.Buffer

	value := edlisp.NewList(
		edlisp.NewSymbol("if"),
		edlisp.NewList(
			edlisp.NewSymbol("search-forward"),
			edlisp.NewString("pattern"),
		),
		edlisp.NewList(
			edlisp.NewSymbol("progn"),
			edlisp.NewList(
				edlisp.NewSymbol("set-mark"),
			),
			edlisp.NewList(
				edlisp.NewSymbol("replace-match"),
				edlisp.NewString("replacement"),
			),
		),
	)

	err := writer.WriteValue(&buf, value)
	if err != nil {
		t.Fatalf("WriteValue failed: %v", err)
	}

	expected := "(if (search-forward \"pattern\") (progn (set-mark) (replace-match \"replacement\")))"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}
