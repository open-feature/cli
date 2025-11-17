# Contributing to OpenFeature CLI

Thank you for your interest in contributing to the OpenFeature CLI! This document provides guidelines and instructions to help you get started with contributing to the project. Whether you're fixing a bug, adding a new feature, or improving documentation, your contributions are greatly appreciated.

## Development Setup

1. **Prerequisites**:
   - Go 1.23 or later
   - Make
   - golangci-lint (will be auto-installed by make commands)

2. **Clone the repository**:
   ```bash
   git clone https://github.com/open-feature/cli.git
   cd cli
   ```

3. **Build the project**:
   ```bash
   make build
   ```

4. **Run tests**:
   ```bash
   make test
   ```

5. **Run all CI checks locally**:
   ```bash
   make ci
   ```

   This will format code, run linting, execute tests, and verify all generated files are up-to-date.

## Available Make Commands

The project includes a comprehensive Makefile with the following commands:

```bash
make help                  # Show all available commands
make build                 # Build the CLI binary
make install               # Install the CLI to your system
make lint                  # Run golangci-lint
make lint-fix              # Run golangci-lint with auto-fix
make test                  # Run unit tests
make test-integration      # Run all integration tests
make generate              # Generate all code (API clients, docs, schema)
make generate-api          # Generate API clients from OpenAPI specs
make generate-docs         # Generate documentation
make generate-schema       # Generate schema
make verify-generate       # Check if all generated files are up to date
make fmt                   # Format Go code
make ci                    # Run all CI checks locally (fmt, lint, test, verify-generate)
```

### Before Submitting a PR

Run the following command to ensure your changes will pass CI:

```bash
make ci
```

This command will:
- Format your code
- Run the linter
- Execute all tests
- Verify all generated files are up-to-date

## Contributing New Generators

We welcome contributions for new generators to extend the functionality of the OpenFeature CLI. Below are the steps to contribute a new generator:

1. **Fork the Repository**: Start by forking the repository to your GitHub account.

2. **Clone the Repository**: Clone the forked repository to your local machine.

3. **Create a New Branch**: Create a new branch for your generator. Use a descriptive name for the branch, such as `feature/add-new-generator`.

4. **Add Your Generator**: Add your generator in the appropriate directory under `/internal/generate/generators/`. For example, if you are adding a generator for Python, you might create a new directory `/internal/generate/generators/python/` and add your files there.

5. **Implement the Generator**: Implement the generator logic. Ensure that your generator follows the existing patterns and conventions used in the project. Refer to the existing generators like `/internal/generate/generators/golang` or `/internal/generate/generators/react` for examples.

6. **Write Tests**: Write tests for your generator to ensure it works as expected. Add your tests in the appropriate test directory, such as `/internal/generate/generators/python/`. Write tests for any commands you may add, too. Add your command tests in the appropriate test directory, such as `cmd/generate_test.go`.

7. **Register the Generator**: After implementing your generator, you need to register it in the CLI under the `generate` command. Follow these steps to register your generator:

   - **Create a New Command Directory**: Create a new directory under `cmd/generate` with the name of your target language. For example, if you are adding a generator for Python, create a new directory `cmd/generate/python/`.

   - **Add Command File**: In the new directory, create a file named `python.go` (replace `python` with the name of your target language). This file will define the CLI command for your generator.

   - **Implement Command**: Implement the command logic in the `python.go` file. Refer to the existing commands like `cmd/generate/golang/golang.go` or `cmd/generate/react/react.go` for examples.

   - **Register Command**: Open the `cmd/generate/generate.go` file and register your new command as a subcommand. Add an import statement for your new command package and call `Root.AddCommand(python.Cmd)` (replace `python` with the name of your target language).

8. **Update Documentation**: Update the documentation to include information about your new generator. This may include updating the README.md and any other relevant documentation files. You can run `make generate-docs` to assist with documentation updates.

9. **Commit and Push**: Commit your changes and push the new branch to your forked repository.

10. **Create a Pull Request**: Create a pull request from your new branch to the main repository. Provide a clear and detailed description of your changes, including the purpose of the new generator and any relevant information.

11. **Address Feedback**: Be responsive to feedback from the maintainers. Make any necessary changes and update your pull request as needed.

### Testing

The OpenFeature CLI includes both unit and integration tests to ensure quality and correctness.

#### Unit Tests

Run the unit tests with:

```bash
go test ./...
```

#### Integration Tests

To verify that generated code compiles correctly, the project includes integration tests. The CLI uses a Dagger-based integration testing framework to test code generation for each supported language:

```bash
# Run all integration tests
make test-integration

# Run tests for a specific language
make test-csharp-dagger
```

For more information on the integration testing framework, see [Integration Testing](./docs/integration-testing.md).

## Contributing to Remote Operations

The CLI uses an OpenAPI-driven architecture for remote flag synchronization. If you're contributing to the remote operations (pull/push commands) or API specifications:

### Modifying the OpenAPI Specification

