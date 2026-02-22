package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/open-feature/cli/test/integration"
)

// Test implements the integration test for the Java generator
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

// Run executes the Java integration test using Dagger
func (t *Test) Run(ctx context.Context, client *dagger.Client) (*dagger.Container, error) {
	// Source code container
	source := client.Host().Directory(t.ProjectDir)
	testFiles := client.Host().Directory(t.TestDir, dagger.HostDirectoryOpts{
		Include: []string{"pom.xml", "src/**/*.java"},
	})

	// Build the CLI
	cli := client.Container().
		From("golang:1.24-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"go", "mod", "tidy"}).
		WithExec([]string{"go", "mod", "download"}).
		WithExec([]string{"go", "build", "-o", "cli", "./cmd/openfeature"})

	// Generate Java client
	generated := cli.WithExec([]string{
		"./cli", "generate", "java",
		"--manifest=/src/sample/sample_manifest.json",
		"--output=/tmp/generated",
		"--package-name=dev.openfeature.generated",
	})

	// Get generated files
	generatedFiles := generated.Directory("/tmp/generated")

	// Test Java compilation with the generated files
	javaContainer := client.Container().
		From("maven:3.9-eclipse-temurin-21-alpine").
		WithWorkdir("/app").
		WithDirectory("/app", testFiles).
		WithDirectory("/app/src/main/java/dev/openfeature/generated", generatedFiles).
		WithExec([]string{"mvn", "clean", "compile", "-B", "-q"}).
		WithExec([]string{"mvn", "exec:java", "-Dexec.mainClass=dev.openfeature.Main", "-q"})

	return javaContainer, nil
}

// Name returns the name of the integration test
func (t *Test) Name() string {
	return "java"
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
	testDir := filepath.Join(projectDir, "test/java-integration")

	// Create and run the Java integration test
	test := New(projectDir, testDir)

	if err := integration.RunTest(ctx, test); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
