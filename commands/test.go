package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dhamidi/texted/edlisp"
	"github.com/dhamidi/texted/edlisp/testing"
)

// NewTestCommand creates the test subcommand
func NewTestCommand() *cobra.Command {
	var verbose bool
	var quiet bool
	var failOnly bool
	var trace bool
	var include []string

	cmd := &cobra.Command{
		Use:   "test [TEST_FILES...]",
		Short: "Run XML-based tests",
		Long: `Run XML-based tests found in the tests/ directory.

Tests can be run individually by specifying test files, or filtered using include patterns.
Output verbosity can be controlled with flags.`,
		Run: func(cmd *cobra.Command, args []string) {
			runTests(args, include, verbose, quiet, failOnly, trace)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose output with before/after buffer states")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Show only pass/fail summary")
	cmd.Flags().BoolVar(&failOnly, "fail-only", false, "Show output only for failing tests")
	cmd.Flags().BoolVar(&trace, "trace", false, "Show trace output for each instruction")
	cmd.Flags().StringSliceVarP(&include, "include", "i", nil, "Include tests matching pattern (can be used multiple times)")

	return cmd
}

// runTests executes the test runner
func runTests(testFiles []string, includePatterns []string, verbose, quiet, failOnly, trace bool) {
	var filesToTest []string

	if len(testFiles) == 0 {
		// Find all XML files in tests/ directory
		matches, err := filepath.Glob("tests/*.xml")
		if err != nil {
			fmt.Printf("Error finding test files: %v\n", err)
			os.Exit(1)
		}
		filesToTest = matches
	} else {
		filesToTest = testFiles
	}

	// Filter by include patterns if specified
	if len(includePatterns) > 0 {
		filtered := make([]string, 0)
		for _, file := range filesToTest {
			basename := filepath.Base(file)
			for _, pattern := range includePatterns {
				matched, err := regexp.MatchString(pattern, basename)
				if err != nil {
					fmt.Printf("Error: invalid pattern %s: %v\n", pattern, err)
					os.Exit(1)
				}
				if matched {
					filtered = append(filtered, file)
					break
				}
			}
		}
		filesToTest = filtered
	}

	if len(filesToTest) == 0 {
		fmt.Println("No test files found")
		return
	}

	// Create test environment
	env := edlisp.NewDefaultEnvironment()

	var results []*testing.TestResult
	totalTests := 0
	passedTests := 0

	for _, testFile := range filesToTest {
		var result *testing.TestResult
		if trace {
			traceCallback := createTraceCallback()
			result = testing.RunTestFileWithTrace(testFile, env, traceCallback)
		} else {
			result = testing.RunTestFile(testFile, env)
		}
		results = append(results, result)
		totalTests++
		if result.Passed {
			passedTests++
		}
	}

	// Output results based on flags
	if !quiet {
		for _, result := range results {
			if failOnly && result.Passed {
				continue
			}
			printTestResult(result, verbose)
		}
	}

	// Print summary
	fmt.Printf("\nTests: %d passed, %d failed, %d total\n", passedTests, totalTests-passedTests, totalTests)

	if passedTests < totalTests {
		os.Exit(1)
	}
}

// printTestResult prints a single test result
func printTestResult(result *testing.TestResult, verbose bool) {
	status := "PASS"
	if !result.Passed {
		status = "FAIL"
	}

	fmt.Printf("[%s] %s\n", status, result.Name)

	if result.Error != nil {
		fmt.Printf("  Error: %s\n", result.Error)
		return
	}

	if !result.Passed || verbose {
		if verbose && result.Expected != "" && result.Actual != "" {
			fmt.Printf("  Expected: %q\n", result.Expected)
			fmt.Printf("  Actual:   %q\n", result.Actual)
		} else if !result.Passed {
			fmt.Printf("  Expected: %q\n", result.Expected)
			fmt.Printf("  Actual:   %q\n", result.Actual)
		}
	}
}

// createTraceCallback creates a trace callback function that logs buffer state.
func createTraceCallback() edlisp.TraceCallback {
	return func(ctx *edlisp.TraceContext) {
		fmt.Printf("TRACE: Executed instruction: %s\n", formatInstruction(ctx.Instruction))
		fmt.Printf("  Buffer content: %q\n", ctx.Buffer.String())
		fmt.Printf("  Point: %d, Mark: %d\n", ctx.Buffer.Point(), ctx.Buffer.Mark())
		fmt.Println()
	}
}

// formatInstruction formats an instruction for trace output.
func formatInstruction(instruction edlisp.Value) string {
	if instruction == nil {
		return "nil"
	}
	
	switch {
	case edlisp.IsA(instruction, edlisp.TheStringKind):
		str := instruction.(*edlisp.String)
		return fmt.Sprintf(`"%s"`, str.Value)
	case edlisp.IsA(instruction, edlisp.TheNumberKind):
		num := instruction.(*edlisp.Number)
		if num.Value == float64(int(num.Value)) {
			return fmt.Sprintf("%d", int(num.Value))
		}
		return fmt.Sprintf("%g", num.Value)
	case edlisp.IsA(instruction, edlisp.TheSymbolKind):
		sym := instruction.(*edlisp.Symbol)
		return sym.Name
	case edlisp.IsA(instruction, edlisp.TheListKind):
		list := instruction.(*edlisp.List)
		var parts []string
		for i := 0; i < list.Len(); i++ {
			parts = append(parts, formatInstruction(list.Get(i)))
		}
		return fmt.Sprintf("(%s)", strings.Join(parts, " "))
	default:
		return fmt.Sprintf("%v", instruction)
	}
}