package texted

import (
	"fmt"
	"os"

	"github.com/dhamidi/texted/edlisp"
	"github.com/dhamidi/texted/edlisp/parser"
)

// EditResult represents the result of editing a single file.
type EditResult struct {
	Filename string
	Success  bool
	Error    error
}

// ExecuteScript executes a texted script on the given input and returns the result.
func ExecuteScript(input, script string) (string, error) {
	buf := edlisp.NewBuffer(input)

	program, err := parser.ParseString(script)
	if err != nil {
		return "", fmt.Errorf("parsing script: %w", err)
	}

	env := edlisp.NewDefaultEnvironment()
	_, err = edlisp.Eval(program, env, buf)
	if err != nil {
		return "", fmt.Errorf("script execution failed: %w", err)
	}

	return buf.String(), nil
}

// ExecuteScriptWithFormat executes a texted script with a specific format on the given input.
func ExecuteScriptWithFormat(input, script, format string) (string, error) {
	if !IsValidFormat(format) {
		return "", fmt.Errorf("invalid script format: %s (must be shell, sexp, or json)", format)
	}

	buf := edlisp.NewBuffer(input)

	var program []edlisp.Value
	var err error

	switch format {
	case "shell", "sexp":
		program, err = parser.ParseString(script)
	case "json":
		program, err = parser.ParseJSONString(script)
	default:
		return "", fmt.Errorf("unsupported script format: %s", format)
	}

	if err != nil {
		return "", fmt.Errorf("parsing script: %w", err)
	}

	env := edlisp.NewDefaultEnvironment()
	_, err = edlisp.Eval(program, env, buf)
	if err != nil {
		return "", fmt.Errorf("script execution failed: %w", err)
	}

	return buf.String(), nil
}

// EditFile applies a texted script to a file.
func EditFile(filename, script string) error {
	content, err := readFile(filename)
	if err != nil {
		return err
	}

	modified, err := ExecuteScript(content, script)
	if err != nil {
		return err
	}

	return writeFile(filename, modified)
}

// EditFileWithFormat applies a texted script with a specific format to a file.
func EditFileWithFormat(filename, script, format string) error {
	content, err := readFile(filename)
	if err != nil {
		return err
	}

	modified, err := ExecuteScriptWithFormat(content, script, format)
	if err != nil {
		return err
	}

	return writeFile(filename, modified)
}

// EditFiles applies a texted script to multiple files.
func EditFiles(files []string, script string) ([]EditResult, error) {
	results := make([]EditResult, 0, len(files))

	for _, filename := range files {
		result := EditResult{Filename: filename}

		err := EditFile(filename, script)
		if err != nil {
			result.Success = false
			result.Error = err
		} else {
			result.Success = true
		}

		results = append(results, result)
	}

	return results, nil
}

// EditFilesWithFormat applies a texted script with a specific format to multiple files.
func EditFilesWithFormat(files []string, script, format string) ([]EditResult, error) {
	results := make([]EditResult, 0, len(files))

	for _, filename := range files {
		result := EditResult{Filename: filename}

		err := EditFileWithFormat(filename, script, format)
		if err != nil {
			result.Success = false
			result.Error = err
		} else {
			result.Success = true
		}

		results = append(results, result)
	}

	return results, nil
}

// IsValidFormat checks if a format string is valid.
func IsValidFormat(format string) bool {
	switch format {
	case "shell", "sexp", "json":
		return true
	default:
		return false
	}
}

// readFile reads the content of a file.
func readFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	return string(content), nil
}

// writeFile writes content to a file.
func writeFile(filename, content string) error {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filename, err)
	}
	return nil
}
