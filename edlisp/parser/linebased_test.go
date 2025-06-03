package parser

import (
	"strings"
	"testing"

	"github.com/dhamidi/texted/edlisp"
)

func TestParseString_EmptyInput(t *testing.T) {
	result, err := ParseString("")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d expressions", len(result))
	}
}

func TestParseString_ShellLikeCommand(t *testing.T) {
	input := `search-forward "doIt"`
	result, err := ParseString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 expression, got %d", len(result))
	}

	list, ok := result[0].(*edlisp.List)
	if !ok {
		t.Fatalf("expected List, got %T", result[0])
	}

	if len(list.Elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(list.Elements))
	}

	sym, ok := list.Elements[0].(*edlisp.Symbol)
	if !ok {
		t.Fatalf("expected Symbol, got %T", list.Elements[0])
	}
	if sym.Name != "search-forward" {
		t.Errorf("expected symbol 'search-forward', got '%s'", sym.Name)
	}

	str, ok := list.Elements[1].(*edlisp.String)
	if !ok {
		t.Fatalf("expected String, got %T", list.Elements[1])
	}
	if str.Value != "doIt" {
		t.Errorf("expected string 'doIt', got '%s'", str.Value)
	}
}

func TestParseString_SExpression(t *testing.T) {
	input := `(search-forward "doIt")`
	result, err := ParseString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 expression, got %d", len(result))
	}

	list, ok := result[0].(*edlisp.List)
	if !ok {
		t.Fatalf("expected List, got %T", result[0])
	}

	if len(list.Elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(list.Elements))
	}

	sym, ok := list.Elements[0].(*edlisp.Symbol)
	if !ok {
		t.Fatalf("expected Symbol, got %T", list.Elements[0])
	}
	if sym.Name != "search-forward" {
		t.Errorf("expected symbol 'search-forward', got '%s'", sym.Name)
	}

	str, ok := list.Elements[1].(*edlisp.String)
	if !ok {
		t.Fatalf("expected String, got %T", list.Elements[1])
	}
	if str.Value != "doIt" {
		t.Errorf("expected string 'doIt', got '%s'", str.Value)
	}
}

func TestParseString_MultipleLines(t *testing.T) {
	input := `search-forward "doIt"
set-mark
search-forward "("
replace-region "helloWorld"`

	result, err := ParseString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 4 {
		t.Fatalf("expected 4 expressions, got %d", len(result))
	}

	// Check first command
	list := result[0].(*edlisp.List)
	if list.Elements[0].(*edlisp.Symbol).Name != "search-forward" {
		t.Errorf("expected 'search-forward', got '%s'", list.Elements[0].(*edlisp.Symbol).Name)
	}

	// Check second command (no arguments)
	list = result[1].(*edlisp.List)
	if len(list.Elements) != 1 {
		t.Errorf("expected 1 element for set-mark, got %d", len(list.Elements))
	}
	if list.Elements[0].(*edlisp.Symbol).Name != "set-mark" {
		t.Errorf("expected 'set-mark', got '%s'", list.Elements[0].(*edlisp.Symbol).Name)
	}
}

func TestParseString_Numbers(t *testing.T) {
	input := `move-point 42`
	result, err := ParseString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	list := result[0].(*edlisp.List)
	num, ok := list.Elements[1].(*edlisp.Number)
	if !ok {
		t.Fatalf("expected Number, got %T", list.Elements[1])
	}
	if num.Value != 42 {
		t.Errorf("expected 42, got %f", num.Value)
	}
}

func TestParseString_FloatingPointNumbers(t *testing.T) {
	input := `set-value 3.14`
	result, err := ParseString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	list := result[0].(*edlisp.List)
	num, ok := list.Elements[1].(*edlisp.Number)
	if !ok {
		t.Fatalf("expected Number, got %T", list.Elements[1])
	}
	if num.Value != 3.14 {
		t.Errorf("expected 3.14, got %f", num.Value)
	}
}

