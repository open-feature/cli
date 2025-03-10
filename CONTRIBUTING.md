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
