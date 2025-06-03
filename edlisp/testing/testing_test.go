package testing

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dhamidi/texted/edlisp"
)

func TestParseTestCase(t *testing.T) {
	xmlContent := `<buffer>initial content</buffer>
<input lang="shell">
command "arg1" "arg2"
</input>
<output>expected output</output>
<error lang="sexp">
</error>`

	testCase, err := ParseTestCase(strings.NewReader(xmlContent))
	if err != nil {
		t.Fatalf("Failed to parse test case: %v", err)
	}

	if testCase.Buffer != "initial content" {
		t.Errorf("Expected buffer 'initial content', got %q", testCase.Buffer)
	}

	if testCase.Input.Lang != "shell" {
		t.Errorf("Expected input lang 'shell', got %q", testCase.Input.Lang)
	}

	expectedInput := `command "arg1" "arg2"`
	if !strings.Contains(testCase.Input.Text, expectedInput) {
		t.Errorf("Expected input to contain %q, got %q", expectedInput, testCase.Input.Text)
	}

	if testCase.Output != "expected output" {
		t.Errorf("Expected output 'expected output', got %q", testCase.Output)
	}

	if testCase.Error.Lang != "sexp" {
		t.Errorf("Expected error lang 'sexp', got %q", testCase.Error.Lang)
	}
}

func TestParseTestCaseWithError(t *testing.T) {
	xmlContent := `<buffer></buffer>
<input lang="shell">
bad-function
</input>
<output></output>
<error lang="sexp">
(undefined-function "bad-function")
</error>`

	testCase, err := ParseTestCase(strings.NewReader(xmlContent))
	if err != nil {
		t.Fatalf("Failed to parse test case: %v", err)
	}

	expectedError := `(undefined-function "bad-function")`
	if !strings.Contains(testCase.Error.Text, expectedError) {
		t.Errorf("Expected error to contain %q, got %q", expectedError, testCase.Error.Text)
	}
}

func TestRunTestSuccess(t *testing.T) {
	// Create a mock environment that always succeeds
	env := &edlisp.Environment{
		Functions: map[string]edlisp.BuiltinFn{
			"mock-command": func(args []edlisp.Value, buffer *edlisp.Buffer) (edlisp.Value, error) {
				// Mock function that sets buffer to expected output
				buffer.SetPoint(1)
				buffer.Insert("success output")
				return edlisp.NewString(""), nil
			},
		},
	}

	testCase := &TestCase{
		Buffer: "",
		Input: Input{
			Lang: "shell",
			Text: `mock-command`,
		},
		Output: "success output",
		Error:  Error{Lang: "sexp", Text: ""},
	}

	result := RunTest(testCase, env)

	if !result.Passed {
		t.Errorf("Expected test to pass, but it failed: %v", result.Error)
	}

	if result.Error != nil {
		t.Errorf("Expected no error, got: %v", result.Error)
	}
}

func TestRunTestOutputMismatch(t *testing.T) {
	// Create a mock environment that produces different output
	env := &edlisp.Environment{
		Functions: map[string]edlisp.BuiltinFn{
			"mock-command": func(args []edlisp.Value, buffer *edlisp.Buffer) (edlisp.Value, error) {
				buffer.SetPoint(1)
				buffer.Insert("actual output")
				return edlisp.NewString(""), nil
			},
		},
	}

	testCase := &TestCase{
		Buffer: "",
		Input: Input{
			Lang: "shell",
			Text: `mock-command`,
		},
		Output: "expected output",
		Error:  Error{Lang: "sexp", Text: ""},
	}

	result := RunTest(testCase, env)

	if result.Passed {
		t.Error("Expected test to fail due to output mismatch")
	}

	if result.Error == nil {
		t.Error("Expected error to be set for failed test")
	}

	if result.Expected != "expected output" {
		t.Errorf("Expected result.Expected to be 'expected output', got %q", result.Expected)
	}

	if result.Actual != "actual output" {
		t.Errorf("Expected result.Actual to be 'actual output', got %q", result.Actual)
	}
}

func TestRunTestExpectedError(t *testing.T) {
	// Create a mock environment that produces an error
	env := &edlisp.Environment{
		Functions: map[string]edlisp.BuiltinFn{
			"error-command": func(args []edlisp.Value, buffer *edlisp.Buffer) (edlisp.Value, error) {
				return nil, fmt.Errorf("undefined-function \"error-command\"")
			},
		},
	}

	testCase := &TestCase{
		Buffer: "",
		Input: Input{
			Lang: "shell",
			Text: `error-command`,
		},
		Output: "",
		Error:  Error{Lang: "sexp", Text: `(undefined-function "error-command")`},
	}

	result := RunTest(testCase, env)

	if !result.Passed {
		t.Errorf("Expected test to pass (error was expected), but it failed: %v", result.Error)
	}
}

func TestRunTestUnexpectedError(t *testing.T) {
	// Create a mock environment that produces an unexpected error
	env := &edlisp.Environment{
		Functions: map[string]edlisp.BuiltinFn{
			"broken-command": func(args []edlisp.Value, buffer *edlisp.Buffer) (edlisp.Value, error) {
				return nil, fmt.Errorf("unexpected system error")
			},
		},
	}

	testCase := &TestCase{
		Buffer: "",
		Input: Input{
			Lang: "shell",
			Text: `broken-command`,
		},
		Output: "some output",
		Error:  Error{Lang: "sexp", Text: ""}, // No error expected
	}

	result := RunTest(testCase, env)

	if result.Passed {
		t.Error("Expected test to fail due to unexpected error")
	}

	if result.Error == nil {
		t.Error("Expected error to be set for failed test")
	}

	if !strings.Contains(result.Error.Error(), "unexpected error") {
		t.Errorf("Expected error message to mention unexpected error, got: %v", result.Error)
	}
}

func TestRunTestParseError(t *testing.T) {
	// Test what happens when input cannot be parsed
	env := NewDefaultEnvironment()

	testCase := &TestCase{
		Buffer: "",
		Input: Input{
			Lang: "shell",
			Text: `"unclosed string`,
		},
		Output: "",
		Error:  Error{Lang: "sexp", Text: ""},
	}

	result := RunTest(testCase, env)

	if result.Passed {
		t.Error("Expected test to fail due to parse error")
	}

	if result.Error == nil {
		t.Error("Expected error to be set for parse failure")
	}

	if !strings.Contains(result.Error.Error(), "parsing input") {
		t.Errorf("Expected error message to mention parsing, got: %v", result.Error)
	}
}

func TestRunTestFileNotFound(t *testing.T) {
	env := NewDefaultEnvironment()
	result := RunTestFile("nonexistent.xml", env)

	if result.Passed {
		t.Error("Expected test to fail for nonexistent file")
	}

	if result.Error == nil {
		t.Error("Expected error to be set for missing file")
	}

	if result.Name != "nonexistent.xml" {
		t.Errorf("Expected result name to be 'nonexistent.xml', got %q", result.Name)
	}
}
