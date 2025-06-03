// Package testing provides test execution facilities for texted scripts.
package testing

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhamidi/texted/edlisp"
	"github.com/dhamidi/texted/edlisp/parser"
)

// TestCase represents a single test case parsed from XML.
type TestCase struct {
	Buffer string `xml:"buffer"`
	Input  Input  `xml:"input"`
	Output string `xml:"output"`
	Error  Error  `xml:"error"`
}

// Input represents the test input with language specification.
type Input struct {
	Lang string `xml:"lang,attr"`
	Text string `xml:",chardata"`
}

// Error represents expected error output with language specification.
type Error struct {
	Lang string `xml:"lang,attr"`
	Text string `xml:",chardata"`
}

// TestResult represents the result of running a test case.
type TestResult struct {
	Name     string
	Passed   bool
	Error    error
	Expected string
	Actual   string
}


// ParseTestCase parses a test case from XML.
func ParseTestCase(r io.Reader) (*TestCase, error) {
	// Since the XML doesn't have a root wrapper, we need to wrap it
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading test case: %w", err)
	}
	
	wrappedXML := "<testcase>" + string(content) + "</testcase>"
	
	var wrapper struct {
		XMLName xml.Name  `xml:"testcase"`
		Buffer  string    `xml:"buffer"`
		Input   Input     `xml:"input"`
		Output  string    `xml:"output"`
		Error   Error     `xml:"error"`
	}
	
	err = xml.Unmarshal([]byte(wrappedXML), &wrapper)
	if err != nil {
		return nil, fmt.Errorf("parsing test case: %w", err)
	}
	
	return &TestCase{
		Buffer: wrapper.Buffer,
		Input:  wrapper.Input,
		Output: wrapper.Output,
		Error:  wrapper.Error,
	}, nil
}

// ParseTestFile parses a test case from a file.
func ParseTestFile(filename string) (*TestCase, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening test file %s: %w", filename, err)
	}
	defer file.Close()

	return ParseTestCase(file)
}

// RunTest executes a single test case.
func RunTest(testCase *TestCase, env *edlisp.Environment) *TestResult {
	name := filepath.Base("unknown")
	result := &TestResult{Name: name}

	// Create buffer with initial content
	buffer := edlisp.NewBuffer(testCase.Buffer)

	// Parse the input script
	var program []edlisp.Value
	var err error

	input := strings.TrimSpace(testCase.Input.Text)
	if input == "" {
		program = []edlisp.Value{}
	} else {
		program, err = parser.ParseString(input)
		if err != nil {
			result.Error = fmt.Errorf("parsing input: %w", err)
			return result
		}
	}

	// Execute the program
	_, evalErr := edlisp.Eval(program, env, buffer)

	// Check for expected error
	expectedError := strings.TrimSpace(testCase.Error.Text)
	if expectedError != "" {
		if evalErr == nil {
			result.Error = fmt.Errorf("expected error %q but got none", expectedError)
			return result
		}
		if !strings.Contains(evalErr.Error(), strings.Trim(expectedError, "()")) {
			result.Error = fmt.Errorf("expected error containing %q but got %q", expectedError, evalErr.Error())
			return result
		}
		result.Passed = true
		return result
	}

	// Check for unexpected error
	if evalErr != nil {
		result.Error = fmt.Errorf("unexpected error: %w", evalErr)
		return result
	}

	// Check output
	expected := testCase.Output
	actual := buffer.String()

	if expected != actual {
		result.Expected = expected
		result.Actual = actual
		result.Error = fmt.Errorf("output mismatch:\nexpected: %q\nactual: %q", expected, actual)
		return result
	}

	result.Passed = true
	return result
}

// RunTestFile runs a test from a file.
func RunTestFile(filename string, env *edlisp.Environment) *TestResult {
	testCase, err := ParseTestFile(filename)
	if err != nil {
		return &TestResult{
			Name:  filepath.Base(filename),
			Error: err,
		}
	}

	result := RunTest(testCase, env)
	result.Name = filepath.Base(filename)
	return result
}

// NewDefaultEnvironment creates a default testing environment with basic functions.
func NewDefaultEnvironment() *edlisp.Environment {
	return edlisp.NewDefaultEnvironment()
}