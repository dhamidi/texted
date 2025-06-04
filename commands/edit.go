package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/dhamidi/texted"
)

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
			return runEdit(cmd, scriptFormat, scriptFile, inPlace, outputFile, backupSuffix, verbose, quiet, dryRun, shell, sexp, json, args)
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

	return cmd
}

// runEdit handles the edit command execution.
func runEdit(cmd *cobra.Command, scriptFormat, scriptFile string, inPlace bool, outputFile, backupSuffix string, verbose, quiet, dryRun, shell, sexp, json bool, files []string) error {
	// Handle format shorthand flags
	if shell {
		scriptFormat = "shell"
	} else if sexp {
		scriptFormat = "sexp"
	} else if json {
		scriptFormat = "json"
	}

	// Validate format
	if !texted.IsValidFormat(scriptFormat) {
		return fmt.Errorf("invalid script format: %s (must be shell, sexp, or json)", scriptFormat)
	}

	// Validate flag combinations
	if outputFile != "" && len(files) > 1 {
		return fmt.Errorf("--output can only be used with a single file")
	}
	if backupSuffix != "" && !inPlace {
		return fmt.Errorf("--backup can only be used with --in-place")
	}
	if outputFile != "" && inPlace {
		return fmt.Errorf("--output and --in-place cannot be used together")
	}

	// Handle expressions
	expressions, err := cmd.Flags().GetStringSlice("expression")
	if err != nil {
		return fmt.Errorf("getting expression flag: %w", err)
	}

	if len(expressions) > 0 {
		return runExpressions(expressions, scriptFormat, verbose, quiet)
	}

	// Get script content
	var script string

	if scriptFile != "" {
		content, err := os.ReadFile(scriptFile)
		if err != nil {
			return fmt.Errorf("reading script file: %w", err)
		}
		script = string(content)
	} else {
		// Get script from --script flag
		script, err = cmd.Flags().GetString("script")
		if err != nil {
			return fmt.Errorf("getting script flag: %w", err)
		}
		if script == "" {
			return fmt.Errorf("either --script, --file, or --expression must be specified")
		}
	}

	// If no files specified, process stdin to stdout
	if len(files) == 0 {
		return processStdin(script, scriptFormat, outputFile, verbose, quiet, dryRun)
	}

	// Process files
	return processFiles(files, script, scriptFormat, inPlace, outputFile, backupSuffix, verbose, quiet, dryRun)
}

// runExpressions handles the --expression flag by evaluating expressions and printing results
func runExpressions(expressions []string, scriptFormat string, verbose, quiet bool) error {
	for i, expr := range expressions {
		if verbose && !quiet {
			fmt.Printf("Evaluating expression %d: %s\n", i+1, expr)
		}

		// For expressions, we'll execute them on an empty buffer and get the result
		result, err := texted.ExecuteScriptWithFormat("", expr, scriptFormat)
		if err != nil {
			if !quiet {
				fmt.Printf("Error in expression %d: %v\n", i+1, err)
			}
			return err
		}

		if !quiet {
			fmt.Printf("%s\n", result)
		}
	}
	return nil
}

// processStdin handles processing stdin to stdout or output file
func processStdin(script, scriptFormat, outputFile string, verbose, quiet, dryRun bool) error {
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("reading stdin: %w", err)
	}

	if verbose && !quiet {
		fmt.Printf("Processing stdin with script in %s format\n", scriptFormat)
	}

	if dryRun {
		if !quiet {
			fmt.Printf("Would process %d bytes from stdin\n", len(content))
		}
		return nil
	}

	result, err := texted.ExecuteScriptWithFormat(string(content), script, scriptFormat)
	if err != nil {
		return err
	}

	if outputFile != "" {
		if verbose && !quiet {
			fmt.Printf("Writing output to %s\n", outputFile)
		}
		return os.WriteFile(outputFile, []byte(result), 0644)
	}

	_, err = os.Stdout.WriteString(result)
	return err
}

