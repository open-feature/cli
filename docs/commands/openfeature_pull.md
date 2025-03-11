<!-- markdownlint-disable-file -->
<!-- WARNING: THIS DOC IS AUTO-GENERATED. DO NOT EDIT! -->
## openfeature pull

Pull a flag manifest from a remote source

### Synopsis

Pull a flag manifest from a remote source.

This command fetches feature flag configurations from a specified remote source and saves them locally as a manifest file.

Supported URL schemes:
- http:// - HTTP remote sources
- https:// - HTTPS remote sources  
- file:// - Local file paths

How it works:
1. Connects to the specified flag source URL
2. Downloads the flag configuration data
3. Validates and processes each flag definition
4. Prompts for missing default values (unless --no-prompt is used)
5. Writes the complete manifest to the local file system

Why pull from a remote source:
- Centralized flag management: Keep all flag definitions in a central repository or service
- Team collaboration: Share flag configurations across team members and environments
- Version control: Track changes to flag configurations over time
- Environment consistency: Ensure the same flag definitions are used across different environments
- Configuration as code: Treat flag definitions as versioned artifacts that can be reviewed and deployed

```
openfeature pull [flags]
```

### Options

```
      --auth-token string        The auth token for the flag source
      --flag-source-url string   The URL of the flag source
  -h, --help                     help for pull
      --no-prompt                Disable interactive prompts for missing default values
```

### Options inherited from parent commands

```
      --debug             Enable debug logging
  -m, --manifest string   Path to the flag manifest (default "flags.json")
      --no-input          Disable interactive prompts
```

### SEE ALSO

* [openfeature](openfeature.md)	 - CLI for OpenFeature.

