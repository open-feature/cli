<!-- markdownlint-disable-file -->
<!-- WARNING: THIS DOC IS AUTO-GENERATED. DO NOT EDIT! -->
## openfeature push

Push flag configurations to a remote source

### Synopsis

The push command syncs local flag configurations to a remote flag management service.

This command reads your local flag manifest and intelligently pushes it to a specified
remote destination. It performs a smart push by:

1. Fetching existing flags from the remote
2. Comparing local flags with remote flags
3. Creating new flags that don't exist remotely
4. Updating existing flags that have changed

This approach ensures idempotent operations and prevents conflicts.

The pushed data follows the Manifest Management API OpenAPI specification defined at:
api/v0/sync.yaml

The API uses individual flag endpoints:
- POST /openfeature/v0/manifest/flags - Creates new flags
- PUT /openfeature/v0/manifest/flags/{key} - Updates existing flags
- GET /openfeature/v0/manifest - Fetches existing flags for comparison

Remote services implementing this API should accept the flag data in the format
specified by the OpenFeature flag manifest schema.

Note: The file:// scheme is not supported for push operations.
For local file operations, use standard shell commands like cp or mv.

```
openfeature push [flags]
```

### Examples

```
  # Push flags to a remote HTTPS endpoint (smart push: creates and updates as needed)
  openfeature push --flag-source-url https://api.example.com --auth-token secret-token

  # Push flags to an HTTP endpoint (development)
  openfeature push --flag-source-url http://localhost:8080

  # Dry run to preview what would be sent
  openfeature push --flag-source-url https://api.example.com --dry-run
```

### Options

```
      --auth-token string        The auth token for the flag destination
      --debug                    Enable debug logging
      --dry-run                  Preview changes without pushing
      --flag-source-url string   The URL of the flag destination
  -h, --help                     help for push
  -m, --manifest string          Path to the flag manifest (default "flags.json")
      --no-input                 Disable interactive prompts
```

### SEE ALSO

* [openfeature](openfeature.md)	 - CLI for OpenFeature.

