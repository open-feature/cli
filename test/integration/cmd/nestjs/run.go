package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/open-feature/cli/test/integration"
)

// Test implements the integration test for the NestJS generator
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

// Run executes the NestJS integration test using Dagger
func (t *Test) Run(ctx context.Context, client *dagger.Client) (*dagger.Container, error) {
	// Source code container
	source := client.Host().Directory(t.ProjectDir)
	testFiles := client.Host().Directory(t.TestDir, dagger.HostDirectoryOpts{
		Include: []string{"package.json", "tsconfig.json", "src/**/*.ts"},
	})

	// Build the CLI
	cli := client.Container().
		From("golang:1.24-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"go", "mod", "download"}).
		WithExec([]string{"go", "build", "-o", "cli", "./cmd/openfeature"})

	// Generate NestJS client
	generated := cli.WithExec([]string{
		"./cli", "generate", "nestjs",
		"--manifest=/src/sample/sample_manifest.json",
		"--output=/tmp/generated",
	})

	// Get generated files
	generatedFiles := generated.Directory("/tmp/generated")

	// Test NestJS compilation with the generated files
	nestjsContainer := client.Container().
		From("node:20-alpine").
		WithWorkdir("/app").
		WithDirectory("/app", testFiles).
		WithDirectory("/app/src/generated", generatedFiles).
		WithExec([]string{"npm", "install"}).
		WithExec([]string{"npm", "run", "build"}).
		WithExec([]string{"node", "dist/main.js"})

	return nestjsContainer, nil
}

// Name returns the name of the integration test
func (t *Test) Name() string {
	return "nestjs"
}

func main() {
	ctx := context.Background()

	// Get project root
	projectDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get project dir: %v\n", err)
		os.Exit(1)
	}

	// Get test directory
	testDir, err := filepath.Abs(filepath.Join(projectDir, "test/nestjs-integration"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get test dir: %v\n", err)
		os.Exit(1)
	}

	// Create and run the NestJS integration test
	test := New(projectDir, testDir)

	if err := integration.RunTest(ctx, test); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
