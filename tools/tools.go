package tools

import (
	"fmt"
	"os"

	"github.com/dhamidi/texted/edlisp"
	"github.com/dhamidi/texted/edlisp/parser"
)

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

func readFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	return string(content), nil
}

func writeFile(filename, content string) error {
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filename, err)
	}
	return nil
}
