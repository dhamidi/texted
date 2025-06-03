package commands

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/dhamidi/texted/edlisp"
	"github.com/dhamidi/texted/edlisp/testing"
)

// NewTestCommand creates the test subcommand
func NewTestCommand() *cobra.Command {
	var verbose bool
	var quiet bool
	var failOnly bool
	var include []string

	cmd := &cobra.Command{
		Use:   "test [TEST_FILES...]",
		Short: "Run XML-based tests",
		Long: `Run XML-based tests found in the tests/ directory.

Tests can be run individually by specifying test files, or filtered using include patterns.
Output verbosity can be controlled with flags.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTests(args, include, verbose, quiet, failOnly)
		},
	}

	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show verbose output with before/after buffer states")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Show only pass/fail summary")
	cmd.Flags().BoolVar(&failOnly, "fail-only", false, "Show output only for failing tests")
	cmd.Flags().StringSliceVarP(&include, "include", "i", nil, "Include tests matching pattern (can be used multiple times)")

	return cmd
}

// runTests executes the test runner
func runTests(testFiles []string, includePatterns []string, verbose, quiet, failOnly bool) error {
	var filesToTest []string

	if len(testFiles) == 0 {
		// Find all XML files in tests/ directory
		matches, err := filepath.Glob("tests/*.xml")
		if err != nil {
			return fmt.Errorf("finding test files: %w", err)
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
					return fmt.Errorf("invalid pattern %s: %w", pattern, err)
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
		return nil
	}

	// Create test environment
	env := edlisp.NewDefaultEnvironment()

	var results []*testing.TestResult
	totalTests := 0
	passedTests := 0

	for _, testFile := range filesToTest {
		result := testing.RunTestFile(testFile, env)
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
		return fmt.Errorf("some tests failed")
	}

	return nil
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