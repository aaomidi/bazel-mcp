package test

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aaomidi/mcp-bazel/bazel"
	"github.com/mark3labs/mcp-go/mcp"
)

// Tool defines the MCP tool for testing Bazel targets.
var Tool = mcp.NewTool("test",
	mcp.WithDescription("Runs tests for a given Bazel target."),
	mcp.WithString("target",
		mcp.Required(),
		mcp.Description("The Bazel target to test (e.g., //foo:test, //path/to/tests/...)."),
	),
	mcp.WithString("project_path",
		mcp.Required(),
		mcp.Description("Where MODULE.bazel or WORKSPACE is located."),
	),
)

// Handler implements the logic for the test tool.
func Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// Build the bazel test arguments
	testArgs := []string{"test", target}

	log.Printf("Executing bazel command with args: %v in directory: [%s]", testArgs, projectPath)

	output, err := bazel.ExecuteCommand(ctx, projectPath, testArgs...)
	if err != nil {
		// Error from ExecuteCommand includes command output which is useful for test failures
		return nil, fmt.Errorf("failed to execute bazel test for target %q: %w", target, err)
	}

	// On success, the output contains the test summary.
	return mcp.NewToolResultText(output), nil
}
