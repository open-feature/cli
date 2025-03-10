# Generators

This directory contains the code generators for different programming languages. Each generator is responsible for generating code based on the OpenFeature flag manifest.

## Structure

Each generator should be placed in its own directory under `/internal/generators`. The directory should be named after the target language (e.g., `golang`, `react`).

Each generator directory should contain the following files:

- `language.go`: This file contains the implementation of the generator logic for the target language. Replace `language` with the name of the target language (e.g., `golang.go`, `react.go`).
- `language.tmpl`: This file contains the template used by the generator to produce the output code. Replace `language` with the name of the target language (e.g., `golang.tmpl`, `react.tmpl`).

## How Generators Work

Each generator consists of two main components: the `language.go` file and the `language.tmpl` file. The `language.go` file contains the logic for processing the feature flag manifest and generating the output code, while the `language.tmpl` file defines the template used to produce the final code.

### `language.go`

The `language.go` file is responsible for reading the feature flag manifest, processing the data, and applying it to the template defined in the `language.tmpl` file. This file typically includes functions for parsing the manifest, preparing the data for the template, and writing the generated code to the appropriate output files.

### `language.tmpl`

The `language.tmpl` file is a text template that defines the structure of the generated code. It uses the Go template syntax to insert data from the feature flag manifest into the appropriate places in the template. The `language.go` file processes this template and fills in the data to produce the final code.

### Example Workflow

1. The `language.go` file reads the feature flag manifest and parses the data.
2. The data is processed and prepared for the template.
3. The `language.go` file applies the data to the `language.tmpl` file using the Go template engine.
4. The generated code is written to the appropriate output files.

By following this pattern, you can create generators for different programming languages that produce consistent and reliable code based on the feature flag manifest.

## Example

Here is an example structure for a Go generator:

```
/internal/generators/
  golang/
    golang.go
    golang.tmpl
```

## Adding a New Generator

To add a new generator, follow these steps:

1. Create a new directory under `/internal/generators` with the name of the target language.
2. Add the `language.go` and `language.tmpl` files to the new directory.
3. Implement the generator logic in the `language.go` file.
4. Create the template in the `language.tmpl` file.
5. Ensure that your generator follows the existing patterns and conventions used in the project.
6. Write tests for your generator to ensure it works as expected.
7. Update the documentation to include information about your new generator.

We appreciate your contributions and look forward to seeing your new generators!