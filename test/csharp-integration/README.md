# C# Integration Testing with Dagger

This directory contains integration tests for the C# code generator.

## Running the tests

You can run the C# integration tests in two ways:

### Using Dagger (Recommended)

```bash
make test-csharp-dagger
```

This will:
1. Build the OpenFeature CLI
2. Generate C# client code using the sample manifest
3. Run the C# compilation test in an isolated environment
4. Report success or failure

### Using the legacy shell script

```bash
make test-csharp
```

or directly:

```bash
./test/csharp-integration/test-compilation.sh
```

The shell script version uses Docker directly and will be deprecated in the future.

## What the test does

The integration test:
1. Builds the OpenFeature CLI
2. Generates C# client code using a sample manifest
3. Compiles the generated code with a sample program
4. Runs the compiled program to verify it works correctly

## Implementation details

The Dagger implementation (`dagger.go`) creates a pipeline that:
1. Builds the CLI in a Go container
2. Generates C# code using the CLI
3. Compiles and runs the C# code in a .NET SDK container
4. Reports success based on the exit code