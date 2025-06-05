package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mark3labs/mcp-go/server"

	"github.com/dhamidi/texted/tools"
)

func NewMCPCommand() *cobra.Command {
	var prefix string

	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Start an MCP server for texted",
		Long: `Start an MCP (Model Context Protocol) server that exposes texted functionality
through standardized tools. The server communicates over stdio and provides:

- edit_file: Apply texted scripts to one or more files
- texted_eval: Transform input text using texted scripts
- texted_doc: Query texted function documentation

The server supports all texted script formats: shell-like syntax, S-expressions, and JSON.

Use the --prefix flag to add a custom prefix to all tool names when registering them.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPServer(prefix)
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "Prefix to add to tool names")

	return cmd
}

func runMCPServer(prefix string) error {
	s := server.NewMCPServer(
		"Texted MCP Server",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	editFileTool := tools.NewEditFileToolWithPrefix(prefix)
	s.AddTool(editFileTool, tools.EditFileHandler)

	textedEvalTool := tools.NewTextedEvalToolWithPrefix(prefix)
	s.AddTool(textedEvalTool, tools.TextedEvalHandler)

	textedDocTool := tools.NewTextedDocToolWithPrefix(prefix)
	s.AddTool(textedDocTool, tools.TextedDocHandler)

	if err := server.ServeStdio(s); err != nil {
		return fmt.Errorf("MCP server error: %w", err)
	}

	return nil
}
