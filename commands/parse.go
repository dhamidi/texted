package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dhamidi/texted/edlisp"
	"github.com/dhamidi/texted/edlisp/parser"
)

// NewParseCommand creates the parse subcommand.
func NewParseCommand() *cobra.Command {
	var inputFormat string
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "parse",
		Short: "Convert scripts between different formats",
		Long: `Convert scripts between different formats without executing them.

Reads a script from stdin in the specified input format and writes the parsed 
version to stdout in the specified output format. This is useful for converting 
between shell-like syntax, S-expressions, and JSON formats.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runParse(inputFormat, outputFormat)
		},
	}

	cmd.Flags().StringVar(&inputFormat, "input-format", "shell", "Input script format: shell, sexp, json")
	cmd.Flags().StringVar(&outputFormat, "output-format", "sexp", "Output script format: shell, sexp, json")

	return cmd
}

// runParse handles the parse command execution.
func runParse(inputFormat, outputFormat string) error {
	// Validate formats
	if !isValidFormat(inputFormat) {
		return fmt.Errorf("invalid input format: %s (must be shell, sexp, or json)", inputFormat)
	}
	if !isValidFormat(outputFormat) {
		return fmt.Errorf("invalid output format: %s (must be shell, sexp, or json)", outputFormat)
	}

	// Parse input from stdin
	expressions, err := parseInput(os.Stdin, inputFormat)
	if err != nil {
		return fmt.Errorf("parsing input: %w", err)
	}

	// Convert to output format and write to stdout
	return writeOutput(os.Stdout, expressions, outputFormat)
}

// isValidFormat checks if a format string is valid.
func isValidFormat(format string) bool {
	switch format {
	case "shell", "sexp", "json":
		return true
	default:
		return false
	}
}

// parseInput parses expressions from a reader based on the input format.
func parseInput(r io.Reader, format string) ([]edlisp.Value, error) {
	switch format {
	case "shell":
		return parser.ParseReader(r)
	case "sexp":
		return parser.ParseReader(r)
	case "json":
		return parser.ParseJSONReader(r)
	default:
		return nil, fmt.Errorf("unsupported input format: %s", format)
	}
}

// writeOutput writes expressions to a writer in the specified output format.
func writeOutput(w io.Writer, expressions []edlisp.Value, format string) error {
	switch format {
	case "shell":
		return writeShellFormat(w, expressions)
	case "sexp":
		return writeSExpFormat(w, expressions)
	case "json":
		return writeJSONFormat(w, expressions)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// writeShellFormat writes expressions in shell-like format.
func writeShellFormat(w io.Writer, expressions []edlisp.Value) error {
	for _, expr := range expressions {
		shellStr, err := toShellString(expr)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, shellStr); err != nil {
			return err
		}
	}
	return nil
}

// writeSExpFormat writes expressions in S-expression format.
func writeSExpFormat(w io.Writer, expressions []edlisp.Value) error {
	for _, expr := range expressions {
		if _, err := fmt.Fprintln(w, expr); err != nil {
			return err
		}
	}
	return nil
}

// writeJSONFormat writes expressions in JSON format.
func writeJSONFormat(w io.Writer, expressions []edlisp.Value) error {
	encoder := json.NewEncoder(w)
	for _, expr := range expressions {
		jsonValue, err := toJSONValue(expr)
		if err != nil {
			return err
		}
		if err := encoder.Encode(jsonValue); err != nil {
			return err
		}
	}
	return nil
}

// toShellString converts an edlisp value to shell-like string format.
func toShellString(value edlisp.Value) (string, error) {
	list, ok := value.(*edlisp.List)
	if !ok {
		return "", fmt.Errorf("can only convert lists to shell format, got %T", value)
	}

	if list.IsEmpty() {
		return "", nil
	}

	var parts []string
	for _, element := range list.Elements {
		part, err := valueToShellToken(element)
		if err != nil {
			return "", err
		}
		parts = append(parts, part)
	}

	return strings.Join(parts, " "), nil
}

// valueToShellToken converts a single value to a shell token.
func valueToShellToken(value edlisp.Value) (string, error) {
	switch v := value.(type) {
	case *edlisp.Symbol:
		return v.Name, nil
	case *edlisp.String:
		return strconv.Quote(v.Value), nil
	case *edlisp.Number:
		return fmt.Sprintf("%v", v), nil
	case *edlisp.List:
		return "", fmt.Errorf("nested lists are not supported in shell format")
	default:
		return "", fmt.Errorf("unsupported value type for shell format: %T", value)
	}
}

// toJSONValue converts an edlisp value to a JSON-compatible value.
func toJSONValue(value edlisp.Value) (interface{}, error) {
	switch v := value.(type) {
	case *edlisp.List:
		var jsonArray []interface{}
		for _, element := range v.Elements {
			jsonValue, err := toJSONValue(element)
			if err != nil {
				return nil, err
			}
			jsonArray = append(jsonArray, jsonValue)
		}
		return jsonArray, nil
	case *edlisp.Symbol:
		return v.Name, nil
	case *edlisp.String:
		return v.Value, nil
	case *edlisp.Number:
		return v.Value, nil
	default:
		return nil, fmt.Errorf("unsupported value type for JSON: %T", value)
	}
}
