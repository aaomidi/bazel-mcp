package deps

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aaomidi/mcp-bazel/bazel"
	"github.com/mark3labs/mcp-go/mcp"
)

// Tool defines the MCP tool for finding dependencies of a Bazel target.
var Tool = mcp.NewTool("deps",
	mcp.WithDescription("Finds the dependencies of a given Bazel target."),
	mcp.WithString("target",
		mcp.Required(),
		mcp.Description("The Bazel target to find dependencies for (e.g., //foo:bar)."),
	),
	mcp.WithString("project_path",
		mcp.Required(),
		mcp.Description("Where MODULE.bazel or WORKSPACE is located."),
	),
	mcp.WithNumber("depth",
		mcp.Description("The maximum depth for dependency search (default: 1 for direct deps). Must be non-negative."),
		mcp.DefaultNumber(1),
	),
)

// Handler implements the logic for the deps tool.
func Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract Arguments
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

	// Extract depth, defaulting to 1
	depthFloat, ok := request.Params.Arguments["depth"].(float64)
	if !ok {
		log.Printf("Warning: depth argument was not a float64, using default 1. Value received: %v (%T)", request.Params.Arguments["depth"], request.Params.Arguments["depth"])
		depthFloat = 1 // Default manually if type assertion fails or not present
	}
	depth := int(depthFloat)

	// Validate depth
	if depth < 0 {
		return nil, fmt.Errorf("depth argument cannot be negative, got %d", depth)
	}
	if depthFloat != float64(depth) {
		log.Printf("Warning: depth received non-integer number %f, using truncated value %d.", depthFloat, depth)
		// Allow truncated non-negative values, consistent with bazel query behavior.
	}

	// Execute Logic
	queryExpr := fmt.Sprintf("deps('%s', %d)", target, depth)
	queryArgs := []string{"query", queryExpr, "--output", "label"}

	log.Printf("Executing bazel command with args: %v in directory: [%s]", queryArgs, projectPath)

	output, err := bazel.ExecuteCommand(ctx, projectPath, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute bazel query for dependencies of target %q with depth %d: %w", target, depth, err)
	}

	// Return the list of dependency targets (one per line).
	return mcp.NewToolResultText(output), nil
}
