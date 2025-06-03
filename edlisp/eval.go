package edlisp

import (
	"fmt"
	"strings"
)

// BuiltinFn represents a built-in function that can be called from texted scripts.
type BuiltinFn func(args []Value, buffer *Buffer) (Value, error)

// TraceContext holds the context information for trace callbacks.
type TraceContext struct {
	Buffer      *Buffer
	Environment *Environment
	Instruction Value
}

// TraceCallback is a function that is called after each instruction is executed.
type TraceCallback func(ctx *TraceContext)

// Environment represents the execution environment for evaluation.
type Environment struct {
	// Functions maps function names to their implementations
	Functions map[string]BuiltinFn
}

// Buffer represents a text buffer for editing operations.
type Buffer struct {
	content         strings.Builder
	point           int
	mark            int
	lastSearchMatch string
	lastSearchStart int
	lastSearchEnd   int
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
	return EvalWithTrace(program, env, buffer, nil)
}

// EvalWithTrace executes a texted program with optional tracing.
func EvalWithTrace(program []Value, env *Environment, buffer *Buffer, traceCallback TraceCallback) (Value, error) {
	var result Value = NewString("")

	for _, expr := range program {
		val, err := evalExpression(expr, env, buffer)
		if err != nil {
			return nil, err
		}
		result = val

		if traceCallback != nil {
			traceCallback(&TraceContext{
				Buffer:      buffer,
				Environment: env,
				Instruction: expr,
			})
		}
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
			evaluatedArg, err := evalExpression(list.Get(i), env, buffer)
			if err != nil {
				return nil, err
			}
			args[i-1] = evaluatedArg
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

	// Add builtin functions
	env.Functions["insert"] = BuiltinInsert
	env.Functions["goto-char"] = BuiltinGotoChar
	env.Functions["goto-line"] = BuiltinGotoLine
	env.Functions["search-forward"] = BuiltinSearchForward
	env.Functions["search-backward"] = BuiltinSearchBackward
	env.Functions["re-search-forward"] = BuiltinReSearchForward
	env.Functions["end-of-buffer"] = BuiltinEndOfBuffer
	env.Functions["set-mark"] = BuiltinSetMark
	env.Functions["mark"] = BuiltinMark
	env.Functions["point"] = BuiltinPoint
	env.Functions["line-number-at-pos"] = BuiltinLineNumberAtPos
	env.Functions["replace-match"] = BuiltinReplaceMatch
	env.Functions["delete-region"] = BuiltinDeleteRegion
	env.Functions["buffer-substring"] = BuiltinBufferSubstring
	env.Functions["mark-word"] = BuiltinMarkWord
	env.Functions["replace-region"] = BuiltinReplaceRegion
	env.Functions["region-beginning"] = BuiltinRegionBeginning
	env.Functions["region-end"] = BuiltinRegionEnd
	env.Functions["forward-char"] = BuiltinForwardChar
	env.Functions["length"] = BuiltinLength
	env.Functions["substring"] = BuiltinSubstring
	env.Functions["string-match"] = BuiltinStringMatch
	env.Functions["backward-char"] = BuiltinBackwardChar
	env.Functions["beginning-of-buffer"] = BuiltinBeginningOfBuffer
	env.Functions["beginning-of-line"] = BuiltinBeginningOfLine
	env.Functions["end-of-line"] = BuiltinEndOfLine
	env.Functions["buffer-size"] = BuiltinBufferSize
	env.Functions["point-max"] = BuiltinPointMax
	env.Functions["point-min"] = BuiltinPointMin
	env.Functions["current-column"] = BuiltinCurrentColumn
	env.Functions["upcase"] = BuiltinUpcase
	env.Functions["downcase"] = BuiltinDowncase
	env.Functions["capitalize"] = BuiltinCapitalize
	env.Functions["concat"] = BuiltinConcat
	env.Functions["set-mark-command"] = BuiltinSetMarkCommand
	env.Functions["exchange-point-and-mark"] = BuiltinExchangePointAndMark
	env.Functions["mark-whole-buffer"] = BuiltinMarkWholeBuffer
	env.Functions["delete-char"] = BuiltinDeleteChar
	env.Functions["delete-backward-char"] = BuiltinDeleteBackwardChar
	env.Functions["forward-word"] = BuiltinForwardWord
	env.Functions["backward-word"] = BuiltinBackwardWord
	env.Functions["looking-at"] = BuiltinLookingAt
	env.Functions["looking-back"] = BuiltinLookingBack
	env.Functions["mark-line"] = BuiltinMarkLine
	env.Functions["delete-line"] = BuiltinDeleteLine
	env.Functions["kill-line"] = BuiltinKillLine
	env.Functions["kill-word"] = BuiltinKillWord
	env.Functions["backward-kill-word"] = BuiltinBackwardKillWord
	env.Functions["re-search-backward"] = BuiltinReSearchBackward
	env.Functions["replace-regexp-in-string"] = BuiltinReplaceRegexpInString

	return env
}
