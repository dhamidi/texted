package tools

import (
	"context"
	"fmt"

	"github.com/dhamidi/texted"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewEditFileTool() mcp.Tool {
	return mcp.NewTool("edit_file",
		mcp.WithDescription("Apply a texted script to one or more files"),
		mcp.WithString("script",
			mcp.Required(),
			mcp.Description("The texted script to execute on each file"),
		),
		mcp.WithArray("files",
			mcp.Required(),
			mcp.Description("List of file paths to edit"),
		),
	)
}

func EditFileHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	script, err := request.RequireString("script")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("script parameter required: %v", err)), nil
	}

	files, err := request.RequireStringSlice("files")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("files parameter required: %v", err)), nil
	}

	if len(files) == 0 {
		return mcp.NewToolResultError("at least one file must be specified"), nil
	}

	editResults, err := texted.EditFiles(files, script)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to edit files: %v", err)), nil
	}

	var results []string
	var errors []string

	for _, result := range editResults {
		if result.Success {
			results = append(results, fmt.Sprintf("Successfully edited %s", result.Filename))
		} else {
			errors = append(errors, fmt.Sprintf("Failed to edit %s: %v", result.Filename, result.Error))
		}
	}

	if len(errors) > 0 {
		message := fmt.Sprintf("Completed with errors:\n")
		for _, result := range results {
			message += fmt.Sprintf("✓ %s\n", result)
		}
		for _, errMsg := range errors {
			message += fmt.Sprintf("✗ %s\n", errMsg)
		}
		return mcp.NewToolResultText(message), nil
	}

	message := "All files edited successfully:\n"
	for _, result := range results {
		message += fmt.Sprintf("✓ %s\n", result)
	}

	return mcp.NewToolResultText(message), nil
}
