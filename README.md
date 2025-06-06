<!-- markdownlint-disable MD033 -->
<!-- x-hide-in-docs-start -->
<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/open-feature/community/0e23508c163a6a1ac8c0ced3e4bd78faafe627c7/assets/logo/horizontal/white/openfeature-horizontal-white.svg" />
    <img align="center" alt="OpenFeature Logo" src="https://raw.githubusercontent.com/open-feature/community/0e23508c163a6a1ac8c0ced3e4bd78faafe627c7/assets/logo/horizontal/black/openfeature-horizontal-black.svg" />
  </picture>
</p>

<h2 align="center">OpenFeature CLI</h2>
<!-- x-hide-in-docs-end -->
<!-- The 'github-badges' class is used in the docs -->
<p align="center" class="github-badges">
  <a href="https://github.com/orgs/open-feature/projects/17">
    <img alt="work-in-progress" src="https://img.shields.io/badge/status-WIP-yellow" />
  </a>
  <a href="https://cloud-native.slack.com/archives/C07DY4TUDK6">
    <img alt="Slack" src="https://img.shields.io/badge/slack-%40cncf%2Fopenfeature-brightgreen?style=flat&logo=slack" />
  </a>
</p>
<!-- x-hide-in-docs-start -->

> [!CAUTION]
> The OpenFeature CLI is experimental!
> Feel free to give it a shot and provide feedback, but expect breaking changes.

[OpenFeature](https://openfeature.dev) is an open specification that provides a vendor-agnostic, community-driven API for feature flagging that works with your favorite feature flag management tool or in-house solution.
<!-- x-hide-in-docs-end -->

## Overview

The OpenFeature CLI is a command-line tool designed to improve the developer experience when working with feature flags.
It helps developers manage feature flags consistently across different environments and programming languages by providing powerful utilities for code generation, flag validation, and more.

## Installation

### via curl

The OpenFeature CLI can be installed using a shell command.
This method is suitable for most Unix-like operating systems.

```bash
curl -fsSL https://openfeature.dev/scripts/install_cli.sh | sh
```

### via Docker

The OpenFeature CLI is available as a Docker image in the [GitHub Container Registry](https://github.com/open-feature/cli/pkgs/container/cli).

You can run the CLI in a Docker container using the following command:

```bash
docker run -it -v $(pwd):/local -w /local ghcr.io/open-feature/cli:latest
```

### via Go

If you have `Go >= 1.23` installed, you can install the CLI using the following command:

```bash
go install github.com/open-feature/cli/cmd/openfeature@latest
```

### via pre-built binaries

Download the appropriate pre-built binary from the [releases page](https://github.com/open-feature/cli/releases).

## Quick Start

1. Create a flag manifest file in your project root:

```bash
cat > flags.json << EOL
{
  "$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
  "flags": {
    "enableMagicButton": {
      "flagType": "boolean",
      "defaultValue": false,
      "description": "Activates a special button that enhances user interaction with magical, intuitive functionalities."
    }
  }
}
EOL
```

> [!NOTE]
> This is for demonstration purposes only.
> In a real-world scenario, you would typically want to fetch this file from a remote flag management service.
> See [here](https://github.com/open-feature/cli/issues/3), more more details.

2. Generate code for your preferred language:

```bash
openfeature generate react
```

See [here](./docs/commands/openfeature_generate.md) for all available options.

3. View the generated code:

```bash
cat openfeature.ts
```

**Congratulations!**
You have successfully generated your first strongly typed flag accessors.
You can now use the generated code in your application to access the feature flags.
This is just scratching the surface of what the OpenFeature CLI can do.
For more advanced usage, read on!

## Commands

The OpenFeature CLI provides the following commands:

### `init`

Initialize a new flag manifest in your project.

```bash
openfeature init
```

See [here](./docs/commands/openfeature_init.md), for all available options.

### `generate`

Generate strongly typed flag accessors for your project.

```bash
# Available languages
openfeature generate

# Basic usage
openfeature generate [language]

# With custom output directory
openfeature generate typescript --output ./src/flags
```

See [here](./docs/commands/openfeature_generate.md), for all available options.

### `version`

Print the version number of the OpenFeature CLI.

```bash
openfeature version
```

See [here](./docs/commands/openfeature_version.md), for all available options.

## Flag Manifest

The flag manifest is a JSON file that defines your feature flags and their properties.
It serves as the source of truth for your feature flags and is used by the CLI to generate strongly typed accessors.
The manifest file should be named `flags.json` and placed in the root of your project.

### Flag Manifest Structure

The flag manifest file should follow the JSON schema defined [here](https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json).

The schema defines the following properties:

- `$schema`: The URL of the JSON schema for validation.
- `flags`: An object containing the feature flags.
  - `flagKey`: A unique key for the flag.
    - `description`: A description of what the flag does.
    - `type`: The type of the flag (e.g., `boolean`, `string`, `number`, `object`).
    - `defaultValue`: The default value of the flag.

### Example Flag Manifest

```json
{
  "$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
  "flags": {
    "uniqueFlagKey": {
      "description": "Description of what this flag does",
      "type": "boolean|string|number|object",
      "defaultValue": "default-value",
    }
  }
}
```

## Configuration

The OpenFeature CLI uses an optional configuration file to override default settings and customize the behavior of the CLI.
This file can be in JSON or YAML format and should be named either `.openfeature.json` or `.openfeature.yaml`.

### Configuration File Structure

```yaml
# Example .openfeature.yaml
manifest: "flags/manifest.json" # Overrides the default manifest path
generate:
  output: "src/flags" # Overrides the default output directory
  # Any language-specific options can be specified here
  # For example, for React:
  react:
    output: "src/flags/react" # Overrides the default React output directory
  # For Go:
  go:
    package: "github.com/myorg/myrepo/flags" # Overrides the default Go package name
    output: "src/flags/go" # Overrides the default Go output directory
```

### Configuration Priority

The CLI uses a layered approach to configuration, allowing you to override settings at different levels.
The configuration is applied in the following order:

```mermaid
flowchart LR
  default("Default Config")
  config("Config File")
  args("Command Line Args")
  default --> config
  config --> args
```

### Get Involved

- **CNCF Slack**: Join the conversation in the [#openfeature](https://cloud-native.slack.com/archives/C0344AANLA1) and [#openfeature-cli](https://cloud-native.slack.com/archives/C07DY4TUDK6) channel
- **Regular Meetings**: Attend our [community calls](https://zoom-lfx.platform.linuxfoundation.org/meetings/openfeature)
- **GitHub Issues**: Report bugs or request features in our [issue tracker](https://github.com/open-feature/cli/issues)
- **Social Media**:
  - Twitter: [@openfeature](https://twitter.com/openfeature)
  - LinkedIn: [OpenFeature](https://www.linkedin.com/company/openfeature/)

For more information, visit our [community page](https://openfeature.dev/community/).

### Support the project

- Give this repo a ⭐️!
- Share your experience and contribute back to the project

### Thanks to everyone who has already contributed

<a href="https://github.com/open-feature/cli/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=open-feature/cli" alt="Pictures of the folks who have contributed to the project" />
</a>

Made with [contrib.rocks](https://contrib.rocks).
<!-- x-hide-in-docs-end -->
