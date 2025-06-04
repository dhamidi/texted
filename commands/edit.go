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
	var scriptFormat string
	var scriptFile string

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
			return runEdit(cmd, scriptFormat, scriptFile, args)
		},
	}

	cmd.Flags().StringVar(&scriptFormat, "script-format", "shell", "Script format: shell, sexp, json")
	cmd.Flags().StringVar(&scriptFile, "script-file", "", "Read script from file instead of --script flag")
	cmd.Flags().String("script", "", "The texted script to execute")

	return cmd
}

// runEdit handles the edit command execution.
func runEdit(cmd *cobra.Command, scriptFormat, scriptFile string, files []string) error {
	// Validate format
	if !texted.IsValidFormat(scriptFormat) {
		return fmt.Errorf("invalid script format: %s (must be shell, sexp, or json)", scriptFormat)
	}

	// Get script content
	var script string
	var err error

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
			return fmt.Errorf("either --script or --script-file must be specified")
		}
	}

	// If no files specified, process stdin to stdout
	if len(files) == 0 {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("reading stdin: %w", err)
		}

		result, err := texted.ExecuteScriptWithFormat(string(content), script, scriptFormat)
		if err != nil {
			return err
		}

		_, err = os.Stdout.WriteString(result)
		return err
	}

	// Process files
	results, err := texted.EditFilesWithFormat(files, script, scriptFormat)
	if err != nil {
		return err
	}

	// Report results
	hasErrors := false
	for _, result := range results {
		if result.Success {
			fmt.Printf("✓ Successfully edited %s\n", result.Filename)
		} else {
			fmt.Printf("✗ Failed to edit %s: %v\n", result.Filename, result.Error)
			hasErrors = true
		}
	}

	if hasErrors {
		return fmt.Errorf("some files could not be edited")
	}

	return nil
}
