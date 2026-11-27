package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/open-feature/cli/test/integration"
)

const bazeliskVersion = "v1.26.0"

// Test implements the integration test for the C++ generator.
type Test struct {
	projectDir string
	TestDir    string
}

// New creates a new Test.
func New(projectDir, testDir string) *Test {
	return &Test{
		projectDir: projectDir,
		TestDir:    testDir,
	}
}

// Run executes the C++ integration test using Dagger.
func (t *Test) Run(ctx context.Context, client *dagger.Client, cli *dagger.Container) (*dagger.Container, error) {
	testFiles := client.Host().Directory(t.TestDir, dagger.HostDirectoryOpts{
		Include: []string{
			".bazelrc",
			".bazelversion",
			".clang-format",
			"BUILD.bazel",
			"MODULE.bazel",
			"main.cpp",
		},
	})

	generated := cli.WithExec([]string{
		"./cli", "generate", "cpp",
		"--manifest=/src/sample/sample_manifest.json",
		"--output=/tmp/generated",
		"--namespace=testnamespace",
	})

	generatedFiles := generated.Directory("/tmp/generated")

	installBazeliskCmd := fmt.Sprintf(
		`ARCH=$(dpkg --print-architecture) && curl -fsSL "https://github.com/bazelbuild/bazelisk/releases/download/%s/bazelisk-linux-${ARCH}" -o /usr/local/bin/bazelisk && chmod +x /usr/local/bin/bazelisk`,
		bazeliskVersion,
	)

	cppContainer := client.Container().
		From("debian:bookworm-slim").
		WithExec([]string{
			"sh", "-c",
			"apt-get update && apt-get install -y --no-install-recommends ca-certificates clang-format curl g++ git && rm -rf /var/lib/apt/lists/*",
		}).
		WithExec([]string{"sh", "-c", installBazeliskCmd}).
		WithDirectory("/app", testFiles).
		WithDirectory("/app/generated", generatedFiles).
		WithWorkdir("/app").
		WithExec([]string{"clang-format", "--style=file", "--dry-run", "-Werror", "generated/openfeature.h"}).
		WithExec([]string{"bazelisk", "test", "--test_output=errors", "//..."})

	return cppContainer, nil
}

// Name returns the name of the integration test.
func (t *Test) Name() string {
	return "cpp"
}

// ProjectDir returns the absolute path to the project root.
func (t *Test) ProjectDir() string {
	return t.projectDir
}

func main() {
	ctx := context.Background()

	projectDir, err := filepath.Abs(os.Getenv("PWD"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get project dir: %v\n", err)
		os.Exit(1)
	}

	testDir, err := filepath.Abs(filepath.Join(projectDir, "test/cpp-integration"))
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
