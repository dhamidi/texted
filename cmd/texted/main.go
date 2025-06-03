package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/dhamidi/texted/commands"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "texted",
		Short: "A scriptable, headless text editor",
		Long:  "Texted is a scriptable, headless text editor for automated file editing. It processes scripts written in shell-like syntax, S-expressions, or JSON format to perform text transformations.",
	}

	// Add subcommands
	rootCmd.AddCommand(commands.NewParseCommand())
	rootCmd.AddCommand(commands.NewTestCommand())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
