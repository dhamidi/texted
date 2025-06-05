package tools

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/dhamidi/texted"
	"github.com/mark3labs/mcp-go/mcp"
)

//go:embed edit_file_description.txt
var editFileDescription string

func NewEditFileTool() mcp.Tool {
	return NewEditFileToolWithPrefix("")
}

func NewEditFileToolWithPrefix(prefix string) mcp.Tool {
	name := "edit_file"
	if prefix != "" {
		name = prefix + name
	}

	return mcp.NewTool(name,
		mcp.WithDescription(editFileDescription),
		mcp.WithString("script",
			mcp.Required(),
			mcp.Description("The texted script to execute on each file"),
		),
		mcp.WithArray("files",
			mcp.Required(),
			mcp.Description("List of file paths to edit"),
		),
		mcp.WithBoolean("loopUntilError",
			mcp.Description("Run the script repeatedly until an error is returned"),
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

	loopUntilError := request.GetBool("loopUntilError", false)

	if len(files) == 0 {
		return mcp.NewToolResultError("at least one file must be specified"), nil
	}

	var editResults []texted.EditResult
	var editErr error
	var iterations int

	if loopUntilError {
		editResults, iterations, editErr = editFilesWithLoop(files, script)
	} else {
		editResults, editErr = texted.EditFiles(files, script)
		iterations = 1
	}

	if editErr != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to edit files: %v", editErr)), nil
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
		message := fmt.Sprintf("Completed with errors after %d iterations:\n", iterations)
		for _, result := range results {
			message += fmt.Sprintf("✓ %s\n", result)
		}
		for _, errMsg := range errors {
			message += fmt.Sprintf("✗ %s\n", errMsg)
		}
		return mcp.NewToolResultText(message), nil
	}

	message := fmt.Sprintf("All files edited successfully after %d iterations:\n", iterations)
	for _, result := range results {
		message += fmt.Sprintf("✓ %s\n", result)
	}

	return mcp.NewToolResultText(message), nil
}

// editFilesWithLoop repeatedly applies a script to files until an error occurs
func editFilesWithLoop(files []string, script string) ([]texted.EditResult, int, error) {
	iterations := 0
	var lastResults []texted.EditResult

	for {
		iterations++
		results, err := texted.EditFiles(files, script)
		if err != nil {
			return lastResults, iterations, err
		}

		// Check if any file had an error
		hasError := false
		for _, result := range results {
			if !result.Success {
				hasError = true
				break
			}
		}

		if hasError {
			return results, iterations, nil
		}

		lastResults = results
	}
}
