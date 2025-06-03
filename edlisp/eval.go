package edlisp

import (
	"fmt"
	"strings"
)

// BuiltinFn represents a built-in function that can be called from texted scripts.
type BuiltinFn func(args []Value, buffer *Buffer) (Value, error)

// Environment represents the execution environment for evaluation.
type Environment struct {
	// Functions maps function names to their implementations
	Functions map[string]BuiltinFn
}

// Buffer represents a text buffer for editing operations.
type Buffer struct {
	content strings.Builder
	point   int
	mark    int
}

// NewBuffer creates a new buffer with the given initial content.
func NewBuffer(content string) *Buffer {
	var b Buffer
	b.content.WriteString(content)
	b.point = 1
	b.mark = 1
	return &b
}

// String returns the current buffer content.
func (b *Buffer) String() string {
	return b.content.String()
}

// Point returns the current cursor position.
func (b *Buffer) Point() int {
	return b.point
}

// Mark returns the current mark position.
func (b *Buffer) Mark() int {
	return b.mark
}

// SetPoint sets the cursor position.
func (b *Buffer) SetPoint(pos int) {
	b.point = pos
}

// SetMark sets the mark position.
func (b *Buffer) SetMark(pos int) {
	b.mark = pos
}

// Insert inserts text at the current point.
func (b *Buffer) Insert(text string) {
	content := b.content.String()
	if b.point <= 1 {
		b.content.Reset()
		b.content.WriteString(text)
		b.content.WriteString(content)
	} else if b.point > len(content)+1 {
		b.content.WriteString(text)
	} else {
		newContent := content[:b.point-1] + text + content[b.point-1:]
		b.content.Reset()
		b.content.WriteString(newContent)
	}
	b.point += len(text)
}

// Eval executes a texted program in the given environment.
func Eval(program []Value, env *Environment, buffer *Buffer) (Value, error) {
	var result Value = NewString("")

	for _, expr := range program {
		val, err := evalExpression(expr, env, buffer)
		if err != nil {
			return nil, err
		}
		result = val
	}

	return result, nil
}

// evalExpression evaluates a single expression.
func evalExpression(expr Value, env *Environment, buffer *Buffer) (Value, error) {
	switch {
	case IsA(expr, TheStringKind):
		return expr, nil
	case IsA(expr, TheNumberKind):
		return expr, nil
	case IsA(expr, TheListKind):
		list := expr.(*List)
		if list.Len() == 0 {
			return nil, fmt.Errorf("empty list")
		}

		firstElem := list.Get(0)
		if !IsA(firstElem, TheSymbolKind) {
			return nil, fmt.Errorf("first element of list must be a symbol")
		}

		symbol := firstElem.(*Symbol)
		fnName := symbol.Name

		fn, exists := env.Functions[fnName]
		if !exists {
			return nil, fmt.Errorf("undefined-function %q", fnName)
		}

		args := make([]Value, list.Len()-1)
		for i := 1; i < list.Len(); i++ {
			args[i-1] = list.Get(i)
		}

		return fn(args, buffer)
	default:
		return nil, fmt.Errorf("unknown expression type")
	}
}

// NewDefaultEnvironment creates a default evaluation environment with basic functions.
func NewDefaultEnvironment() *Environment {
	env := &Environment{
		Functions: make(map[string]BuiltinFn),
	}

	// Add basic insert function
	env.Functions["insert"] = func(args []Value, buffer *Buffer) (Value, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("insert expects 1 argument, got %d", len(args))
		}

		if !IsA(args[0], TheStringKind) {
			return nil, fmt.Errorf("insert expects a string argument")
		}

		str := args[0].(*String)
		buffer.Insert(str.Value)
		return NewString(""), nil
	}

	return env
}