package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/dhamidi/texted/edlisp"
)

// NewDocCommand creates the doc subcommand.
func NewDocCommand() *cobra.Command {
	var category string
	var verbose bool

	cmd := &cobra.Command{
		Use:   "doc [function-name]",
		Short: "Show documentation for functions",
		Long: `Show documentation for texted functions.

Without arguments, lists all available functions.
With a function name, shows detailed documentation for that function.
Use --category to filter by function category.`,
		Example: `  texted doc                    # List all functions
  texted doc search-forward     # Show docs for search-forward
  texted doc --category search  # Show all search functions
  texted doc --verbose          # Show detailed list with summaries`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoc(args, category, verbose)
		},
	}

	cmd.Flags().StringVar(&category, "category", "", "Filter functions by category")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show function summaries in list view")

	return cmd
}

// runDoc handles the doc command execution.
func runDoc(args []string, category string, verbose bool) error {
	// If a specific function is requested
	if len(args) == 1 {
		return showFunctionDoc(args[0])
	}

	// If too many arguments
	if len(args) > 1 {
		return fmt.Errorf("doc command accepts at most one argument (function name)")
	}

	// List functions
	return listFunctions(category, verbose)
}

// showFunctionDoc displays detailed documentation for a specific function.
func showFunctionDoc(functionName string) error {
	doc, exists := edlisp.GetDocumentation(functionName)
	if !exists {
		return fmt.Errorf("no documentation found for function: %s", functionName)
	}

	// Function header
	fmt.Printf("# %s\n\n", doc.Name)
	fmt.Printf("**%s**\n\n", doc.Summary)

	// Description
	if doc.Description != "" {
		fmt.Printf("## Description\n\n%s\n\n", doc.Description)
	}

	// Parameters
	if len(doc.Parameters) > 0 {
		fmt.Printf("## Parameters\n\n")
		for _, param := range doc.Parameters {
			optional := ""
			if param.Optional {
				optional = " (optional)"
			}
			fmt.Printf("- **%s** (%s)%s: %s\n", param.Name, param.Type, optional, param.Description)
		}
		fmt.Printf("\n")
	}

	// Examples
	if len(doc.Examples) > 0 {
		fmt.Printf("## Examples\n\n")
		for i, example := range doc.Examples {
			if i > 0 {
				fmt.Printf("\n")
			}
			fmt.Printf("### %s\n\n", example.Description)
			if example.Buffer != "" {
				fmt.Printf("**Initial buffer:**\n```\n%s\n```\n\n", example.Buffer)
			}
			fmt.Printf("**Command:**\n```\n%s\n```\n\n", example.Input)
			fmt.Printf("**Result:**\n%s\n", example.Output)
		}
		fmt.Printf("\n")
	}

	// Category
	if doc.Category != "" {
		fmt.Printf("## Category\n\n%s\n\n", doc.Category)
	}

	// See also
	if len(doc.SeeAlso) > 0 {
		fmt.Printf("## See Also\n\n")
		for _, related := range doc.SeeAlso {
			fmt.Printf("- %s\n", related)
		}
		fmt.Printf("\n")
	}

	return nil
}

// listFunctions shows a list of available functions, optionally filtered by category.
func listFunctions(category string, verbose bool) error {
	var docs []edlisp.FunctionDoc

	if category != "" {
		docs = edlisp.GetDocumentationByCategory(category)
		if len(docs) == 0 {
			return fmt.Errorf("no functions found in category: %s", category)
		}
	} else {
		docs = edlisp.GetAllDocumentation()
	}

	if len(docs) == 0 {
		fmt.Println("No documented functions found.")
		return nil
	}

	// Show header
	if category != "" {
		fmt.Printf("Functions in category '%s':\n\n", category)
	} else {
		fmt.Printf("Available functions (%d total):\n\n", len(docs))
	}

	if verbose {
		// Group by category for verbose output
		categoryMap := make(map[string][]edlisp.FunctionDoc)
		for _, doc := range docs {
			cat := doc.Category
			if cat == "" {
				cat = "other"
			}
			categoryMap[cat] = append(categoryMap[cat], doc)
		}

		for catName, catDocs := range categoryMap {
			fmt.Printf("## %s\n\n", strings.Title(catName))
			for _, doc := range catDocs {
				fmt.Printf("- **%s**: %s\n", doc.Name, doc.Summary)
			}
			fmt.Printf("\n")
		}
	} else {
		// Simple list
		for _, doc := range docs {
			fmt.Printf("  %s\n", doc.Name)
		}
		fmt.Printf("\nUse 'texted doc <function-name>' for detailed documentation.\n")
		fmt.Printf("Use 'texted doc --verbose' to see function summaries.\n")
	}

	return nil
}
