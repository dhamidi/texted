// Package edlisp provides documentation types for builtin functions.
// This file defines the core structures used to document function behavior,
// parameters, examples, and related information.
package edlisp

// FunctionDoc represents comprehensive documentation for a single function.
// It includes all information needed to understand and use the function.
type FunctionDoc struct {
	// Name is the function name as used in scripts (e.g., "search-forward")
	Name string

	// Summary is a one-line description of what the function does
	Summary string

	// Description provides detailed explanation of the function's behavior,
	// including side effects, error conditions, and usage notes
	Description string

	// Parameters documents each function parameter
	Parameters []ParameterDoc

	// Examples provides practical usage examples
	Examples []ExampleDoc

	// SeeAlso lists related functions that users might find useful
	SeeAlso []string

	// Category groups functions by purpose (e.g., "movement", "editing", "search")
	Category string
}

// ParameterDoc documents a single function parameter.
type ParameterDoc struct {
	// Name is the parameter name for documentation purposes
	Name string

	// Type describes the expected value type ("string", "number", etc.)
	Type string

	// Description explains what this parameter controls
	Description string

	// Optional indicates whether this parameter can be omitted
	Optional bool
}

// ExampleDoc provides a concrete usage example for a function.
type ExampleDoc struct {
	// Description explains what this example demonstrates
	Description string

	// Input is the command or script to execute
	Input string

	// Buffer shows the initial buffer state (empty string if not relevant)
	Buffer string

	// Output describes the expected result or buffer state after execution
	Output string
}