// processFiles handles processing one or more files
func processFiles(files []string, script, scriptFormat string, inPlace bool, outputFile, backupSuffix string, verbose, quiet, dryRun bool) error {
	if len(files) == 1 && outputFile != "" {
		// Single file with output redirection
		return processSingleFileToOutput(files[0], script, scriptFormat, outputFile, verbose, quiet, dryRun)
	}

	if !inPlace && outputFile == "" {
		// Default behavior: output to stdout (for single file) or error for multiple
		if len(files) == 1 {
			return processSingleFileToStdout(files[0], script, scriptFormat, verbose, quiet, dryRun)
		}
		return fmt.Errorf("multiple files require --in-place or --output flag")
	}

	// In-place editing
	return processFilesInPlace(files, script, scriptFormat, backupSuffix, verbose, quiet, dryRun)
}

// processSingleFileToOutput processes a single file and writes to specified output
func processSingleFileToOutput(filename, script, scriptFormat, outputFile string, verbose, quiet, dryRun bool) error {
	if verbose && !quiet {
		fmt.Printf("Processing %s -> %s\n", filename, outputFile)
	}

	if dryRun {
		if !quiet {
			fmt.Printf("Would process %s and write to %s\n", filename, outputFile)
		}
		return nil
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filename, err)
	}

	result, err := texted.ExecuteScriptWithFormat(string(content), script, scriptFormat)
	if err != nil {
		return fmt.Errorf("processing %s: %w", filename, err)
	}

	return os.WriteFile(outputFile, []byte(result), 0644)
}

// processSingleFileToStdout processes a single file and writes to stdout
func processSingleFileToStdout(filename, script, scriptFormat string, verbose, quiet, dryRun bool) error {
	if verbose && !quiet {
		fmt.Printf("Processing %s -> stdout\n", filename)
	}

	if dryRun {
		if !quiet {
			fmt.Printf("Would process %s and write to stdout\n", filename)
		}
		return nil
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filename, err)
	}

	result, err := texted.ExecuteScriptWithFormat(string(content), script, scriptFormat)
	if err != nil {
		return fmt.Errorf("processing %s: %w", filename, err)
	}

	_, err = os.Stdout.WriteString(result)
	return err
}

// processFilesInPlace processes files in place with optional backup
func processFilesInPlace(files []string, script, scriptFormat, backupSuffix string, verbose, quiet, dryRun bool) error {
	hasErrors := false
	for _, filename := range files {
		if verbose && !quiet {
			fmt.Printf("Processing %s in place\n", filename)
		}

		if dryRun {
			if !quiet {
				fmt.Printf("Would edit %s in place", filename)
				if backupSuffix != "" {
					fmt.Printf(" (backup to %s%s)", filename, backupSuffix)
				}
				fmt.Println()
			}
			continue
		}

		// Create backup if requested
		if backupSuffix != "" {
			backupFile := filename + backupSuffix
			if verbose && !quiet {
				fmt.Printf("Creating backup %s\n", backupFile)
			}

			content, err := os.ReadFile(filename)
			if err != nil {
				if !quiet {
					fmt.Printf("✗ Failed to read %s for backup: %v\n", filename, err)
				}
				hasErrors = true
				continue
			}

			if err := os.WriteFile(backupFile, content, 0644); err != nil {
				if !quiet {
					fmt.Printf("✗ Failed to create backup %s: %v\n", backupFile, err)
				}
				hasErrors = true
				continue
			}
		}

		// Process the file
		content, err := os.ReadFile(filename)
		if err != nil {
			if !quiet {
				fmt.Printf("✗ Failed to read %s: %v\n", filename, err)
			}
			hasErrors = true
			continue
		}

		result, err := texted.ExecuteScriptWithFormat(string(content), script, scriptFormat)
		if err != nil {
			if !quiet {
				fmt.Printf("✗ Failed to process %s: %v\n", filename, err)
			}
			hasErrors = true
			continue
		}

		if err := os.WriteFile(filename, []byte(result), 0644); err != nil {
			if !quiet {
				fmt.Printf("✗ Failed to write %s: %v\n", filename, err)
			}
			hasErrors = true
			continue
		}

		if !quiet {
			fmt.Printf("✓ Successfully edited %s\n", filename)
		}
	}

	if hasErrors {
		return fmt.Errorf("some files could not be edited")
	}

	return nil
}