1. **Edit the specification**: Update the OpenAPI spec at `api/v0/sync.yaml`
2. **Regenerate the client**: Run `make generate-api` to regenerate the client code
3. **Update the wrapper**: Modify `internal/api/sync/client.go` if needed
4. **Test your changes**: Add or update tests in `internal/api/sync/`

### Adding New Remote Operations

1. **Define in OpenAPI**: Add the new operation to `api/v0/sync.yaml`
2. **Regenerate**: Run `make generate-api`
3. **Implement wrapper method**: Add the method in `internal/api/sync/client.go`
4. **Create/update command**: Add or modify commands in `internal/cmd/`
5. **Add tests**: Include unit tests and integration tests
6. **Update documentation**: Update relevant documentation including command docs

### API Compatibility

When modifying the OpenAPI specification:

- **Backwards Compatibility**: Ensure changes don't break existing integrations
- **Versioning**: Use proper API versioning (currently v0)
- **Documentation**: Update the specification descriptions and examples
- **Schema Validation**: Ensure all request/response schemas are properly defined

For detailed information about the OpenAPI client pattern, see the [OpenAPI Client Pattern documentation](./docs/openapi-client-pattern.md).

## CI/CD Pipeline

The project uses GitHub Actions for continuous integration. The following workflows run automatically:

### PR Validation Workflow
- **Generated Files Check**: Ensures all generated files (OpenAPI client, docs, schema) are up-to-date
- **Format Check**: Verifies all Go code is properly formatted
- **Tests**: Runs all unit tests
- **Integration Tests**: Runs language-specific integration tests

### PR Test Workflow
- **Unit Tests**: Runs all unit tests
- **Integration Tests**: Runs Dagger-based integration tests for all supported languages

### PR Lint Workflow
- **golangci-lint**: Runs comprehensive linting with golangci-lint v1.64

## Setting Up Lefthook

To streamline the setup of Git hooks for this project, we utilize [Lefthook](https://github.com/evilmartians/lefthook). Lefthook automates pre-commit and pre-push checks, ensuring consistent enforcement of best practices across the team. These checks include code formatting, documentation generation, and running tests.

This tool is particularly helpful for new contributors or those returning to the project after some time, as it provides a seamless way to align with the project's standards. By catching issues early in your local development environment, Lefthook helps you address potential problems before opening a Pull Request, saving time and effort for both contributors and maintainers.

### Installation

1. Install Lefthook using Homebrew:

   ```bash
   brew install lefthook
   ```

2. Install the Lefthook configuration into your Git repository:

   ```bash
   lefthook install
   ```

### Pre-Commit Hook

The pre-commit hook is configured to run the following check:

1. **Code Formatting**: Ensures all files are properly formatted using `go fmt`. Any changes made by `go fmt` will be automatically staged.

### Pre-Push Hook

The pre-push hook is configured to run the following checks:

1. **Generated Files Check**: Verifies all generated files are up-to-date:
   - **OpenAPI Client**: Ensures `internal/api/client/` is current with the OpenAPI spec
   - **Documentation**: Ensures `docs/` is current with the latest command structure
   - **Schema**: Ensures `schema/` files are up-to-date
2. **Tests**: Executes `make test` to verify that all tests pass

If any of these checks fail, the push will be blocked. Run `make generate` to update all generated files and commit the changes.

### Running Hooks Manually

You can manually run the hooks using the following commands:

- Pre-commit hook:

  ```bash
  lefthook run pre-commit
  ```

- Pre-push hook:

  ```bash
  lefthook run pre-push
  ```

## Templates

### Data

The `TemplateData` struct is used to pass data to the templates.

### Built-in template functions

The following functions are automatically included in the templates:

#### ToPascal

Converts a string to `PascalCase`

```go
{{ "hello world" | ToPascal }} // HelloWorld
```

#### ToCamel

Converts a string to `camelCase`

```go
{{ "hello world" | ToCamel }} // helloWorld
```

#### ToKebab

Converts a string to `kebab-case`

```go
{{ "hello world" | ToKebab }} // hello-world
```

#### ToSnake

Converts a string to `snake_case`

```go
{{ "hello world" | ToSnake }} // hello_world
```

#### ToScreamingSnake

Converts a string to `SCREAMING_SNAKE_CASE`

```go
{{ "hello world" | ToScreamingSnake }} // HELLO_WORLD
```

#### ToUpper

Converts a string to `UPPER CASE`

```go
{{ "hello world" | ToUpper }} // HELLO WORLD
```

#### ToLower

Converts a string to `lower case`

```go
{{ "HELLO WORLD" | ToLower }} // hello world
```

#### ToTitle

Converts a string to `Title Case`

```go
{{ "hello world" | ToTitle }} // Hello World
```

#### Quote

Wraps a string in double quotes

```go
{{ "hello world" | Quote }} // "hello world"
```

#### QuoteString

Wraps only strings in double quotes

```go
{{ "hello world" | QuoteString }} // "hello world"
{{ 123 | QuoteString }} // 123
```

### Custom template functions

You can add custom template functions by passing a `FuncMap` to the `GenerateFile` function.
