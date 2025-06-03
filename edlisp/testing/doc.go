// Package testing provides test execution facilities for texted scripts.
//
// This package can parse XML-based test cases and execute them using a configurable
// environment with built-in functions. The test format includes:
//
//   - <buffer>initial content</buffer> - Initial buffer content
//   - <input lang="shell">script</input> - Script to execute
//   - <output>expected output</output> - Expected buffer content after execution
//   - <error lang="sexp">error</error> - Expected error if any
//
// Example usage:
//
//	env := testing.NewDefaultEnvironment()
//
//	// Add custom functions
//	env.Functions["search-forward"] = func(args []edlisp.Value, buffer *edlisp.Buffer) (edlisp.Value, error) {
//		// Implementation here
//		return edlisp.NewString(""), nil
//	}
//
//	result := testing.RunTestFile("test.xml", env)
//	if !result.Passed {
//		fmt.Printf("Test failed: %v\n", result.Error)
//	}
//
// The package uses edlisp.Eval function that implements the evaluation
// semantics described in docs/spec.md, supporting:
//
//   - Symbol resolution to functions
//   - String and number literal evaluation
//   - List evaluation as function calls
//   - Basic error handling
package testing
