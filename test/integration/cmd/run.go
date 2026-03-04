package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("=== Running all integration tests ===")

	tests := []struct {
		name string
		path string
	}{
		{"C#", "github.com/open-feature/cli/test/integration/cmd/csharp"},
		{"Go", "github.com/open-feature/cli/test/integration/cmd/go"},
		{"NodeJS", "github.com/open-feature/cli/test/integration/cmd/nodejs"},
		{"Angular", "github.com/open-feature/cli/test/integration/cmd/angular"},
		{"React", "github.com/open-feature/cli/test/integration/cmd/react"},
	}

	for _, test := range tests {
		fmt.Printf("--- Running %s integration test ---\n", test.name)
		cmd := exec.Command("go", "run", test.path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running %s integration test: %v\n", test.name, err)
			os.Exit(1)
		}
	}

	fmt.Println("=== All integration tests passed successfully ===")
}
