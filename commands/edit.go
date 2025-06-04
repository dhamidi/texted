package commands

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dhamidi/texted"
	"github.com/dhamidi/texted/edlisp"
	"github.com/dhamidi/texted/edlisp/parser"
	"github.com/dhamidi/texted/edlisp/writer"
)

// runEditArgs holds the arguments for the runEdit function
type runEditArgs struct {
	cmd          *cobra.Command
	scriptFormat string
	scriptFile   string
	inPlace      bool
	outputFile   string
	backupSuffix string
	verbose      bool
	quiet        bool
	dryRun       bool
	shell        bool
	sexp         bool
	json         bool
	outputFormat string
	files        []string
}

// runExpressionsArgs holds the arguments for the runExpressions function
type runExpressionsArgs struct {
	expressions  []string
	scriptFormat string
	outputFormat string
	verbose      bool
	quiet        bool
	files        []string
}

// evaluateExpressionsOnContentArgs holds the arguments for the evaluateExpressionsOnContent function
type evaluateExpressionsOnContentArgs struct {
	expressions  []string
	scriptFormat string
	outputFormat string
	verbose      bool
	quiet        bool
	content      string
	source       string
}

// processStdinArgs holds the arguments for the processStdin function
type processStdinArgs struct {
	script       string
	scriptFormat string
	outputFile   string
	verbose      bool
	quiet        bool
	dryRun       bool
}

// processFilesArgs holds the arguments for the processFiles function
type processFilesArgs struct {
	files        []string
	script       string
	scriptFormat string
	inPlace      bool
	outputFile   string
	backupSuffix string
	verbose      bool
	quiet        bool
	dryRun       bool
}

// processSingleFileToOutputArgs holds the arguments for the processSingleFileToOutput function
type processSingleFileToOutputArgs struct {
	filename     string
	script       string
	scriptFormat string
	outputFile   string
	verbose      bool
	quiet        bool
	dryRun       bool
}

// processSingleFileToStdoutArgs holds the arguments for the processSingleFileToStdout function
type processSingleFileToStdoutArgs struct {
	filename     string
	script       string
	scriptFormat string
	verbose      bool
	quiet        bool
	dryRun       bool
}

// processFilesInPlaceArgs holds the arguments for the processFilesInPlace function
type processFilesInPlaceArgs struct {
	files        []string
	script       string
	scriptFormat string
	backupSuffix string
	verbose      bool
	quiet        bool
	dryRun       bool
}

// NewEditCommand creates the edit subcommand.
func NewEditCommand() *cobra.Command {
	var (
		scriptFormat string
		scriptFile   string
		inPlace      bool
		outputFile   string
		backupSuffix string
		verbose      bool
		quiet        bool
		dryRun       bool
		shell        bool
		sexp         bool
		json         bool
		outputFormat string
	)

	cmd := &cobra.Command{
		Use:   "edit [flags] [files...]",
		Short: "Apply texted scripts to files",
		Long: `Apply texted scripts to one or more files or process stdin to stdout.

If no files are specified, the script is applied to stdin and the result is written to stdout.
If files are specified, the script is applied to each file in place.

Script formats:
  shell:  Default shell-like syntax (e.g., "search-forward hello")
  sexp:   S-expression syntax (e.g., "(search-forward \"hello\")")
  json:   JSON array syntax (e.g., ["search-forward", "hello"])`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEdit(&runEditArgs{
				cmd:          cmd,
				scriptFormat: scriptFormat,
				scriptFile:   scriptFile,
				inPlace:      inPlace,
				outputFile:   outputFile,
				backupSuffix: backupSuffix,
				verbose:      verbose,
				quiet:        quiet,
				dryRun:       dryRun,
				shell:        shell,
				sexp:         sexp,
				json:         json,
				outputFormat: outputFormat,
				files:        args,
			})
		},
	}

	// Script Input Options
	cmd.Flags().StringP("script", "s", "", "The texted script to execute")
	cmd.Flags().StringVarP(&scriptFile, "file", "f", "", "Read script from file")
	cmd.Flags().StringSliceP("expression", "e", nil, "Execute single expression and print result (can be used multiple times)")

	// Input/Output Options
	cmd.Flags().BoolVarP(&inPlace, "in-place", "i", false, "Edit files in place (modify original files)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write output to FILE (single file mode only)")
	cmd.Flags().StringVar(&backupSuffix, "backup", "", "Create backup files with SUFFIX when using --in-place")

	// Script Format Options
	cmd.Flags().StringVar(&scriptFormat, "format", "shell", "Specify script format: shell, sexp, json")
	cmd.Flags().BoolVar(&shell, "shell", false, "Force shell-like syntax parsing")
	cmd.Flags().BoolVar(&sexp, "sexp", false, "Force S-expression syntax parsing")
	cmd.Flags().BoolVar(&json, "json", false, "Force JSON syntax parsing")

	// Behavior Options
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress all output except errors")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Show what would be done without making changes")
	cmd.Flags().StringVar(&outputFormat, "output-format", "shell", "Output format for expression results: shell, sexp, json")

	return cmd
}

