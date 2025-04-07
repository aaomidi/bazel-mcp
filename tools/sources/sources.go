package sources

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aaomidi/mcp-bazel/bazel"
	"github.com/mark3labs/mcp-go/mcp"
)

// Tool defines the MCP tool for finding source files of a Bazel target.
var Tool = mcp.NewTool("sources",
	mcp.WithDescription("Finds the direct source files associated with a given Bazel target."),
	mcp.WithString("target",
		mcp.Required(),
		mcp.Description("The Bazel target to find sources for (e.g., //foo:bar)."),
	),
	mcp.WithString("project_path",
		mcp.Required(),
		mcp.Description("Where MODULE.bazel or WORKSPACE is located."),
	),
)

// Handler implements the logic for the sources tool.
func Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// --- Extract Arguments --- (Direct extraction, no common helpers)
	target, ok := request.Params.Arguments["target"].(string)
	if !ok {
		return nil, errors.New("target must be a string")
	}
	if target == "" {
		return nil, errors.New("target cannot be empty")
	}

	projectPath, ok := request.Params.Arguments["project_path"].(string)
	if !ok {
		return nil, errors.New("project_path must be a string")
	}
	if projectPath == "" {
		return nil, errors.New("project_path cannot be empty")
	}

	// --- Execute Logic ---
	queryExpr := fmt.Sprintf("kind('source file', deps('%s'))", target)
	queryArgs := []string{"query", queryExpr, "--output", "label"}

	log.Printf("Executing bazel command with args: %v in directory: [%s]", queryArgs, projectPath)

	output, err := bazel.ExecuteCommand(ctx, projectPath, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute bazel query for sources of target %q: %w", target, err)
	}

	// Return the list of source files (one per line, as label output does).
	return mcp.NewToolResultText(output), nil
}
