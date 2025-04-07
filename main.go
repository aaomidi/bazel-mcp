package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"

	"github.com/aaomidi/mcp-bazel/bazel"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func _main() error {
	s := server.NewMCPServer(
		"bazel-mcp",
		"0.0.1",
	)

	rdepsTarget := mcp.NewTool("reverse-dependencies",
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

	s.AddTool(rdepsTarget, reverseDependenciesToolHandler)

	return server.ServeStdio(s)
}

// executeCommand runs a shell command in the given directory and returns its stdout.
// If the command fails, it returns an error including stderr.
func executeCommand(ctx context.Context, workingDir string, command string) (string, error) {
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	cmd.Dir = workingDir

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("command failed: %v\nstderr: %s", err, string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to run command: %w", err)
	}
	return string(output), nil
}

func reverseDependenciesToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	var cmdStr string
	queryScope := "//..." // Default scope for rdeps

	if maxDepth > 0 {
		// Use depth-limited rdeps query: rdeps(universe, target, depth)
		cmdStr = fmt.Sprintf("bazel query 'rdeps(%s, %s, %d)' --output graph", queryScope, bazelTarget, maxDepth)
	} else {
		// Use standard rdeps query (depth <= 0 means unlimited): rdeps(universe, target)
		cmdStr = fmt.Sprintf("bazel query 'rdeps(%s, %s)' --output graph", queryScope, bazelTarget)
	}

	log.Printf("Executing command: [%s] in directory: [%s]", cmdStr, projectPath)

	output, err := executeCommand(ctx, projectPath, cmdStr)
	if err != nil {
		// Provide more context in the error message
		return nil, fmt.Errorf("failed to execute bazel query for target %q (derived from input %q) with depth %d: %w", bazelTarget, targetInput, maxDepth, err)
	}

	return mcp.NewToolResultText(output), nil
}

func main() {
	if err := _main(); err != nil {
		log.Fatalf("%+v", err)
	}
}
