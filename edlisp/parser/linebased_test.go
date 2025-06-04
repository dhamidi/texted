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

func TestParseString_NestedSExpression_BufferSubstring(t *testing.T) {
	input := `buffer-substring (region-beginning) (region-end)`
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

	// Check first element is symbol 'buffer-substring'
	sym := list.Elements[0].(*edlisp.Symbol)
	if sym.Name != "buffer-substring" {
		t.Errorf("expected 'buffer-substring', got '%s'", sym.Name)
	}

	// Check second element is nested list (region-beginning)
	nestedList1 := list.Elements[1].(*edlisp.List)
	if len(nestedList1.Elements) != 1 {
		t.Fatalf("expected 1 element in first nested list, got %d", len(nestedList1.Elements))
	}
	if nestedList1.Elements[0].(*edlisp.Symbol).Name != "region-beginning" {
		t.Errorf("expected 'region-beginning' in first nested list")
	}

	// Check third element is nested list (region-end)
	nestedList2 := list.Elements[2].(*edlisp.List)
	if len(nestedList2.Elements) != 1 {
		t.Fatalf("expected 1 element in second nested list, got %d", len(nestedList2.Elements))
	}
	if nestedList2.Elements[0].(*edlisp.Symbol).Name != "region-end" {
		t.Errorf("expected 'region-end' in second nested list")
	}
}

func TestParseString_DeeplyNestedSExpression(t *testing.T) {
	input := `(outer (middle (inner "value")) 42)`
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

	// Check first element is symbol 'outer'
	sym := list.Elements[0].(*edlisp.Symbol)
	if sym.Name != "outer" {
		t.Errorf("expected 'outer', got '%s'", sym.Name)
	}

	// Check second element is nested list (middle ...)
	middleList := list.Elements[1].(*edlisp.List)
	if len(middleList.Elements) != 2 {
		t.Fatalf("expected 2 elements in middle list, got %d", len(middleList.Elements))
	}
	if middleList.Elements[0].(*edlisp.Symbol).Name != "middle" {
		t.Errorf("expected 'middle' in middle list")
	}

	// Check the inner nested list
	innerList := middleList.Elements[1].(*edlisp.List)
	if len(innerList.Elements) != 2 {
		t.Fatalf("expected 2 elements in inner list, got %d", len(innerList.Elements))
	}
	if innerList.Elements[0].(*edlisp.Symbol).Name != "inner" {
		t.Errorf("expected 'inner' in inner list")
	}
	if innerList.Elements[1].(*edlisp.String).Value != "value" {
		t.Errorf("expected 'value' string in inner list")
	}

	// Check third element is number 42
	num := list.Elements[2].(*edlisp.Number)
	if num.Value != 42 {
		t.Errorf("expected 42, got %f", num.Value)
	}
}

func TestParseString_MultipleNestedExpressions(t *testing.T) {
	input := `(replace-region (buffer-substring (region-beginning) (region-end)))`
	result, err := ParseString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 expression, got %d", len(result))
	}

	list := result[0].(*edlisp.List)
	if len(list.Elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(list.Elements))
	}

	// Check first element is symbol 'replace-region'
	sym := list.Elements[0].(*edlisp.Symbol)
	if sym.Name != "replace-region" {
		t.Errorf("expected 'replace-region', got '%s'", sym.Name)
	}

	// Check second element is nested buffer-substring expression
	bufferSubList := list.Elements[1].(*edlisp.List)
	if len(bufferSubList.Elements) != 3 {
		t.Fatalf("expected 3 elements in buffer-substring list, got %d", len(bufferSubList.Elements))
	}
	if bufferSubList.Elements[0].(*edlisp.Symbol).Name != "buffer-substring" {
		t.Errorf("expected 'buffer-substring' in nested list")
	}

	// Check the nested function calls within buffer-substring
	regionBeginList := bufferSubList.Elements[1].(*edlisp.List)
	if regionBeginList.Elements[0].(*edlisp.Symbol).Name != "region-beginning" {
		t.Errorf("expected 'region-beginning'")
	}

	regionEndList := bufferSubList.Elements[2].(*edlisp.List)
	if regionEndList.Elements[0].(*edlisp.Symbol).Name != "region-end" {
		t.Errorf("expected 'region-end'")
	}
}

func TestParseString_NestedWithMixedTypes(t *testing.T) {
	input := `(goto-char (+ (point) 10))`
	result, err := ParseString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 expression, got %d", len(result))
	}

	list := result[0].(*edlisp.List)
	if len(list.Elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(list.Elements))
	}

	// Check first element is symbol 'goto-char'
	sym := list.Elements[0].(*edlisp.Symbol)
	if sym.Name != "goto-char" {
		t.Errorf("expected 'goto-char', got '%s'", sym.Name)
	}

	// Check second element is nested arithmetic expression
	mathList := list.Elements[1].(*edlisp.List)
	if len(mathList.Elements) != 3 {
		t.Fatalf("expected 3 elements in math expression, got %d", len(mathList.Elements))
	}

	if mathList.Elements[0].(*edlisp.Symbol).Name != "+" {
		t.Errorf("expected '+' operator")
	}

	pointList := mathList.Elements[1].(*edlisp.List)
	if pointList.Elements[0].(*edlisp.Symbol).Name != "point" {
		t.Errorf("expected 'point' function call")
	}

	num := mathList.Elements[2].(*edlisp.Number)
	if num.Value != 10 {
		t.Errorf("expected 10, got %f", num.Value)
	}
}

