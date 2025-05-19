package integration

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

// Test defines the interface for all integration tests
type Test interface {
	// Run executes the integration test with the given Dagger client
	Run(ctx context.Context, client *dagger.Client) (*dagger.Container, error)
	// Name returns the name of the integration test
	Name() string
}

// RunTest runs a single integration test
func RunTest(ctx context.Context, test Test) error {
	// Initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return fmt.Errorf("failed to connect to Dagger engine: %w", err)
	}
	defer client.Close()

	fmt.Printf("=== Running %s integration test ===\n", test.Name())

	// Run the integration test
	container, err := test.Run(ctx, client)
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
