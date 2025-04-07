package bazel

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

// ResolveInputToBazelTarget takes a user input string (which can be a Bazel target label
// or a file path) and a project path, and returns the canonical Bazel target label.
// It handles resolving absolute/relative file paths against the project path.
func ResolveInputToBazelTarget(targetInput, projectPath string) (string, error) {
	if targetInput == "" {
		return "", fmt.Errorf("target input cannot be empty")
	}
	if projectPath == "" {
		return "", fmt.Errorf("project path cannot be empty")
	}

	if strings.HasPrefix(targetInput, "//") || strings.HasPrefix(targetInput, "@") {
		// Assume it's already a valid Bazel target label
		log.Printf("Input %q is already a Bazel label.", targetInput)
		return targetInput, nil
	}

	// --- Assume it's a file path ---
	log.Printf("Input %q is not a label, attempting to resolve as file path relative to project %q", targetInput, projectPath)

	var absTargetPath string
	var err error

	// Ensure projectPath is absolute and clean for reliable operations
	absProjectPath, err := filepath.Abs(projectPath)
	if err != nil {
		return "", fmt.Errorf("could not determine absolute path for project path %q: %w", projectPath, err)
	}
	absProjectPath = filepath.Clean(absProjectPath)

	// Check if the input target is already absolute
	if filepath.IsAbs(targetInput) {
		absTargetPath = filepath.Clean(targetInput)
	} else {
		// If not absolute, join it with the project path and then resolve fully
		joinedPath := filepath.Join(absProjectPath, targetInput)
		absTargetPath, err = filepath.Abs(joinedPath)
		if err != nil {
			return "", fmt.Errorf("could not resolve absolute path for relative target %q in project %q: %w", targetInput, absProjectPath, err)
		}
		absTargetPath = filepath.Clean(absTargetPath) // Clean the final absolute path
	}

	// Now calculate the path relative to the absolute project path.
	relativePath, err := filepath.Rel(absProjectPath, absTargetPath)
	if err != nil {
		// This might happen if the resolved absTargetPath isn't under absProjectPath
		return "", fmt.Errorf("target path %q (resolved from %q) could not be made relative to project path %q: %w", absTargetPath, targetInput, absProjectPath, err)
	}

	// Check if Rel produced a path starting with '..'. This means the target is outside the project dir.
	if strings.HasPrefix(relativePath, "..") || filepath.IsAbs(relativePath) {
		return "", fmt.Errorf("target path %q (resolved from %q) is outside the project directory %q", absTargetPath, targetInput, absProjectPath)
	}

	// Now 'relativePath' is confirmed relative to the project root and within it.
	// Convert to Bazel's expected format (forward slashes).
	cleanPath := filepath.ToSlash(relativePath)

	// Construct the Bazel target label //path/to:file or //:file
	dir := filepath.Dir(cleanPath)
	base := filepath.Base(cleanPath)

	var bazelTarget string
	if dir == "." {
		bazelTarget = fmt.Sprintf("//:%s", base)
	} else {
		bazelTarget = fmt.Sprintf("//%s:%s", dir, base)
	}
	log.Printf("Interpreted input %q as file path, converted to Bazel target %q relative to project %q", targetInput, bazelTarget, absProjectPath)
	return bazelTarget, nil
}
