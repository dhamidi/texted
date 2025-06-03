package edlisp

import (
	"testing"
)

func TestBuffer(t *testing.T) {
	buffer := NewBuffer("initial")

	if buffer.String() != "initial" {
		t.Errorf("Expected 'initial', got %q", buffer.String())
	}

	// Insert at point 1 (beginning)
	buffer.Insert(" text")
	if buffer.String() != " textinitial" {
		t.Errorf("Expected ' textinitial', got %q", buffer.String())
	}

	// After insert, point moved to 6 (after " text")
	// Setting point to 11 should be after "t" in " textinitial"
	buffer.SetPoint(11)
	buffer.Insert("added")
	expected := " textiniti" + "added" + "al"
	if buffer.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buffer.String())
	}
}

func TestEvalString(t *testing.T) {
	env := NewDefaultEnvironment()
	buffer := NewBuffer("")

	program := []Value{NewString("hello")}
	result, err := Eval(program, env, buffer)

	if err != nil {
		t.Fatalf("Eval failed: %v", err)
	}

	if !IsA(result, TheStringKind) {
		t.Errorf("Expected string result, got %T", result)
	}

	str := result.(*String)
	if str.Value != "hello" {
		t.Errorf("Expected 'hello', got %q", str.Value)
	}
}

func TestEvalNumber(t *testing.T) {
	env := NewDefaultEnvironment()
	buffer := NewBuffer("")

	program := []Value{NewNumber(42)}
	result, err := Eval(program, env, buffer)

	if err != nil {
		t.Fatalf("Eval failed: %v", err)
	}

	if !IsA(result, TheNumberKind) {
		t.Errorf("Expected number result, got %T", result)
	}

	num := result.(*Number)
	if num.Value != 42 {
		t.Errorf("Expected 42, got %f", num.Value)
	}
}

func TestEvalInsertFunction(t *testing.T) {
	env := NewDefaultEnvironment()
	buffer := NewBuffer("")

	// Test the insert function
	insertCall := NewList(
		NewSymbol("insert"),
		NewString("hello, world"),
	)

	program := []Value{insertCall}
	_, err := Eval(program, env, buffer)

	if err != nil {
		t.Fatalf("Eval failed: %v", err)
	}

	if buffer.String() != "hello, world" {
		t.Errorf("Expected 'hello, world', got %q", buffer.String())
	}
}

func TestEvalUndefinedFunction(t *testing.T) {
	env := NewDefaultEnvironment()
	buffer := NewBuffer("")

	// Test undefined function
	undefinedCall := NewList(
		NewSymbol("undefined-function"),
	)

	program := []Value{undefinedCall}
	_, err := Eval(program, env, buffer)

	if err == nil {
		t.Error("Expected error for undefined function")
	}

	if err.Error() != `undefined-function "undefined-function"` {
		t.Errorf("Expected specific error message, got %q", err.Error())
	}
}

func TestEvalRecursiveArguments(t *testing.T) {
	env := NewDefaultEnvironment()
	buffer := NewBuffer("hello world")

	// Set buffer state: point at 1, mark at 6 (selecting "hello")
	buffer.SetPoint(1)
	buffer.SetMark(6)

	// Test buffer-substring with function calls as arguments
	// This should evaluate (region-beginning) and (region-end) first
	substringCall := NewList(
		NewSymbol("buffer-substring"),
		NewList(NewSymbol("region-beginning")),
		NewList(NewSymbol("region-end")),
	)

	program := []Value{substringCall}
	result, err := Eval(program, env, buffer)

	if err != nil {
		t.Fatalf("Eval failed: %v", err)
	}

	if !IsA(result, TheStringKind) {
		t.Errorf("Expected string result, got %T", result)
	}

	str := result.(*String)
	if str.Value != "hello" {
		t.Errorf("Expected 'hello', got %q", str.Value)
	}
}

func TestEvalNestedFunctionCalls(t *testing.T) {
	env := NewDefaultEnvironment()
	buffer := NewBuffer("test")

	// Test nested function calls: (length (concat "a" "b"))
	concatCall := NewList(
		NewSymbol("concat"),
		NewString("a"),
		NewString("b"),
	)

	lengthCall := NewList(
		NewSymbol("length"),
		concatCall,
	)

	program := []Value{lengthCall}
	result, err := Eval(program, env, buffer)

	if err != nil {
		t.Fatalf("Eval failed: %v", err)
	}

	if !IsA(result, TheNumberKind) {
		t.Errorf("Expected number result, got %T", result)
	}

	num := result.(*Number)
	if num.Value != 2 {
		t.Errorf("Expected 2, got %f", num.Value)
	}
}
