package integration

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

// GoBaseImage is the Go container image used to build the CLI inside
// integration test pipelines. Centralized here so the version is bumped
// in a single place when go.mod's required Go version changes.
const GoBaseImage = "golang:1.26-alpine"

// GoGenerateMinCompatImage is the minimum Go version the generated Go client
// code must compile against. Used to verify forward compatibility of the
// generate output.
const GoGenerateMinCompatImage = "golang:1.25-alpine"

// Test defines the interface for all integration tests
type Test interface {
	// Run executes the integration test with the given Dagger client and pre-built CLI container
	Run(ctx context.Context, client *dagger.Client, cli *dagger.Container) (*dagger.Container, error)
	// Name returns the name of the integration test
	Name() string
	// ProjectDir returns the absolute path to the project root
	ProjectDir() string
}

// buildOpenFeatureCLI compiles the OpenFeature CLI binary inside a Dagger container.
// The returned container has the binary at /src/cli and the project source mounted at /src.
func buildOpenFeatureCLI(client *dagger.Client, source *dagger.Directory) *dagger.Container {
	return client.Container().
		From(GoBaseImage).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"go", "mod", "download"}).
		WithExec([]string{"go", "build", "-o", "cli", "./cmd/openfeature"})
}

// RunTest builds the CLI once and runs a single integration test
func RunTest(ctx context.Context, test Test) error {
	// Initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return fmt.Errorf("failed to connect to Dagger engine: %w", err)
	}
	defer client.Close()

	fmt.Printf("=== Running %s integration test ===\n", test.Name())

	// Build the CLI once (shared across all tests)
	source := client.Host().Directory(test.ProjectDir())
	cli := buildOpenFeatureCLI(client, source)

	// Run the integration test with the pre-built CLI
	container, err := test.Run(ctx, client, cli)
	if err != nil {
		return fmt.Errorf("failed to run %s integration test: %w", test.Name(), err)
	}

	// Execute the pipeline and wait for it to complete
	_, err = container.Stdout(ctx)
	if err != nil {
		return fmt.Errorf("%s integration test failed: %w", test.Name(), err)
	}

	fmt.Printf("=== Success: %s integration test passed ===\n", test.Name())
	return nil
}
