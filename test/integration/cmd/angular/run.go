package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/open-feature/cli/test/integration"
)

type Test struct {
	ProjectDir string
	TestDir    string
}

func New(projectDir, testDir string) *Test {
	return &Test{
		ProjectDir: projectDir,
		TestDir:    testDir,
	}
}

func (t *Test) Run(ctx context.Context, client *dagger.Client) (*dagger.Container, error) {
	// Mount the project source
	source := client.Host().Directory(t.ProjectDir)

	// Mount the test files
	testFiles := client.Host().Directory(t.TestDir, dagger.HostDirectoryOpts{
		Include: []string{
			"package.json",
			"tsconfig.json",
			"vitest.config.ts",
			"setup.ts",
			"test.component.ts",
			"specs/**/*.ts",
		},
	})

	// Build the CLI in a Go container
	cli := client.Container().
		From("golang:1.24-alpine").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"go", "build", "-o", "cli", "./cmd/openfeature"})

	// Generate the Angular code
	generated := cli.WithExec([]string{
		"./cli", "generate", "angular",
		"--manifest=/src/sample/sample_manifest.json",
		"--output=/tmp/generated",
	})

	// Get the generated files
	generatedFiles := generated.Directory("/tmp/generated")

	// Create the Angular test container
	nodeContainer := client.Container().
		From("node:22-alpine").
		// Install necessary build tools for native modules
		WithExec([]string{"apk", "add", "--no-cache", "python3", "make", "g++"}).
		// Copy test files
		WithDirectory("/app", testFiles).
		// Copy generated files
		WithDirectory("/app/generated", generatedFiles).
		WithWorkdir("/app").
		// Install dependencies
		WithExec([]string{"npm", "install"}).
		// Run the tests
		WithExec([]string{"npm", "test"})

	return nodeContainer, nil
}

func (t *Test) Name() string {
	return "angular"
}

func main() {
	ctx := context.Background()

	projectDir, err := filepath.Abs(os.Getenv("PWD"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get project dir: %v\n", err)
		os.Exit(1)
	}

	testDir, err := filepath.Abs(filepath.Join(projectDir, "test/angular-integration"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get test dir: %v\n", err)
		os.Exit(1)
	}

	test := New(projectDir, testDir)

	if err := integration.RunTest(ctx, test); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
