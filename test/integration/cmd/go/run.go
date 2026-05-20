package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/open-feature/cli/test/integration"
)

// Test implements the integration test for the Go generator
type Test struct {
	// projectDir is the absolute path to the root of the project
	projectDir string
	// TestDir is the absolute path to the test directory
	TestDir string
}

// New creates a new Test
func New(projectDir, testDir string) *Test {
	return &Test{
		projectDir: projectDir,
		TestDir:    testDir,
	}
}

// Run executes the Go integration test using Dagger
func (t *Test) Run(ctx context.Context, client *dagger.Client, cli *dagger.Container) (*dagger.Container, error) {
	// Test source files
	testFiles := client.Host().Directory(t.TestDir, dagger.HostDirectoryOpts{
		Include: []string{"test.go", "go.mod"},
	})

	// Generate Go client using the pre-built CLI
	generated := cli.WithExec([]string{
		"./cli", "generate", "go",
		"--manifest=/src/sample/sample_manifest.json",
		"--output=/tmp/generated",
		"--package-name=openfeature",
	})

	// Get generated files
	generatedFiles := generated.Directory("/tmp/generated")

	// Test Go compilation with the generated files
	goContainer := client.Container().
		From(integration.GoGenerateMinCompatImage).
		WithWorkdir("/app").
		WithDirectory("/app", testFiles).
		WithDirectory("/app/openfeature", generatedFiles).
		WithExec([]string{"go", "mod", "tidy"}).
		WithExec([]string{"go", "build", "-o", "test", "-v"}).
		WithExec([]string{"./test"})

	return goContainer, nil
}

// Name returns the name of the integration test
func (t *Test) Name() string {
	return "go"
}

// ProjectDir returns the absolute path to the project root
func (t *Test) ProjectDir() string {
	return t.projectDir
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
	testDir, err := filepath.Abs(filepath.Join(projectDir, "test/go-integration"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get test dir: %v\n", err)
		os.Exit(1)
	}

	// Create and run the Go integration test
	test := New(projectDir, testDir)

	if err := integration.RunTest(ctx, test); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
