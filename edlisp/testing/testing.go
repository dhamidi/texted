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
	Result Result `xml:"result"`
	Error  Error  `xml:"error"`
}

// Input represents the test input with language specification.
type Input struct {
	Lang string `xml:"lang,attr"`
	Text string `xml:",chardata"`
}

// Result represents expected result value with language specification.
type Result struct {
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
		Result  Result    `xml:"result"`
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
		Result: wrapper.Result,
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
	return RunTestWithTrace(testCase, env, nil)
}

// RunTestWithTrace executes a single test case with optional tracing.
func RunTestWithTrace(testCase *TestCase, env *edlisp.Environment, traceCallback edlisp.TraceCallback) *TestResult {
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
	evalResult, evalErr := edlisp.EvalWithTrace(program, env, buffer, traceCallback)

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

	// Check result if specified
	expectedResult := strings.TrimSpace(testCase.Result.Text)
	if expectedResult != "" {
		// Parse the expected result using the appropriate parser based on lang attribute
		var expectedValues []edlisp.Value
		var err error
		
		lang := testCase.Result.Lang
		if lang == "" {
			lang = "sexp" // Default to sexp
		}
		
		expectedValues, err = parser.ParseFormat(lang, expectedResult)
		if err != nil {
			result.Error = fmt.Errorf("parsing expected result with format %s: %w", lang, err)
			return result
		}
		
		// Get the expected value (should be a single expression)
		var expectedValue edlisp.Value
		if len(expectedValues) == 0 {
			expectedValue = nil
		} else if len(expectedValues) == 1 {
			expectedValue = expectedValues[0]
		} else {
			result.Error = fmt.Errorf("expected result should contain exactly one expression, got %d", len(expectedValues))
			return result
		}
		
		// Compare using Equal function
		if !edlisp.Equal(expectedValue, evalResult) {
			result.Expected = formatValue(expectedValue)
			result.Actual = formatValue(evalResult)
			result.Error = fmt.Errorf("result mismatch:\nexpected: %q\nactual: %q", result.Expected, result.Actual)
			return result
		}
	}

	result.Passed = true
	return result
}

// RunTestFile runs a test from a file.
func RunTestFile(filename string, env *edlisp.Environment) *TestResult {
	return RunTestFileWithTrace(filename, env, nil)
}

// RunTestFileWithTrace runs a test from a file with optional tracing.
func RunTestFileWithTrace(filename string, env *edlisp.Environment, traceCallback edlisp.TraceCallback) *TestResult {
	testCase, err := ParseTestFile(filename)
	if err != nil {
		return &TestResult{
			Name:  filepath.Base(filename),
			Error: err,
		}
	}

	result := RunTestWithTrace(testCase, env, traceCallback)
	result.Name = filepath.Base(filename)
	return result
}

// formatValue converts an edlisp.Value to its string representation for comparison.
func formatValue(value edlisp.Value) string {
	if value == nil {
		return "nil"
	}
	
	switch {
	case edlisp.IsA(value, edlisp.TheStringKind):
		str := value.(*edlisp.String)
		return fmt.Sprintf(`"%s"`, str.Value)
	case edlisp.IsA(value, edlisp.TheNumberKind):
		num := value.(*edlisp.Number)
		// Format as integer if it's a whole number
		if num.Value == float64(int(num.Value)) {
			return fmt.Sprintf("%d", int(num.Value))
		}
		return fmt.Sprintf("%g", num.Value)
	case edlisp.IsA(value, edlisp.TheSymbolKind):
		sym := value.(*edlisp.Symbol)
		return sym.Name
	case edlisp.IsA(value, edlisp.TheListKind):
		list := value.(*edlisp.List)
		var parts []string
		for i := 0; i < list.Len(); i++ {
			parts = append(parts, formatValue(list.Get(i)))
		}
		return fmt.Sprintf("(%s)", strings.Join(parts, " "))
	default:
		return fmt.Sprintf("%v", value)
	}
}


// NewDefaultEnvironment creates a default testing environment with basic functions.
func NewDefaultEnvironment() *edlisp.Environment {
	return edlisp.NewDefaultEnvironment()
}