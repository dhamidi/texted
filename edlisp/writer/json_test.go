package writer

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dhamidi/texted/edlisp"
)

func TestJSONWriter_WriteValue_Symbol(t *testing.T) {
	writer := &JSONWriter{}
	var buf bytes.Buffer

	value := edlisp.NewSymbol("search-forward")
	err := writer.WriteValue(&buf, value)
	if err != nil {
		t.Fatalf("WriteValue failed: %v", err)
	}

	expected := "\"search-forward\"\n"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestJSONWriter_WriteValue_String(t *testing.T) {
	writer := &JSONWriter{}
	var buf bytes.Buffer

	value := edlisp.NewString("hello world")
	err := writer.WriteValue(&buf, value)
	if err != nil {
		t.Fatalf("WriteValue failed: %v", err)
	}

	expected := "\"hello world\"\n"
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestJSONWriter_WriteValue_Number(t *testing.T) {
	tests := []struct {
		name     string
		value    *edlisp.Number
		expected string
	}{
		{"integer", edlisp.NewIntNumber(42), "42\n"},
		{"float", edlisp.NewNumber(3.14), "3.14\n"},
		{"zero", edlisp.NewIntNumber(0), "0\n"},
		{"negative", edlisp.NewIntNumber(-123), "-123\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &JSONWriter{}
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

func TestJSONWriter_WriteValue_List(t *testing.T) {
	tests := []struct {
		name     string
		value    *edlisp.List
		expected string
	}{
		{
			"empty list",
			edlisp.NewEmptyList(),
			"[]\n",
		},
		{
			"simple command",
			edlisp.NewList(
				edlisp.NewSymbol("search-forward"),
				edlisp.NewString("text"),
			),
			"[\"search-forward\",\"text\"]\n",
		},
		{
			"command with number",
			edlisp.NewList(
				edlisp.NewSymbol("move"),
				edlisp.NewIntNumber(5),
			),
			"[\"move\",5]\n",
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
			"[\"progn\",[\"search-forward\",\"text\"],[\"replace-match\",\"replacement\"]]\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &JSONWriter{}
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

func TestJSONWriter_Write_MultipleExpressions(t *testing.T) {
	writer := &JSONWriter{}
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
		"[\"search-forward\",\"doIt\"]",
		"[\"set-mark\"]",
		"[\"search-forward\",\"(\"]",
		"[\"replace-region\",\"helloWorld\"]",
	}

	for i, expected := range expectedLines {
		if lines[i] != expected {
			t.Errorf("Line %d: expected %q, got %q", i, expected, lines[i])
		}
	}
}

func TestJSONWriter_valueToJSON_ValidJSON(t *testing.T) {
	writer := &JSONWriter{}

	value := edlisp.NewList(
		edlisp.NewSymbol("search-forward"),
		edlisp.NewString("text with \"quotes\" and \n newlines"),
		edlisp.NewNumber(3.14159),
	)

	jsonValue, err := writer.valueToJSON(value)
	if err != nil {
		t.Fatalf("valueToJSON failed: %v", err)
	}

	jsonBytes, err := json.Marshal(jsonValue)
	if err != nil {
		t.Fatalf("Generated invalid JSON: %v", err)
	}

	var unmarshaled interface{}
	err = json.Unmarshal(jsonBytes, &unmarshaled)
	if err != nil {
		t.Fatalf("JSON roundtrip failed: %v", err)
	}
}

func TestJSONWriter_WriteValue_UnsupportedType(t *testing.T) {
	writer := &JSONWriter{}
	var buf bytes.Buffer

	err := writer.WriteValue(&buf, nil)
	if err == nil {
		t.Fatal("Expected error for nil value, got none")
	}

	if !strings.Contains(err.Error(), "unsupported value type") {
		t.Errorf("Expected 'unsupported value type' error, got: %v", err)
	}
}
