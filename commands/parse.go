package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/dhamidi/texted/edlisp"
	"github.com/dhamidi/texted/edlisp/parser"
	"github.com/dhamidi/texted/edlisp/writer"
)

// runParseArgs holds the arguments for the runParse function
type runParseArgs struct {
	inputFormat  string
	outputFormat string
}

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
			return runParse(&runParseArgs{
				inputFormat:  inputFormat,
				outputFormat: outputFormat,
			})
		},
	}

	cmd.Flags().StringVar(&inputFormat, "input-format", "shell", "Input script format: shell, sexp, json")
	cmd.Flags().StringVar(&outputFormat, "output-format", "sexp", "Output script format: shell, sexp, json")

	return cmd
}

// runParse handles the parse command execution.
func runParse(args *runParseArgs) error {
	// Validate formats
	if !isValidFormat(args.inputFormat) {
		return fmt.Errorf("invalid input format: %s (must be shell, sexp, or json)", args.inputFormat)
	}
	if !isValidFormat(args.outputFormat) {
		return fmt.Errorf("invalid output format: %s (must be shell, sexp, or json)", args.outputFormat)
	}

	// Parse input from stdin
	expressions, err := parseInput(os.Stdin, args.inputFormat)
	if err != nil {
		return fmt.Errorf("parsing input: %w", err)
	}

	// Convert to output format and write to stdout
	return writeOutput(os.Stdout, expressions, args.outputFormat)
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
	writerInstance, err := writer.NewWriter(writer.Format(format))
	if err != nil {
		return fmt.Errorf("unsupported output format: %s", format)
	}
	return writerInstance.Write(w, expressions)
}
