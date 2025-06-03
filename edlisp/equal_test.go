package edlisp

import "testing"

func TestEqual_NilValues(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Value
		expected bool
	}{
		{"both nil", nil, nil, true},
		{"first nil", nil, NewString("test"), false},
		{"second nil", NewString("test"), nil, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Equal(test.a, test.b)
			if result != test.expected {
				t.Errorf("Equal(%v, %v) = %v, want %v", test.a, test.b, result, test.expected)
			}
		})
	}
}

func TestEqual_Strings(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Value
		expected bool
	}{
		{"equal strings", NewString("hello"), NewString("hello"), true},
		{"different strings", NewString("hello"), NewString("world"), false},
		{"empty strings", NewString(""), NewString(""), true},
		{"empty vs non-empty", NewString(""), NewString("hello"), false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Equal(test.a, test.b)
			if result != test.expected {
				t.Errorf("Equal(%v, %v) = %v, want %v", test.a, test.b, result, test.expected)
			}
		})
	}
}

func TestEqual_Numbers(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Value
		expected bool
	}{
		{"equal integers", NewNumber(42), NewNumber(42), true},
		{"equal floats", NewNumber(3.14), NewNumber(3.14), true},
		{"different numbers", NewNumber(42), NewNumber(43), false},
		{"int vs float same value", NewIntNumber(5), NewNumber(5.0), true},
		{"zero values", NewNumber(0), NewNumber(0), true},
		{"negative numbers", NewNumber(-10), NewNumber(-10), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Equal(test.a, test.b)
			if result != test.expected {
				t.Errorf("Equal(%v, %v) = %v, want %v", test.a, test.b, result, test.expected)
			}
		})
	}
}

func TestEqual_Symbols(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Value
		expected bool
	}{
		{"equal symbols", NewSymbol("test"), NewSymbol("test"), true},
		{"different symbols", NewSymbol("foo"), NewSymbol("bar"), false},
		{"empty symbol names", NewSymbol(""), NewSymbol(""), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Equal(test.a, test.b)
			if result != test.expected {
				t.Errorf("Equal(%v, %v) = %v, want %v", test.a, test.b, result, test.expected)
			}
		})
	}
}

func TestEqual_Lists(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Value
		expected bool
	}{
		{
			"empty lists",
			NewEmptyList(),
			NewEmptyList(),
			true,
		},
		{
			"equal single element lists",
			NewList(NewString("hello")),
			NewList(NewString("hello")),
			true,
		},
		{
			"different single element lists",
			NewList(NewString("hello")),
			NewList(NewString("world")),
			false,
		},
		{
			"equal multi-element lists",
			NewList(NewString("hello"), NewNumber(42), NewSymbol("test")),
			NewList(NewString("hello"), NewNumber(42), NewSymbol("test")),
			true,
		},
		{
			"different length lists",
			NewList(NewString("hello")),
			NewList(NewString("hello"), NewNumber(42)),
			false,
		},
		{
			"same length different elements",
			NewList(NewString("hello"), NewNumber(42)),
			NewList(NewString("hello"), NewNumber(43)),
			false,
		},
		{
			"nested lists equal",
			NewList(NewString("outer"), NewList(NewString("inner"))),
			NewList(NewString("outer"), NewList(NewString("inner"))),
			true,
		},
		{
			"nested lists different",
			NewList(NewString("outer"), NewList(NewString("inner1"))),
			NewList(NewString("outer"), NewList(NewString("inner2"))),
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Equal(test.a, test.b)
			if result != test.expected {
				t.Errorf("Equal(%v, %v) = %v, want %v", test.a, test.b, result, test.expected)
			}
		})
	}
}

func TestEqual_DifferentKinds(t *testing.T) {
	tests := []struct {
		name string
		a, b Value
	}{
		{"string vs number", NewString("42"), NewNumber(42)},
		{"string vs symbol", NewString("test"), NewSymbol("test")},
		{"number vs symbol", NewNumber(42), NewSymbol("42")},
		{"list vs string", NewList(NewString("hello")), NewString("hello")},
		{"list vs number", NewList(NewNumber(42)), NewNumber(42)},
		{"list vs symbol", NewList(NewSymbol("test")), NewSymbol("test")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Equal(test.a, test.b)
			if result {
				t.Errorf("Equal(%v, %v) = true, want false (different kinds should not be equal)", test.a, test.b)
			}
		})
	}
}

func TestEqual_ComplexNestedStructures(t *testing.T) {
	// Test deeply nested lists
	deepList1 := NewList(
		NewString("level1"),
		NewList(
			NewString("level2"),
			NewList(
				NewString("level3"),
				NewNumber(42),
			),
		),
	)
	
	deepList2 := NewList(
		NewString("level1"),
		NewList(
			NewString("level2"),
			NewList(
				NewString("level3"),
				NewNumber(42),
			),
		),
	)
	
	deepList3 := NewList(
		NewString("level1"),
		NewList(
			NewString("level2"),
			NewList(
				NewString("level3"),
				NewNumber(43), // Different number
			),
		),
	)

	if !Equal(deepList1, deepList2) {
		t.Error("Expected deeply nested identical lists to be equal")
	}

	if Equal(deepList1, deepList3) {
		t.Error("Expected deeply nested lists with different leaf values to be unequal")
	}
}