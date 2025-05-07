# C# Integration Testing

This directory contains integration tests for the C# code generator.

## Running the tests

Run the C# integration tests with Dagger:

```bash
make test-csharp-dagger
```

This will:
1. Build the OpenFeature CLI
2. Generate C# client code using the sample manifest
3. Run the C# compilation test in an isolated environment
4. Report success or failure

## What the test does

The integration test:
1. Builds the OpenFeature CLI inside a container
2. Generates C# client code using a sample manifest
3. Compiles the generated code with a sample program
4. Runs the compiled program to verify it works correctly

## Test Files

- `CompileTest.csproj`: .NET project file for compilation testing
- `Program.cs`: Test program that uses the generated code
- `expected/`: Directory containing expected output files (used for verification)

## Implementation

The C# integration test uses Dagger to create a reproducible test environment:

1. It builds the CLI in a Go container
2. Generates C# code using the CLI
3. Tests the generated code in a .NET container

The implementation is located in `test/integration/cmd/csharp/run.go`.

For more implementation details, see the main [test/README.md](../README.md) file.