// runEdit handles the edit command execution.
func runEdit(args *runEditArgs) error {
	// Handle format shorthand flags
	if args.shell {
		args.scriptFormat = "shell"
	} else if args.sexp {
		args.scriptFormat = "sexp"
	} else if args.json {
		args.scriptFormat = "json"
	}

	// Validate format
	if !texted.IsValidFormat(args.scriptFormat) {
		return fmt.Errorf("invalid script format: %s (must be shell, sexp, or json)", args.scriptFormat)
	}

	// Validate flag combinations
	if args.outputFile != "" && len(args.files) > 1 {
		return fmt.Errorf("--output can only be used with a single file")
	}
	if args.backupSuffix != "" && !args.inPlace {
		return fmt.Errorf("--backup can only be used with --in-place")
	}
	if args.outputFile != "" && args.inPlace {
		return fmt.Errorf("--output and --in-place cannot be used together")
	}

	// Handle expressions
	expressions, err := args.cmd.Flags().GetStringSlice("expression")
	if err != nil {
		return fmt.Errorf("getting expression flag: %w", err)
	}

	if len(expressions) > 0 {
		return runExpressions(&runExpressionsArgs{
			expressions:  expressions,
			scriptFormat: args.scriptFormat,
			outputFormat: args.outputFormat,
			verbose:      args.verbose,
			quiet:        args.quiet,
			files:        args.files,
		})
	}

	// Get script content
	var script string

	if args.scriptFile != "" {
		content, err := os.ReadFile(args.scriptFile)
		if err != nil {
			return fmt.Errorf("reading script file: %w", err)
		}
		script = string(content)
	} else {
		// Get script from --script flag
		script, err = args.cmd.Flags().GetString("script")
		if err != nil {
			return fmt.Errorf("getting script flag: %w", err)
		}
		if script == "" {
			return fmt.Errorf("either --script, --file, or --expression must be specified")
		}
	}

	// If no files specified, process stdin to stdout
	if len(args.files) == 0 {
		return processStdin(&processStdinArgs{
			script:       script,
			scriptFormat: args.scriptFormat,
			outputFile:   args.outputFile,
			verbose:      args.verbose,
			quiet:        args.quiet,
			dryRun:       args.dryRun,
		})
	}

	// Process files
	return processFiles(&processFilesArgs{
		files:        args.files,
		script:       script,
		scriptFormat: args.scriptFormat,
		inPlace:      args.inPlace,
		outputFile:   args.outputFile,
		backupSuffix: args.backupSuffix,
		verbose:      args.verbose,
		quiet:        args.quiet,
		dryRun:       args.dryRun,
	})
}

