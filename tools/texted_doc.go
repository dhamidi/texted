package tools

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/dhamidi/texted/edlisp"
	"github.com/mark3labs/mcp-go/mcp"
)

//go:embed texted_doc_description.txt
var textedDocDescription string

func NewTextedDocTool() mcp.Tool {
	return mcp.NewTool("texted_doc",
		mcp.WithDescription(textedDocDescription),
		mcp.WithString("function_name",
			mcp.Description("Specific function name to show documentation for"),
		),
		mcp.WithString("category",
			mcp.Description("Filter functions by category (e.g., 'search', 'movement', 'editing')"),
		),
		mcp.WithString("verbose",
			mcp.Description("Show function summaries in list view (true/false)"),
		),
	)
}

func TextedDocHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	functionName := request.GetString("function_name", "")
	category := request.GetString("category", "")
	verboseStr := request.GetString("verbose", "false")
	verbose := verboseStr == "true"

	// Validate mutually exclusive parameters
	if functionName != "" && category != "" {
		return mcp.NewToolResultError("function_name and category parameters are mutually exclusive"), nil
	}

	// Handle specific function documentation
	if functionName != "" {
		return handleFunctionDoc(functionName)
	}

	// Handle category filtering
	if category != "" {
		return handleCategoryListing(category, verbose)
	}

	// Handle listing all functions
	return handleAllFunctions(verbose)
}

func handleFunctionDoc(functionName string) (*mcp.CallToolResult, error) {
	doc, exists := edlisp.GetDocumentation(functionName)
	if !exists {
		// Get similar function names for suggestions
		allDocs := edlisp.GetAllDocumentation()
		var suggestions []string
		for _, d := range allDocs {
			if strings.Contains(d.Name, functionName) || strings.Contains(functionName, d.Name) {
				suggestions = append(suggestions, d.Name)
				if len(suggestions) >= 3 {
					break
				}
			}
		}

		errorMsg := fmt.Sprintf("Function '%s' not found", functionName)
		if len(suggestions) > 0 {
			errorMsg += fmt.Sprintf(". Did you mean: %s", strings.Join(suggestions, ", "))
		}
		return mcp.NewToolResultError(errorMsg), nil
	}

	output := formatFunctionDoc(doc)
	return mcp.NewToolResultText(output), nil
}

func handleCategoryListing(category string, verbose bool) (*mcp.CallToolResult, error) {
	docs := edlisp.GetDocumentationByCategory(category)
	if len(docs) == 0 {
		// Get available categories for suggestions
		categories := edlisp.GetCategories()
		errorMsg := fmt.Sprintf("No functions found in category '%s'", category)
		if len(categories) > 0 {
			errorMsg += fmt.Sprintf(". Available categories: %s", strings.Join(categories, ", "))
		}
		return mcp.NewToolResultError(errorMsg), nil
	}

	output := formatFunctionList(docs, fmt.Sprintf("Functions in category '%s'", category), verbose)
	return mcp.NewToolResultText(output), nil
}

func handleAllFunctions(verbose bool) (*mcp.CallToolResult, error) {
	docs := edlisp.GetAllDocumentation()
	if len(docs) == 0 {
		return mcp.NewToolResultError("No functions documented"), nil
	}

	output := formatFunctionList(docs, "All functions", verbose)
	return mcp.NewToolResultText(output), nil
}

func formatFunctionDoc(doc edlisp.FunctionDoc) string {
	var output strings.Builder

	// Function name and summary
	output.WriteString(fmt.Sprintf("# %s\n\n", doc.Name))
	if doc.Summary != "" {
		output.WriteString(fmt.Sprintf("**%s**\n\n", doc.Summary))
	}

	// Description
	if doc.Description != "" {
		output.WriteString(fmt.Sprintf("%s\n\n", doc.Description))
	}

	// Parameters
	if len(doc.Parameters) > 0 {
		output.WriteString("## Parameters\n\n")
		for _, param := range doc.Parameters {
			output.WriteString(fmt.Sprintf("- **%s** (%s)", param.Name, param.Type))
			if param.Description != "" {
				output.WriteString(fmt.Sprintf(": %s", param.Description))
			}
			output.WriteString("\n")
		}
		output.WriteString("\n")
	}

	// Examples
	if len(doc.Examples) > 0 {
		output.WriteString("## Examples\n\n")
		for _, example := range doc.Examples {
			if example.Description != "" {
				output.WriteString(fmt.Sprintf("%s\n\n", example.Description))
			}
			output.WriteString("```\n")
			output.WriteString(fmt.Sprintf("%s\n", example.Input))
			output.WriteString("```\n\n")
		}
	}

	// Category
	if doc.Category != "" {
		output.WriteString(fmt.Sprintf("**Category:** %s\n\n", doc.Category))
	}

	// See also
	if len(doc.SeeAlso) > 0 {
		output.WriteString("**See also:** ")
		output.WriteString(strings.Join(doc.SeeAlso, ", "))
		output.WriteString("\n")
	}

	return strings.TrimSpace(output.String())
}

func formatFunctionList(docs []edlisp.FunctionDoc, title string, verbose bool) string {
	var output strings.Builder

	// Title and count
	output.WriteString(fmt.Sprintf("# %s\n\n", title))
	output.WriteString(fmt.Sprintf("Found %d function(s):\n\n", len(docs)))

	if verbose {
		// Group by category for verbose output
		categoryMap := make(map[string][]edlisp.FunctionDoc)
		for _, doc := range docs {
			category := doc.Category
			if category == "" {
				category = "Uncategorized"
			}
			categoryMap[category] = append(categoryMap[category], doc)
		}

		// Sort categories
		categories := make([]string, 0, len(categoryMap))
		for category := range categoryMap {
			categories = append(categories, category)
		}

		for _, category := range categories {
			output.WriteString(fmt.Sprintf("## %s\n\n", category))
			for _, doc := range categoryMap[category] {
				output.WriteString(fmt.Sprintf("- **%s**", doc.Name))
				if doc.Summary != "" {
					output.WriteString(fmt.Sprintf(": %s", doc.Summary))
				}
				output.WriteString("\n")
			}
			output.WriteString("\n")
		}
	} else {
		// Simple list for non-verbose output
		for _, doc := range docs {
			output.WriteString(fmt.Sprintf("- %s\n", doc.Name))
		}
	}

	return strings.TrimSpace(output.String())
}
