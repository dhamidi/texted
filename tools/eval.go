package tools

import (
	"context"
	"fmt"

	"github.com/dhamidi/texted"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewEvalTool() mcp.Tool {
	return mcp.NewTool("eval",
		mcp.WithDescription("Transform input text using a texted script"),
		mcp.WithString("input",
			mcp.Required(),
			mcp.Description("Input text to transform"),
		),
		mcp.WithString("script",
			mcp.Required(),
			mcp.Description("The texted script to execute on the input"),
		),
	)
}

func EvalHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := request.RequireString("input")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("input parameter required: %v", err)), nil
	}

	script, err := request.RequireString("script")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("script parameter required: %v", err)), nil
	}

	output, err := texted.ExecuteScript(input, script)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("script execution failed: %v", err)), nil
	}

	return mcp.NewToolResultText(output), nil
}
