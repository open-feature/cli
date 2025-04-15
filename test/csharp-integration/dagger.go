package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	// Initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to Dagger engine: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	// Get project root directory
	projectDir, err := filepath.Abs("../..")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get project dir: %v\n", err)
		os.Exit(1)
	}

	// Get integration test directory
	testDir, err := filepath.Abs(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get test dir: %v\n", err)
		os.Exit(1)
	}

	// Source code container
	source := client.Host().Directory(projectDir)

	// Build the CLI
	cli := client.Container().
		From("golang:1.21-alpine").
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
		WithDirectory("/app/expected", generatedFiles).
		WithDirectory("/app/test", client.Host().Directory(testDir, dagger.HostDirectoryOpts{
			Include: []string{"CompileTest.csproj", "Program.cs"},
		})).
		WithWorkdir("/app").
		WithExec([]string{"cp", "/app/test/CompileTest.csproj", "."}).
		WithExec([]string{"cp", "/app/test/Program.cs", "."}).
		WithExec([]string{"dotnet", "restore"}).
		WithExec([]string{"dotnet", "build"}).
		WithExec([]string{"dotnet", "run"})

	// Execute the pipeline
	_, err = dotnetContainer.ExitCode(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Pipeline failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Success: C# code compiles and executes correctly ===")
}