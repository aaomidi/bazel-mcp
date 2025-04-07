package bazel

import (
	"context"
	"fmt"
	"os/exec"
)

// ExecuteCommand runs a bazel command in the given directory and returns its combined stdout/stderr.
// The first argument in 'args' should be the bazel command (e.g., "query", "build", "test").
// If the command fails, it returns an error including the combined output.
func ExecuteCommand(ctx context.Context, workingDir string, args ...string) (string, error) {
	// Command is assumed to be "bazel"
	cmd := exec.CommandContext(ctx, "bazel", args...)
	cmd.Dir = workingDir

	output, err := cmd.CombinedOutput() // Use CombinedOutput to get stdout and stderr
	if err != nil {
		// The error already includes command details and potentially exit code.
		// CombinedOutput includes stderr in the output on error.
		return string(output), fmt.Errorf("bazel command failed: %w\nArgs: %v\nOutput:\n%s", err, args, string(output))
	}
	return string(output), nil
}
