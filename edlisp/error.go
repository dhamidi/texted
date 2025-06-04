package edlisp

import (
	"fmt"
	"strings"
)

// ExecutionError captures the full execution state at the time of an error occurrence.
type ExecutionError struct {
	// OriginalError is the underlying error that occurred
	OriginalError error

	// Program is the complete parsed program being executed
	Program []Value

	// CurrentInstruction is the instruction that caused the error
	CurrentInstruction Value

	// InstructionIndex is the index of the current instruction in the program
	InstructionIndex int

	// BufferContents is a snapshot of the buffer contents when the error occurred
	BufferContents string

	// Point is the cursor position when the error occurred
	Point int

	// Mark is the mark position when the error occurred
	Mark int

	// LastSearchMatch stores the last search match string
	LastSearchMatch string

	// LastSearchStart stores the start position of the last search match
	LastSearchStart int

	// LastSearchEnd stores the end position of the last search match
	LastSearchEnd int

	// Environment contains the function registry at the time of error
	Environment *Environment
}

// Error implements the error interface.
func (e *ExecutionError) Error() string {
	var b strings.Builder

	// Start with the original error message
	b.WriteString(e.OriginalError.Error())

	// Add context information
	b.WriteString(" (at instruction ")
	b.WriteString(fmt.Sprintf("%d", e.InstructionIndex))
	b.WriteString(": ")

	// Add instruction details using built-in String() method
	if e.CurrentInstruction != nil {
		b.WriteString(fmt.Sprintf("%v", e.CurrentInstruction))
	} else {
		b.WriteString("(nil instruction)")
	}

	b.WriteString(", point=")
	b.WriteString(fmt.Sprintf("%d", e.Point))
	b.WriteString(", mark=")
	b.WriteString(fmt.Sprintf("%d", e.Mark))
	b.WriteString(")")

	return b.String()
}

// Unwrap returns the original error for error chain compatibility.
func (e *ExecutionError) Unwrap() error {
	return e.OriginalError
}

// NewExecutionError creates a new ExecutionError with the current execution state.
func NewExecutionError(originalError error, program []Value, instructionIndex int, currentInstruction Value, buffer *Buffer, env *Environment) *ExecutionError {
	return &ExecutionError{
		OriginalError:      originalError,
		Program:            program,
		CurrentInstruction: currentInstruction,
		InstructionIndex:   instructionIndex,
		BufferContents:     buffer.String(),
		Point:              buffer.Point(),
		Mark:               buffer.Mark(),
		LastSearchMatch:    buffer.lastSearchMatch,
		LastSearchStart:    buffer.lastSearchStart,
		LastSearchEnd:      buffer.lastSearchEnd,
		Environment:        env,
	}
}