func TestParseString_NestedSExpression(t *testing.T) {
	input := `(progn (search-forward "test") (set-mark))`
	result, err := ParseString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 expression, got %d", len(result))
	}

	list := result[0].(*edlisp.List)
	if len(list.Elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(list.Elements))
	}

	// Check first element is symbol 'progn'
	sym := list.Elements[0].(*edlisp.Symbol)
	if sym.Name != "progn" {
		t.Errorf("expected 'progn', got '%s'", sym.Name)
	}

	// Check second element is nested list
	nestedList := list.Elements[1].(*edlisp.List)
	if nestedList.Elements[0].(*edlisp.Symbol).Name != "search-forward" {
		t.Errorf("expected 'search-forward' in nested list")
	}
}

func TestParseString_EscapedQuotes(t *testing.T) {
	input := `search-forward "say \"hello\""`
	result, err := ParseString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if len(result) == 0 {
		t.Fatal("expected at least one result")
	}

	list, ok := result[0].(*edlisp.List)
	if !ok {
		t.Fatalf("expected List, got %T", result[0])
	}

	if len(list.Elements) < 2 {
		t.Fatalf("expected at least 2 elements, got %d", len(list.Elements))
	}

	str, ok := list.Elements[1].(*edlisp.String)
	if !ok {
		t.Fatalf("expected String, got %T", list.Elements[1])
	}

	expected := `say "hello"`
	if str.Value != expected {
		t.Errorf("expected %q, got %q", expected, str.Value)
	}
}

func TestParseString_LeadingWhitespace(t *testing.T) {
	input := `   search-forward "test"
	set-mark`

	result, err := ParseString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 expressions, got %d", len(result))
	}
}

func TestParseString_EmptyLines(t *testing.T) {
	input := `search-forward "test"

set-mark

replace-match "hello"`

	result, err := ParseString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 expressions, got %d", len(result))
	}
}

func TestParseString_ErrorCases(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"unterminated string", `search-forward "unterminated`},
		{"unterminated list", `(search-forward "test"`},
		{"unexpected closing paren", `search-forward "test")`},
		{"invalid string escape", `search-forward "invalid\x"`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseString(tc.input)
			if err == nil {
				t.Errorf("expected error for input: %q", tc.input)
			}
		})
	}
}

func TestParseReader(t *testing.T) {
	input := `search-forward "test"
(set-mark)
replace-match "hello"`

	reader := strings.NewReader(input)
	result, err := ParseReader(reader)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 expressions, got %d", len(result))
	}
}

func TestTokenize(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{`search-forward "test"`, []string{"search-forward", `"test"`}},
		{`(set-mark)`, []string{"(", "set-mark", ")"}},
		{`func arg1 arg2`, []string{"func", "arg1", "arg2"}},
		{`"quoted string"`, []string{`"quoted string"`}},
		{`"string with spaces"`, []string{`"string with spaces"`}},
		{`(nested (list))`, []string{"(", "nested", "(", "list", ")", ")"}},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := tokenize(tc.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if len(result) != len(tc.expected) {
				t.Errorf("expected %d tokens, got %d", len(tc.expected), len(result))
			}
			for i, expected := range tc.expected {
				if i >= len(result) || result[i] != expected {
					t.Errorf("token %d: expected %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	testCases := []struct {
		input    string
		expected interface{}
	}{
		{`"hello"`, "hello"},
		{`"hello world"`, "hello world"},
		{`42`, 42.0},
		{`3.14`, 3.14},
		{`symbol`, "symbol"},
		{`search-forward`, "search-forward"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := parseToken(tc.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			switch expected := tc.expected.(type) {
			case string:
				switch v := result.(type) {
				case *edlisp.String:
					if v.Value != expected {
						t.Errorf("expected string %q, got %q", expected, v.Value)
					}
				case *edlisp.Symbol:
					if v.Name != expected {
						t.Errorf("expected symbol %q, got %q", expected, v.Name)
					}
				default:
					t.Errorf("expected string or symbol, got %T", result)
				}
			case float64:
				if num, ok := result.(*edlisp.Number); ok {
					if num.Value != expected {
						t.Errorf("expected number %f, got %f", expected, num.Value)
					}
				} else {
					t.Errorf("expected number, got %T", result)
				}
			}
		})
	}
}
