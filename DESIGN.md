# OpenFeature CLI Design

This document describes the design considerations and goals for the OpenFeature CLI tool.

## Why Code Generation?

Code generation automates the creation of strongly typed flag accessors, minimizing configuration errors and providing a better developer experience.
By generating these accessors, developers can avoid issues related to incorrect flag names or types, resulting in more reliable and maintainable code.

Benefits of the code generation approach:

- **Type Safety**: Catch flag-related errors at compile time instead of runtime
- **IDE Support**: Get autocomplete and documentation for your flags
- **Refactoring Support**: Rename flags and the changes propagate throughout your codebase
- **Discoverability**: Make it easier for developers to find and use available flags

## Goals

- **Unified Flag Manifest Format**: Establish a standardized flag manifest format that can be easily converted from existing configurations.
- **Strongly Typed Flag Accessors**: Develop a CLI tool to generate strongly typed flag accessors for multiple programming languages.
- **Modular and Extensible Design**: Create a format that allows for future extensions and modularization of flags.
- **Language Agnostic**: Support multiple programming languages through a common flag manifest format.
- **Provider Independence**: Work with any feature flag provider that can be adapted to the OpenFeature API.
- **Remote Flag Synchronization**: Enable bidirectional synchronization with remote flag management services through a standardized API.

## Non-Goals

- **Full Provider Integration**: The initial scope does not include creating tools to convert provider-specific configurations to the new flag manifest format.
- **Validation of Flag Configs**: The project will not initially focus on validating flag configurations for consistency with the flag manifest.
- **General-Purpose Configuration**: The project will not aim to create a general-purpose configuration tool for feature flags beyond the scope of the code generation tool.
- **Runtime Flag Management**: The CLI is not intended to replace provider SDKs for runtime flag evaluation.

## Architecture Patterns

### OpenAPI Client Pattern

The CLI uses an OpenAPI-driven architecture for all remote operations. This pattern provides several benefits:

#### Benefits

1. **Type Safety**: Generated clients from OpenAPI specs ensure compile-time checking of API requests and responses
2. **Self-Documenting**: The OpenAPI specification serves as both implementation guide and documentation
3. **Provider Agnostic**: Any service implementing the Manifest Management API can integrate with the CLI
4. **Maintainability**: Changes to the API are made in one place (the spec) and propagate to all consumers
5. **Extensibility**: New endpoints and operations can be added without breaking existing functionality

#### Implementation

The pattern follows this architecture:

```
OpenAPI Spec (api/v0/sync.yaml)
    ↓
Generated Client (internal/api/client/*.gen.go)
    ↓
Wrapper Client (internal/api/sync/client.go)
    ↓
CLI Commands (internal/cmd/pull.go, push.go)
```

For detailed implementation guidelines, see the [OpenAPI Client Pattern documentation](./docs/openapi-client-pattern.md).

### Code Generation Pattern

The CLI generates strongly-typed flag accessors for multiple languages. This pattern:

1. **Parses** the flag manifest JSON/YAML
2. **Validates** the flag configurations against the schema
3. **Transforms** the data into language-specific structures
4. **Generates** code using Go templates
5. **Formats** the output according to language conventions

Each generator follows a consistent interface, making it easy to add support for new languages.
