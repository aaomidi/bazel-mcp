package bazel

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveInputToBazelTarget(t *testing.T) {
	// Create a temporary directory structure for testing relative paths
	tempDir, err := os.MkdirTemp("", "bazel_target_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after test

	projectRoot := filepath.Join(tempDir, "myproject")
	if err := os.Mkdir(projectRoot, 0755); err != nil {
		t.Fatalf("Failed to create project root: %v", err)
	}
	// Create a subdirectory and a dummy file
	subdir := filepath.Join(projectRoot, "src", "app")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}
	_, err = os.Create(filepath.Join(subdir, "main.go"))
	if err != nil {
		t.Fatalf("Failed to create dummy file: %v", err)
	}
	_, err = os.Create(filepath.Join(projectRoot, "BUILD.bazel"))
	if err != nil {
		t.Fatalf("Failed to create dummy file: %v", err)
	}

	absoluteFilePath := filepath.Join(subdir, "main.go")
	absoluteRootFile := filepath.Join(projectRoot, "BUILD.bazel")

	tests := []struct {
		name        string
		targetInput string
		projectPath string
		wantTarget  string
		wantErr     bool
	}{
		{
			name:        "Valid Bazel Label",
			targetInput: "//src/app:main",
			projectPath: projectRoot,
			wantTarget:  "//src/app:main",
			wantErr:     false,
		},
		{
			name:        "Valid External Repo Label",
			targetInput: "@my_repo//pkg:lib",
			projectPath: projectRoot,
			wantTarget:  "@my_repo//pkg:lib",
			wantErr:     false,
		},
		{
			name:        "Relative Path to File",
			targetInput: filepath.Join("src", "app", "main.go"), // Use filepath.Join for OS independence
			projectPath: projectRoot,
			wantTarget:  "//src/app:main.go",
			wantErr:     false,
		},
		{
			name:        "Absolute Path to File",
			targetInput: absoluteFilePath,
			projectPath: projectRoot,
			wantTarget:  "//src/app:main.go",
			wantErr:     false,
		},
		{
			name:        "Relative Path to Root File",
			targetInput: "BUILD.bazel",
			projectPath: projectRoot,
			wantTarget:  "//:BUILD.bazel",
			wantErr:     false,
		},
		{
			name:        "Absolute Path to Root File",
			targetInput: absoluteRootFile,
			projectPath: projectRoot,
			wantTarget:  "//:BUILD.bazel",
			wantErr:     false,
		},
		{
			name:        "Path Outside Project (Relative)",
			targetInput: filepath.Join("..", "other_file"),
			projectPath: projectRoot,
			wantTarget:  "",
			wantErr:     true,
		},
		{
			name:        "Path Outside Project (Absolute)",
			targetInput: filepath.Join(tempDir, "outside.txt"), // File directly in tempDir
			projectPath: projectRoot,
			wantTarget:  "",
			wantErr:     true,
		},
		{
			name:        "Empty Target Input",
			targetInput: "",
			projectPath: projectRoot,
			wantTarget:  "",
			wantErr:     true,
		},
		{
			name:        "Empty Project Path",
			targetInput: "src/app/main.go",
			projectPath: "",
			wantTarget:  "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTarget, err := ResolveInputToBazelTarget(tt.targetInput, tt.projectPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveInputToBazelTarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTarget != tt.wantTarget {
				t.Errorf("ResolveInputToBazelTarget() gotTarget = %v, want %v", gotTarget, tt.wantTarget)
			}
		})
	}
}
