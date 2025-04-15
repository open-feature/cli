package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/open-feature/cli/test/integration"
)

// Test implements the integration test for the C# generator
type Test struct {
	// ProjectDir is the absolute path to the root of the project
	ProjectDir string
	// TestDir is the absolute path to the test directory
	TestDir string
}

// New creates a new Test
func New(projectDir, testDir string) *Test {
	return &Test{
		ProjectDir: projectDir,
		TestDir:    testDir,
	}
}

// Run executes the C# integration test using Dagger
func (t *Test) Run(ctx context.Context, client *dagger.Client) (*dagger.Container, error) {
	// Source code container
	source := client.Host().Directory(t.ProjectDir)
	testFiles := client.Host().Directory(t.TestDir, dagger.HostDirectoryOpts{
		Include: []string{"CompileTest.csproj", "Program.cs"},
	})

	// Build the CLI
	cli := client.Container().
		From("golang:1.24-alpine").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"go", "build", "-o", "cli"})

	// Generate C# client
	generated := cli.WithExec([]string{
		"./cli", "generate", "csharp",
		"--manifest=/src/sample/sample_manifest.json",
		"--output=/tmp/generated",
		"--namespace=TestNamespace",
	})

	// Get generated files
	generatedFiles := generated.Directory("/tmp/generated")

	// Test C# compilation with the generated files
	dotnetContainer := client.Container().
		From("mcr.microsoft.com/dotnet/sdk:8.0").
		WithDirectory("/app/generated", generatedFiles).
		WithDirectory("/app", testFiles).
		WithWorkdir("/app").
		WithExec([]string{"dotnet", "restore"}).
		WithExec([]string{"dotnet", "build"}).
		WithExec([]string{"dotnet", "run"})

	return dotnetContainer, nil
}

// Name returns the name of the integration test
func (t *Test) Name() string {
	return "csharp"
}

func main() {
	ctx := context.Background()

	// Get project root
	projectDir, err := filepath.Abs(os.Getenv("PWD"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get project dir: %v\n", err)
		os.Exit(1)
	}

	// Get test directory
	testDir, err := filepath.Abs(filepath.Join(projectDir, "test/csharp-integration"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get test dir: %v\n", err)
		os.Exit(1)
	}

	// Create and run the C# integration test
	test := New(projectDir, testDir)

	if err := integration.RunTest(ctx, test); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
