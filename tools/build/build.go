package build

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aaomidi/mcp-bazel/bazel"
	"github.com/mark3labs/mcp-go/mcp"
)

// Tool defines the MCP tool for building Bazel targets.
var Tool = mcp.NewTool("build",
	mcp.WithDescription("Builds a given Bazel target."),
	mcp.WithString("target",
		mcp.Required(),
		mcp.Description("The Bazel target to build (e.g., //foo:bar)."),
	),
	mcp.WithString("project_path",
		mcp.Required(),
		mcp.Description("Where MODULE.bazel or WORKSPACE is located."),
	),
)

// Handler implements the logic for the build tool.
func Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	target, ok := request.Params.Arguments["target"].(string)
	if !ok {
		return nil, errors.New("target must be a string")
	}
	// Basic validation: Ensure target isn't empty (more robust validation could be added)
	if target == "" {
		return nil, errors.New("target cannot be empty")
	}

	projectPath, ok := request.Params.Arguments["project_path"].(string)
	if !ok {
		return nil, errors.New("project_path must be a string")
	}

	// Build the bazel build arguments
	buildArgs := []string{"build", target}

	log.Printf("Executing bazel command with args: %v in directory: [%s]", buildArgs, projectPath)

	output, err := bazel.ExecuteCommand(ctx, projectPath, buildArgs...)
	if err != nil {
		// Error from ExecuteCommand includes command output which is useful for build failures
		return nil, fmt.Errorf("failed to execute bazel build for target %q: %w", target, err)
	}

	// On success, the output might contain build information or be empty.
	return mcp.NewToolResultText(output), nil
}
