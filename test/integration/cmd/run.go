package main

import (
	"fmt"
	"os"
	"os/exec"
)

// runIntegrationTest runs a single integration test for the specified language
func runIntegrationTest(language string) error {
	cmd := exec.Command("go", "run", fmt.Sprintf("github.com/open-feature/cli/test/integration/cmd/%s", language))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running %s integration test: %w", language, err)
	}
	return nil
}

func main() {
	// List of all integration tests to run
	tests := []string{"csharp", "go", "nodejs", "angular", "nestjs"}

	fmt.Println("=== Running all integration tests ===")

	// Run each integration test
	for _, test := range tests {
		if err := runIntegrationTest(test); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("=== All integration tests passed successfully ===")
}