// runExpressions handles the --expression flag by evaluating expressions and printing results
func runExpressions(args *runExpressionsArgs) error {
	// If no files specified, read from stdin
	if len(args.files) == 0 {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("reading stdin: %w", err)
		}
		return evaluateExpressionsOnContent(&evaluateExpressionsOnContentArgs{
			expressions:  args.expressions,
			scriptFormat: args.scriptFormat,
			outputFormat: args.outputFormat,
			verbose:      args.verbose,
			quiet:        args.quiet,
			content:      string(content),
			source:       "stdin",
		})
	}

	// If multiple files, evaluate expressions on each file
	for _, filename := range args.files {
		content, err := os.ReadFile(filename)
		if err != nil {
			if !args.quiet {
				fmt.Printf("Error reading %s: %v\n", filename, err)
			}
			return fmt.Errorf("reading %s: %w", filename, err)
		}

		if len(args.files) > 1 && !args.quiet {
			fmt.Printf("=== %s ===\n", filename)
		}

		err = evaluateExpressionsOnContent(&evaluateExpressionsOnContentArgs{
			expressions:  args.expressions,
			scriptFormat: args.scriptFormat,
			outputFormat: args.outputFormat,
			verbose:      args.verbose,
			quiet:        args.quiet,
			content:      string(content),
			source:       filename,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// evaluateExpressionsOnContent evaluates expressions on the given content
func evaluateExpressionsOnContent(args *evaluateExpressionsOnContentArgs) error {
	buffer := edlisp.NewBuffer(args.content)
	env := edlisp.NewDefaultEnvironment()

	for i, expr := range args.expressions {
		if args.verbose && !args.quiet {
			fmt.Printf("Evaluating expression %d on %s: %s\n", i+1, args.source, expr)
		}

		// Parse the expression based on format
		var program []edlisp.Value
		var err error

		switch args.scriptFormat {
		case "shell", "sexp":
			program, err = parser.ParseString(expr)
		case "json":
			program, err = parser.ParseJSONString(expr)
		default:
			return fmt.Errorf("unsupported script format: %s", args.scriptFormat)
		}

		if err != nil {
			if !args.quiet {
				fmt.Printf("Error parsing expression %d: %v\n", i+1, err)
			}
			return fmt.Errorf("parsing expression: %w", err)
		}

		// Execute the expression and get the result value (not buffer content)
		result, err := edlisp.Eval(program, env, buffer)
		if err != nil {
			if !args.quiet {
				fmt.Printf("Error in expression %d: %v\n", i+1, err)
			}
			return err
		}

		if !args.quiet {
			// Convert the result value to string representation using specified format
			var writerFormat writer.Format
			switch args.outputFormat {
			case "shell":
				writerFormat = writer.FormatShell
			case "sexp":
				writerFormat = writer.FormatSExp
			case "json":
				writerFormat = writer.FormatJSON
			default:
				return fmt.Errorf("invalid output format: %s (must be shell, sexp, or json)", args.outputFormat)
			}

			w, err := writer.NewWriter(writerFormat)
			if err != nil {
				return fmt.Errorf("creating writer: %w", err)
			}

			var buf strings.Builder

			// For shell format, wrap the result in a list since shell writer expects lists
			if args.outputFormat == "shell" {
				list := edlisp.NewList(result)
				err = w.WriteValue(&buf, list)
			} else {
				err = w.WriteValue(&buf, result)
			}

			if err != nil {
				return fmt.Errorf("writing result: %w", err)
			}

			fmt.Printf("%s\n", buf.String())
		}
	}
	return nil
}

// processStdin handles processing stdin to stdout or output file
func processStdin(args *processStdinArgs) error {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}

	if args.verbose && !args.quiet {
		fmt.Printf("Processing stdin with script in %s format\n", args.scriptFormat)
	}

	if args.dryRun {
		if !args.quiet {
			fmt.Printf("Would process %d bytes from stdin\n", len(content))
		}
		return nil
	}

	result, err := texted.ExecuteScriptWithFormat(string(content), args.script, args.scriptFormat)
	if err != nil {
		return err
	}

	if args.outputFile != "" {
		if args.verbose && !args.quiet {
			fmt.Printf("Writing output to %s\n", args.outputFile)
		}
		return os.WriteFile(args.outputFile, []byte(result), 0644)
	}

	_, err = os.Stdout.WriteString(result)
	return err
}

// processFiles handles processing one or more files
func processFiles(args *processFilesArgs) error {
	if len(args.files) == 1 && args.outputFile != "" {
		// Single file with output redirection
		return processSingleFileToOutput(&processSingleFileToOutputArgs{
			filename:     args.files[0],
			script:       args.script,
			scriptFormat: args.scriptFormat,
			outputFile:   args.outputFile,
			verbose:      args.verbose,
			quiet:        args.quiet,
			dryRun:       args.dryRun,
		})
	}

	if !args.inPlace && args.outputFile == "" {
		// Default behavior: output to stdout (for single file) or error for multiple
		if len(args.files) == 1 {
			return processSingleFileToStdout(&processSingleFileToStdoutArgs{
				filename:     args.files[0],
				script:       args.script,
				scriptFormat: args.scriptFormat,
				verbose:      args.verbose,
				quiet:        args.quiet,
				dryRun:       args.dryRun,
			})
		}
		return fmt.Errorf("multiple files require --in-place or --output flag")
	}

	// In-place editing
	return processFilesInPlace(&processFilesInPlaceArgs{
		files:        args.files,
		script:       args.script,
		scriptFormat: args.scriptFormat,
		backupSuffix: args.backupSuffix,
		verbose:      args.verbose,
		quiet:        args.quiet,
		dryRun:       args.dryRun,
	})
}

// processSingleFileToOutput processes a single file and writes to specified output
func processSingleFileToOutput(args *processSingleFileToOutputArgs) error {
	if args.verbose && !args.quiet {
		fmt.Printf("Processing %s -> %s\n", args.filename, args.outputFile)
	}

	if args.dryRun {
		if !args.quiet {
			fmt.Printf("Would process %s and write to %s\n", args.filename, args.outputFile)
		}
		return nil
	}

	content, err := os.ReadFile(args.filename)
	if err != nil {
		return fmt.Errorf("reading %s: %w", args.filename, err)
	}

	result, err := texted.ExecuteScriptWithFormat(string(content), args.script, args.scriptFormat)
	if err != nil {
		return fmt.Errorf("processing %s: %w", args.filename, err)
	}

	return os.WriteFile(args.outputFile, []byte(result), 0644)
}

// processSingleFileToStdout processes a single file and writes to stdout
func processSingleFileToStdout(args *processSingleFileToStdoutArgs) error {
	if args.verbose && !args.quiet {
		fmt.Printf("Processing %s -> stdout\n", args.filename)
	}

	if args.dryRun {
		if !args.quiet {
			fmt.Printf("Would process %s and write to stdout\n", args.filename)
		}
		return nil
	}

	content, err := os.ReadFile(args.filename)
	if err != nil {
		return fmt.Errorf("reading %s: %w", args.filename, err)
	}

	result, err := texted.ExecuteScriptWithFormat(string(content), args.script, args.scriptFormat)
	if err != nil {
		return fmt.Errorf("processing %s: %w", args.filename, err)
	}

	_, err = os.Stdout.WriteString(result)
	return err
}

// processFilesInPlace processes files in place with optional backup
func processFilesInPlace(args *processFilesInPlaceArgs) error {
	hasErrors := false
	for _, filename := range args.files {
		if args.verbose && !args.quiet {
			fmt.Printf("Processing %s in place\n", filename)
		}

		if args.dryRun {
			if !args.quiet {
				fmt.Printf("Would edit %s in place", filename)
				if args.backupSuffix != "" {
					fmt.Printf(" (backup to %s%s)", filename, args.backupSuffix)
				}
				fmt.Println()
			}
			continue
		}

		// Create backup if requested
		if args.backupSuffix != "" {
			backupFile := filename + args.backupSuffix
			if args.verbose && !args.quiet {
				fmt.Printf("Creating backup %s\n", backupFile)
			}

			content, err := os.ReadFile(filename)
			if err != nil {
				if !args.quiet {
					fmt.Printf("✗ Failed to read %s for backup: %v\n", filename, err)
				}
				hasErrors = true
				continue
			}

			if err := os.WriteFile(backupFile, content, 0644); err != nil {
				if !args.quiet {
					fmt.Printf("✗ Failed to create backup %s: %v\n", backupFile, err)
				}
				hasErrors = true
				continue
			}
		}

		// Process the file
		content, err := os.ReadFile(filename)
		if err != nil {
			if !args.quiet {
				fmt.Printf("✗ Failed to read %s: %v\n", filename, err)
			}
			hasErrors = true
			continue
		}

		result, err := texted.ExecuteScriptWithFormat(string(content), args.script, args.scriptFormat)
		if err != nil {
			if !args.quiet {
				fmt.Printf("✗ Failed to process %s: %v\n", filename, err)
			}
			hasErrors = true
			continue
		}

		if err := os.WriteFile(filename, []byte(result), 0644); err != nil {
			if !args.quiet {
				fmt.Printf("✗ Failed to write %s: %v\n", filename, err)
			}
			hasErrors = true
			continue
		}

		if !args.quiet {
			fmt.Printf("✓ Successfully edited %s\n", filename)
		}
	}

	if hasErrors {
		return fmt.Errorf("some files could not be edited")
	}

	return nil
}
