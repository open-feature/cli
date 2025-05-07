# OpenFeature CLI Integration Testing

This directory contains integration tests for validating the OpenFeature CLI generators.

## Integration Test Structure

The integration tests use [Dagger](https://dagger.io/) to create reproducible test environments without needing to install dependencies locally.

Each integration test:

1. Builds the CLI from source
2. Generates code using a sample manifest file
3. Compiles and tests the generated code in a language-specific container
4. Reports success or failure

## Running Tests

### Run all integration tests

```bash
make test-integration
```

### Run a specific integration test

```bash
# For C# tests
make test-csharp-dagger
```

## Adding a New Integration Test

To add an integration test for a new generator:

1. Create a combined implementation and runner file in `test/integration/cmd/<language>/run.go`
2. Update the main runner in `test/integration/cmd/run.go` to execute your new test
3. Add a Makefile target for running your test individually

See the step-by-step guide in [new-language.md](new-language.md) for detailed instructions.

## How It Works

The testing framework uses the following components:

- `test/integration/integration.go`: Defines the `Test` interface and common utilities
- `test/integration/cmd/run.go`: Runner for all integration tests that executes each language-specific test
- `test/integration/cmd/<language>/run.go`: Combined implementation and runner for each language
- `test/<language>-integration/`: Contains language-specific test files (code samples, project files)

Each integration test uses Dagger to:

1. Build the CLI in a clean environment
2. Generate code using a sample manifest
3. Compile and test the generated code in a language-specific container
4. Report success or failure

## Benefits Over Shell Scripts

Using Dagger for integration tests provides several advantages:

1. **Reproducibility**: Tests run in containerized environments that are identical locally and in CI
2. **Language Support**: Easy to add new language tests with the same pattern
3. **Improved Debugging**: Clear separation of build, generate, and test steps
4. **Parallelization**: Tests can run in parallel when executed in different containers
5. **No Dependencies**: No need to install language-specific tooling locally