func TestParseString_SemicolonSeparatedCommands(t *testing.T) {
	input := `goto-char 6; delete-char 2`
	result, err := ParseString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 expressions, got %d", len(result))
	}

	// Check first command: goto-char 6
	list1 := result[0].(*edlisp.List)
	if len(list1.Elements) != 2 {
		t.Fatalf("expected 2 elements in first command, got %d", len(list1.Elements))
	}

	sym1 := list1.Elements[0].(*edlisp.Symbol)
	if sym1.Name != "goto-char" {
		t.Errorf("expected 'goto-char', got '%s'", sym1.Name)
	}

	num1 := list1.Elements[1].(*edlisp.Number)
	if num1.Value != 6 {
		t.Errorf("expected 6, got %f", num1.Value)
	}

	// Check second command: delete-char 2
	list2 := result[1].(*edlisp.List)
	if len(list2.Elements) != 2 {
		t.Fatalf("expected 2 elements in second command, got %d", len(list2.Elements))
	}

	sym2 := list2.Elements[0].(*edlisp.Symbol)
	if sym2.Name != "delete-char" {
		t.Errorf("expected 'delete-char', got '%s'", sym2.Name)
	}

	num2 := list2.Elements[1].(*edlisp.Number)
	if num2.Value != 2 {
		t.Errorf("expected 2, got %f", num2.Value)
	}
}

func TestParseString_SemicolonWithSpaces(t *testing.T) {
	input := `set-mark ; goto-char 10 ; delete-char 1`
	result, err := ParseString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 expressions, got %d", len(result))
	}

	// Check commands are parsed correctly
	commands := []string{"set-mark", "goto-char", "delete-char"}
	for i, expectedCmd := range commands {
		list := result[i].(*edlisp.List)
		sym := list.Elements[0].(*edlisp.Symbol)
		if sym.Name != expectedCmd {
			t.Errorf("command %d: expected '%s', got '%s'", i, expectedCmd, sym.Name)
		}
	}
}

func TestParseString_SemicolonInQuotes(t *testing.T) {
	input := `search-forward "text;with;semicolons"; replace-match "new"`
	result, err := ParseString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 expressions, got %d", len(result))
	}

	// Check first command has string with semicolons
	list1 := result[0].(*edlisp.List)
	str1 := list1.Elements[1].(*edlisp.String)
	if str1.Value != "text;with;semicolons" {
		t.Errorf("expected 'text;with;semicolons', got '%s'", str1.Value)
	}

	// Check second command
	list2 := result[1].(*edlisp.List)
	sym2 := list2.Elements[0].(*edlisp.Symbol)
	if sym2.Name != "replace-match" {
		t.Errorf("expected 'replace-match', got '%s'", sym2.Name)
	}
}

func TestParseString_MultipleSemicolonsAndEmptyCommands(t *testing.T) {
	input := `goto-char 5;; delete-char 1; ; set-mark`
	result, err := ParseString(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Should parse 3 commands (empty commands are skipped)
	if len(result) != 3 {
		t.Fatalf("expected 3 expressions, got %d", len(result))
	}

	// Verify the commands
	commands := []string{"goto-char", "delete-char", "set-mark"}
	for i, expectedCmd := range commands {
		list := result[i].(*edlisp.List)
		sym := list.Elements[0].(*edlisp.Symbol)
		if sym.Name != expectedCmd {
			t.Errorf("command %d: expected '%s', got '%s'", i, expectedCmd, sym.Name)
		}
	}
}

func TestSplitOnSemicolons(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{`command1; command2`, []string{"command1", " command2"}},
		{`cmd "arg;with;semi"; cmd2`, []string{`cmd "arg;with;semi"`, " cmd2"}},
		{`cmd1;;cmd2`, []string{"cmd1", "", "cmd2"}},
		{`single-command`, []string{"single-command"}},
		{`cmd "escaped\"semicolon;here"; cmd2`, []string{`cmd "escaped\"semicolon;here"`, " cmd2"}},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := splitOnSemicolons(tc.input)
			if len(result) != len(tc.expected) {
				t.Errorf("expected %d parts, got %d", len(tc.expected), len(result))
			}
			for i, expected := range tc.expected {
				if i >= len(result) || result[i] != expected {
					t.Errorf("part %d: expected %q, got %q", i, expected, result[i])
				}
			}
		})
	}
}
