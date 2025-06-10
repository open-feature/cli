package main

import (
	"context"
	"dagger.io/dagger"
	"fmt"
	"github.com/open-feature/cli/test/integration"
	"os"
	"path/filepath"
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
	source := client.Host().Directory(t.ProjectDir)
	testFiles := client.Host().Directory(t.TestDir, dagger.HostDirectoryOpts{
		Include: []string{"package.json", "test.ts"},
	})

	cli := client.Container().
		From("golang:1.24-alpine").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"go", "build", "-o", "cli", "./cmd/openfeature"})

	generated := cli.WithExec([]string{
		"./cli", "generate", "nodejs",
		"--manifest=/src/sample/sample_manifest.json",
		"--output=/tmp/generated",
	})

	generatedFiles := generated.Directory("/tmp/generated")

	nodeContainer := client.Container().
		From("node:22-alpine").
		WithExec([]string{"npm", "install", "-g", "typescript"}).
		WithDirectory("/app/generated", generatedFiles).
		WithDirectory("/app", testFiles).
		WithWorkdir("/app").
		WithExec([]string{"npm", "install"}).
		WithExec([]string{"npm", "test"})

	return nodeContainer, nil
}
func (t *Test) Name() string {
	return "nodejs"
}
func main() {
	ctx := context.Background()

	projectDir, err := filepath.Abs(os.Getenv("PWD"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get project dir: %v\n", err)
		os.Exit(1)
	}

	testDir, err := filepath.Abs(filepath.Join(projectDir, "test/nodejs-integration"))
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
