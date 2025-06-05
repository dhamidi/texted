package tools

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/dhamidi/texted"
	"github.com/dhamidi/texted/edlisp"
	"github.com/dhamidi/texted/edlisp/parser"
	"github.com/dhamidi/texted/edlisp/writer"
	"github.com/mark3labs/mcp-go/mcp"
)

//go:embed texted_eval_description.txt
var textedEvalDescription string

func NewTextedEvalTool() mcp.Tool {
	return NewTextedEvalToolWithPrefix("")
}

func NewTextedEvalToolWithPrefix(prefix string) mcp.Tool {
	name := "texted_eval"
	if prefix != "" {
		name = prefix + name
	}

	return mcp.NewTool(name,
		mcp.WithDescription(textedEvalDescription),
		mcp.WithString("input",
			mcp.Required(),
			mcp.Description("Input text to transform"),
		),
		mcp.WithString("script",
			mcp.Required(),
			mcp.Description("The texted script to execute on the input"),
		),
		mcp.WithString("output",
			mcp.Description("Output type: 'buffer' (default) returns transformed text, 'expression' returns last evaluated expression value"),
		),
	)
}

func TextedEvalHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input, err := request.RequireString("input")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("input parameter required: %v", err)), nil
	}

	script, err := request.RequireString("script")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("script parameter required: %v", err)), nil
	}

	outputMode := request.GetString("output", "buffer")
	if outputMode != "buffer" && outputMode != "expression" {
		return mcp.NewToolResultError("output parameter must be 'buffer' or 'expression'"), nil
	}

	if outputMode == "buffer" {
		// Use existing ExecuteScript for buffer mode
		output, err := texted.ExecuteScript(input, script)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("script execution failed: %v", err)), nil
		}
		return mcp.NewToolResultText(output), nil
	}

	// Expression mode - need to get the return value
	buf := edlisp.NewBuffer(input)

	program, err := parser.ParseString(script)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("parsing script: %v", err)), nil
	}

	env := edlisp.NewDefaultEnvironment()
	result, err := edlisp.Eval(program, env, buf)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("script execution failed: %v", err)), nil
	}

	// Format the result value as string using sexp writer
	var output strings.Builder
	sexpWriter := &writer.SExpWriter{}
	err = sexpWriter.WriteValue(&output, result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to format result: %v", err)), nil
	}

	return mcp.NewToolResultText(output.String()), nil
}
