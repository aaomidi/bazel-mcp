package rdeps

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aaomidi/mcp-bazel/bazel"
	"github.com/mark3labs/mcp-go/mcp"
)

// Tool defines the MCP tool for finding reverse dependencies.
var Tool = mcp.NewTool("reverse-dependencies",
	mcp.WithDescription("Given a bazel target, or file path, find all other bazel targets that depend on it."),
	mcp.WithString("target",
		mcp.Required(),
		mcp.Description("The target to find reverse dependencies for."),
	),
	mcp.WithString("project_path",
		mcp.Required(),
		mcp.Description("Where MODULE.bazel or WORKSPACE is located."),
	),
	mcp.WithNumber("max_depth",
		mcp.Description("The maximum depth to search for reverse dependencies. Set to -1 to search indefinitely. 1 finds the immediate targets."),
		mcp.DefaultNumber(-1),
	),
)

// Handler implements the logic for the reverse dependencies tool.
func Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	targetInput, ok := request.Params.Arguments["target"].(string)
	if !ok {
		return nil, errors.New("target must be a string")
	}

	projectPath, ok := request.Params.Arguments["project_path"].(string)
	if !ok {
		return nil, errors.New("project_path must be a string")
	}

	// MCP sends numbers as float64, so we need to retrieve and convert
	maxDepthFloat, ok := request.Params.Arguments["max_depth"].(float64)
	if !ok {
		log.Printf("Warning: max_depth argument was not a float64, using default -1. Value received: %v", request.Params.Arguments["max_depth"])
		maxDepthFloat = -1 // Default manually if type assertion fails
	}
	maxDepth := int(maxDepthFloat)

	bazelTarget, err := bazel.ResolveInputToBazelTarget(targetInput, projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve input %q to bazel target: %w", targetInput, err)
	}

	// Build the bazel query arguments
	queryArgs := []string{"query"}
	queryScope := "//..." // Default scope for rdeps
	queryExpr := ""

	if maxDepth > 0 {
		// Use depth-limited rdeps query: rdeps(universe, target, depth)
		queryExpr = fmt.Sprintf("rdeps(%s, %s, %d)", queryScope, bazelTarget, maxDepth)
	} else {
		// Use standard rdeps query (depth <= 0 means unlimited): rdeps(universe, target)
		queryExpr = fmt.Sprintf("rdeps(%s, %s)", queryScope, bazelTarget)
	}
	queryArgs = append(queryArgs, queryExpr, "--output", "graph")

	log.Printf("Executing bazel command with args: %v in directory: [%s]", queryArgs, projectPath)

	output, err := bazel.ExecuteCommand(ctx, projectPath, queryArgs...)
	if err != nil {
		// Error from ExecuteCommand already includes details and output
		return nil, fmt.Errorf("failed to execute bazel query for target %q (derived from input %q) with depth %d: %w", bazelTarget, targetInput, maxDepth, err)
	}

	return mcp.NewToolResultText(output), nil
}
