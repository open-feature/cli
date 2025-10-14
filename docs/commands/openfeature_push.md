<!-- markdownlint-disable-file -->
<!-- WARNING: THIS DOC IS AUTO-GENERATED. DO NOT EDIT! -->
## openfeature push

Push flag configurations to a remote source

### Synopsis

The push command uploads local flag configurations to a remote flag management service.

This command reads your local flag manifest and pushes it to a specified remote destination.
It supports HTTP and HTTPS protocols for remote endpoints.

Note: The file:// scheme is not supported for push operations.
For local file operations, use standard shell commands like cp or mv.

```
openfeature push [flags]
```

### Examples

```
  # Push flags to a remote HTTPS endpoint
  openfeature push --flag-destination-url https://api.example.com/flags --auth-token secret-token

  # Push flags to an HTTP endpoint (development)
  openfeature push --flag-destination-url http://localhost:8080/flags

  # Push using PUT method instead of POST
  openfeature push --flag-destination-url https://api.example.com/flags/my-app --method PUT

  # Dry run to preview what would be sent
  openfeature push --flag-destination-url https://api.example.com/flags --dry-run
```

### Options

```
      --auth-token string             The auth token for the flag destination
      --debug                         Enable debug logging
      --dry-run                       Preview changes without pushing
      --flag-destination-url string   The URL of the flag destination
  -h, --help                          help for push
  -m, --manifest string               Path to the flag manifest (default "flags.json")
      --method string                 HTTP method to use (POST or PUT) (default "POST")
      --no-input                      Disable interactive prompts
```

### SEE ALSO

* [openfeature](openfeature.md)	 - CLI for OpenFeature.

