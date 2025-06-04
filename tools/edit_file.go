package tools

import (
	"context"
	"fmt"

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

	var results []string
	var errors []string

	for _, filename := range files {
		content, err := readFile(filename)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to read %s: %v", filename, err))
			continue
		}

		modified, err := ExecuteScript(content, script)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to execute script on %s: %v", filename, err))
			continue
		}

		err = writeFile(filename, modified)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to write %s: %v", filename, err))
			continue
		}

		results = append(results, fmt.Sprintf("Successfully edited %s", filename))
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
