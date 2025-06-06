package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {

	// Run the language-specific tests
	fmt.Println("=== Running all integration tests ===")

	// Run the C# integration test
	csharpCmd := exec.Command("go", "run", "github.com/open-feature/cli/test/integration/cmd/csharp")
	csharpCmd.Stdout = os.Stdout
	csharpCmd.Stderr = os.Stderr
	if err := csharpCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running C# integration test: %v\n", err)
		os.Exit(1)
	}

	// Run the Go integration test
	goCmd := exec.Command("go", "run", "github.com/open-feature/cli/test/integration/cmd/go")
	goCmd.Stdout = os.Stdout
	goCmd.Stderr = os.Stderr
	if err := goCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running Go integration test: %v\n", err)
		os.Exit(1)
	}

	// Add more tests here as they are available

	fmt.Println("=== All integration tests passed successfully ===")
}
