# Adding a New Generator Integration Test

This guide explains how to add integration tests for a new generator.

## Directory Structure

The integration testing framework has the following directory structure:

```
test/
  integration/              # Core integration test framework
    integration.go          # Test interface definition
    cmd/                    # Command-line runners and implementations
      run.go                # Runner for all tests
      csharp/               # C# specific implementation and runner
        run.go
      python/               # Python specific implementation and runner (future)
        run.go
  csharp-integration/       # C# test files
  python-integration/       # Python test files (future)
```

## Step 1: Create a generator-specific implementation and runner

Create a file at `test/integration/cmd/python/run.go`:

```go
package main

import (
 "context"
 "fmt"
 "os"
 "path/filepath"
 
 "dagger.io/dagger"
 "github.com/open-feature/cli/test/integration"
)

// Test implements the integration test for the Python generator
type Test struct {
 ProjectDir string
 TestDir    string
}

// New creates a new Test
func New(projectDir, testDir string) *Test {
 return &Test{
  ProjectDir: projectDir,
  TestDir:    testDir,
 }
}

// Run executes the Python integration test
func (t *Test) Run(ctx context.Context, client *dagger.Client) (*dagger.Container, error) {
 // Source code container
 source := client.Host().Directory(t.ProjectDir)
 testFiles := client.Host().Directory(t.TestDir, dagger.HostDirectoryOpts{
  Include: []string{"test_openfeature.py", "requirements.txt"},
 })

 // Build the CLI
 cli := client.Container().
  From("golang:1.24-alpine").
  WithDirectory("/src", source).
  WithWorkdir("/src").
  WithExec([]string{"go", "build", "-o", "cli"})

 // Generate Python client
 generated := cli.WithExec([]string{
  "./cli", "generate", "python",
  "--manifest=/src/sample/sample_manifest.json",
  "--output=/tmp/generated",
  "--package=openfeature_test",
 })

 // Get generated files
 generatedFiles := generated.Directory("/tmp/generated")

 // Test Python with the generated files
 pythonContainer := client.Container().
  From("python:3.11-slim").
  WithDirectory("/app/openfeature", generatedFiles).
  WithDirectory("/app/test", testFiles).
  WithWorkdir("/app").
  WithExec([]string{"pip", "install", "-r", "test/requirements.txt"}).
  WithExec([]string{"python", "-m", "pytest", "test/test_openfeature.py", "-v"})

 return pythonContainer, nil
}

// Name returns the name of the integration test
func (t *Test) Name() string {
 return "python"
}

func main() {
 ctx := context.Background()

 // Get project root
 projectDir, err := filepath.Abs(os.Getenv("PWD"))
 if err != nil {
  fmt.Fprintf(os.Stderr, "Failed to get project dir: %v\n", err)
  os.Exit(1)
 }

 // Get test directory
 testDir, err := filepath.Abs(filepath.Join(projectDir, "test/python-integration"))
 if err != nil {
  fmt.Fprintf(os.Stderr, "Failed to get test dir: %v\n", err)
  os.Exit(1)
 }

 // Create and run the Python integration test
 test := New(projectDir, testDir)

 if err := integration.RunTest(ctx, test); err != nil {
  fmt.Fprintf(os.Stderr, "Error: %v\n", err)
  os.Exit(1)
 }
}
```

## Step 2: Add the test to the all-integration runner

Update `test/integration/cmd/run.go` to include your test:

```go
package main

import (
 "fmt"
 "os"
 "os/exec"
)

func main() {
 // Run the generator-specific tests
 fmt.Println("=== Running all integration tests ===")
 
 // Run the C# integration test
 csharpCmd := exec.Command("go", "run", "github.com/open-feature/cli/test/integration/cmd/csharp")
 csharpCmd.Stdout = os.Stdout
 csharpCmd.Stderr = os.Stderr
 if err := csharpCmd.Run(); err != nil {
  fmt.Fprintf(os.Stderr, "Error running C# integration test: %v\n", err)
  os.Exit(1)
 }
 
 // Run the Python integration test
 pythonCmd := exec.Command("go", "run", "github.com/open-feature/cli/test/integration/cmd/python")
 pythonCmd.Stdout = os.Stdout
 pythonCmd.Stderr = os.Stderr
 if err := pythonCmd.Run(); err != nil {
  fmt.Fprintf(os.Stderr, "Error running Python integration test: %v\n", err)
  os.Exit(1)
 }
 
 // Add more tests here as they are available
 
 fmt.Println("=== All integration tests passed successfully ===")
}
```

## Step 3: Create test files

Create the following directory structure with your test files:

```
test/
  python-integration/
    requirements.txt
    test_openfeature.py
    README.md
```

## Step 4: Add a Makefile target

Update the Makefile with a new target:

```makefile
.PHONY: test-python-dagger
test-python-dagger:
 @echo "Running Python integration test with Dagger..."
 @go run ./test/integration/cmd/python/run.go
```

## Step 5: Update the documentation

Update `test/README.md` to include your new test.