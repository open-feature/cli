# OpenAPI Client Pattern for Remote Operations

This document describes the OpenAPI client generation pattern used in the OpenFeature CLI for implementing remote API operations. This pattern provides type-safe, well-documented client code for interacting with remote services.

## Overview

The OpenFeature CLI uses a code generation approach for remote API operations:

1. **Define**: Remote API functionality is defined in an OpenAPI specification (YAML)
2. **Generate**: A type-safe Go client is generated from the specification
3. **Wrap**: A thin wrapper provides convenience methods and integrates with CLI infrastructure
4. **Use**: Commands consume the wrapped client for all remote operations

This pattern ensures consistency, type safety, and maintainability across all remote operations in the CLI.

## Current Implementation

The current implementation consists of:

- **OpenAPI Specification**: `api/v0/sync.yaml` - Defines the Manifest Management API v0
- **Generated Client**: `internal/api/client/sync_client.gen.go` - Auto-generated from the spec
- **Wrapper Client**: `internal/api/sync/client.go` - Provides convenience methods and retry logic
- **Test Suite**: `internal/api/sync/retry_test.go` - Tests retry logic and error handling
- **CLI Commands**: `internal/cmd/pull.go` and `internal/cmd/push.go` - Use the wrapped client

### Dependencies

The implementation relies on these key libraries:

- **[oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)**: Generates type-safe Go clients from OpenAPI specs
- **[GoRetry](https://github.com/kriscoleman/GoRetry)**: Provides retry logic with exponential backoff for transient failures
- **[gock](https://github.com/h2non/gock)**: HTTP mocking for comprehensive testing without real servers

## Architecture

```
┌─────────────────────┐
│  OpenAPI Spec       │
│  (api/v0/sync.yaml) │
└──────────┬──────────┘
           │
           │ oapi-codegen
           ▼
┌─────────────────────┐
│  Generated Client   │
│  (internal/api/     │
│   client/*.gen.go)  │
└──────────┬──────────┘
           │
           │ wrapped by
           ▼
┌─────────────────────┐
│  Sync Client        │
│  (internal/api/     │
│   sync/client.go)   │
└──────────┬──────────┘
           │
           │ used by
           ▼
┌─────────────────────┐
│  CLI Commands       │
│  (internal/cmd/     │
│   pull.go,push.go)  │
└─────────────────────┘
```

## Benefits

### 1. Type Safety

- Compile-time checking of API requests and responses
- Eliminates string-based URL construction errors
- Strongly typed request/response models

### 2. Documentation as Code

- API contract is clearly defined in OpenAPI spec
- Self-documenting code through generated types
- Examples and descriptions in the spec translate to code comments

### 3. Maintainability

- Changes to API are made in one place (the spec)
- Regenerate client to update all consumers
- Clear separation between generated and custom code

### 4. Consistency

- All remote operations follow the same pattern
- Standardized error handling across endpoints
- Uniform logging and debugging capabilities

### 5. Extensibility

- Easy to add new endpoints (just update the spec)
- Supports multiple API versions side-by-side
- Foundation for plugin systems and provider implementations

## Implementation Guide

### Step 1: Define the OpenAPI Specification

Create or update the OpenAPI specification in `api/v0/` directory:

```yaml
# api/v0/sync.yaml
openapi: 3.0.3
info:
  title: Manifest Management API
  version: 1.0.0
  description: |
    CRUD endpoints for managing feature flags through the OpenFeature CLI.

paths:
  /openfeature/v0/manifest:
    get:
      summary: Get Project Manifest
      description: Returns the project manifest containing active flags
      security:
        - bearerAuth: []
      responses:
        '200':
          description: Manifest exported successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ManifestEnvelope'
```

**Best Practices:**

- Use semantic versioning for the API version
- Include detailed descriptions for all operations
- Define reusable schemas in `components/schemas`
- Specify all possible response codes
- Document authentication requirements

### Step 2: Generate the Client

Add a Makefile target to generate the client:

```makefile
# Generate OpenAPI clients
generate-sync-client:
 oapi-codegen -package syncclient \
  -generate types,client \
  -o internal/api/client/sync_client.gen.go \
  api/v0/sync.yaml
```

Run the generator:

```bash
make generate-sync-client
```

This produces type-safe Go code in `internal/api/client/sync_client.gen.go` including:

- Request/response types
- Client interface
- HTTP request builders
- Response parsers

**Important:** Never manually edit generated files. Always regenerate from the spec.

### Step 3: Create a Wrapper Client

Create a wrapper in `internal/api/sync/client.go` that:

- Provides convenience methods
- Handles common concerns (logging, retries, error mapping)
- Converts between API and internal models

```go
// internal/api/sync/client.go
package sync

import (
 "context"
 "encoding/json"
 "fmt"
 "net/http"

 syncclient "github.com/open-feature/cli/internal/api/client"
 "github.com/open-feature/cli/internal/flagset"
 "github.com/open-feature/cli/internal/logger"
)

// Client wraps the generated OpenAPI client with convenience methods
type Client struct {
 apiClient *syncclient.ClientWithResponses
 authToken string
}

// NewClient creates a new sync client
func NewClient(baseURL string, authToken string) (*Client, error) {
 // Configure HTTP client options
 var opts []syncclient.ClientOption

 if authToken != "" {
  opts = append(opts, syncclient.WithRequestEditorFn(
   func(ctx context.Context, req *http.Request) error {
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
    return nil
   }))
 }

 apiClient, err := syncclient.NewClientWithResponses(baseURL, opts...)
 if err != nil {
  return nil, fmt.Errorf("failed to create API client: %w", err)
 }

 return &Client{
  apiClient: apiClient,
  authToken: authToken,
 }, nil
}

// PullFlags fetches flags from the remote API (without retry logic)
func (c *Client) PullFlags(ctx context.Context) (*flagset.Flagset, error) {
 logger.Default.Debug("Fetching flags using sync API client")

 resp, err := c.apiClient.GetOpenfeatureV0ManifestWithResponse(ctx)
 if err != nil {
  return nil, fmt.Errorf("failed to fetch manifest: %w", err)
 }

 // Debug logging
 if resp.HTTPResponse != nil {
  logger.Default.Debug(fmt.Sprintf("Pull response: HTTP %d",
   resp.HTTPResponse.StatusCode))
 }

 // Handle error responses
 if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode >= 300 {
  if resp.JSON401 != nil {
   return nil, fmt.Errorf("authentication failed: %s",
    resp.JSON401.Error.Message)
  }
  // ... handle other error cases
 }

 // Convert API model to internal model
 return convertToFlagset(resp.JSON200)
}
```

**Key responsibilities of the wrapper:**

- Authentication/authorization setup
- Request/response logging for debugging
- Error handling and user-friendly error messages
- Model conversion (API types ↔ internal types)
- Retry logic for write operations (POST/PUT/DELETE) using GoRetry library
- Context propagation
- Smart push logic (comparing local vs remote flags before making changes)

The wrapper also defines result types for operations:

```go
// PushResult tracks changes made during a push operation
type PushResult struct {
 Created   []flagset.Flag  // Flags created on the remote
 Updated   []flagset.Flag  // Flags updated on the remote
 Unchanged []flagset.Flag  // Flags that didn't need changes
}
```

### Step 4: Use in Commands

Update commands to use the wrapped client:

```go
// internal/cmd/pull.go
func GetPullCmd() *cobra.Command {
 return &cobra.Command{
  Use:   "pull",
  Short: "Pull flag configurations from a remote source",
  RunE: func(cmd *cobra.Command, args []string) error {
   flagSourceUrl := config.GetFlagSourceUrl(cmd)
   authToken := config.GetAuthToken(cmd)

   // Use the sync API client
   flags, err := manifest.LoadFromSyncAPI(flagSourceUrl, authToken)
   if err != nil {
    return fmt.Errorf("error fetching flags: %w", err)
   }

   // Process flags...
   return nil
  },
 }
}
```

## URL Handling

The generated client expects the **base URL only** (domain and optionally port). The client automatically constructs full paths from the OpenAPI spec.

**Correct:**

```bash
# Base URL only
openfeature pull --flag-source-url https://api.example.com
openfeature pull --flag-source-url https://api.example.com:8080
```

**Incorrect:**

```bash
# Don't include API paths in the URL
openfeature pull --flag-source-url https://api.example.com/openfeature/v0/manifest
```

The generated client automatically:

- Adds a trailing slash to the base URL
- Constructs operation paths from the spec (e.g., `/openfeature/v0/manifest`)
- Handles path parameters and query strings

## Error Handling

The pattern provides multiple levels of error handling:

### 1. Network Errors

```go
resp, err := c.apiClient.GetOpenfeatureV0ManifestWithResponse(ctx)
if err != nil {
 return nil, fmt.Errorf("failed to fetch manifest: %w", err)
}
```

### 2. HTTP Status Errors

```go
if resp.HTTPResponse.StatusCode >= 400 {
 // Access typed error responses
 if resp.JSON401 != nil {
  return nil, fmt.Errorf("authentication failed: %s",
   resp.JSON401.Error.Message)
 }
}
```

### 3. Schema Validation

The generated client automatically validates responses against the OpenAPI schema.

## Debugging

Enable debug logging to see full HTTP request/response details:

```bash
openfeature pull --flag-source-url https://api.example.com --debug
```

The wrapper client should log:

- Request URLs and headers
- Request bodies (sanitized)
- Response status codes
- Response bodies (truncated if large)

Example debug output:

```
[DEBUG] Loading flags from sync API at https://api.example.com
[DEBUG] Pull response: HTTP 200 - OK
[DEBUG] Response body: {"flags":[{"key":"feature-x",...}]}
[DEBUG] Successfully pulled 5 flags
```

## Testing

### HTTP Mocking with Gock

The CLI uses [gock](https://github.com/h2non/gock) for mocking HTTP requests in tests. This provides a clean way to test retry logic, error handling, and different response scenarios without needing a real server.

**Why Gock?**
- **Simple API**: Fluent interface for setting up mock responses
- **Request Matching**: Can match on URL, headers, body, query params
- **Response Stacking**: Queue multiple responses for testing retry logic
- **Network Interception**: Works at the HTTP transport level, no code changes needed
- **Clean Test Isolation**: Each test can set up its own mock expectations

```go
// internal/api/sync/retry_test.go
import "github.com/h2non/gock"

func TestRetryLogic(t *testing.T) {
 defer gock.Off() // Clean up after test

 // First attempt: 500 error (will trigger retry)
 gock.New("https://api.example.com").
  Post("/openfeature/v0/manifest/flags").
  Reply(500).
  JSON(map[string]interface{}{
   "error": map[string]interface{}{
    "message": "Internal Server Error",
    "status":  500,
   },
  })

 // Second attempt: Success
 gock.New("https://api.example.com").
  Post("/openfeature/v0/manifest/flags").
  Reply(201).
  JSON(map[string]interface{}{
   "flag": map[string]interface{}{
    "key": "test-flag",
    "type": "boolean",
    "defaultValue": true,
   },
   "updatedAt": "2024-03-02T09:45:03.000Z",
  })

 // Create client and test retry behavior
 client, _ := NewClient("https://api.example.com", "test-token")

 // This will fail once, then succeed on retry
 result, err := client.PushFlags(ctx, localFlags, remoteFlags, false)
 assert.NoError(t, err)
 assert.Len(t, result.Created, 1)
}
```

**Key Testing Patterns with Gock:**

1. **Testing Retry Logic**: Stack multiple responses for the same endpoint to simulate failures followed by success
2. **Testing Error Handling**: Mock specific HTTP status codes and error responses
3. **Request Validation**: Use `MatchBody()` to verify request payloads
4. **Header Validation**: Use `MatchHeader()` to ensure authentication headers are sent

```go
// Testing authentication header
gock.New("https://api.example.com").
 Post("/openfeature/v0/manifest/flags").
 MatchHeader("Authorization", "Bearer test-token").
 Reply(201).
 JSON(successResponse)

// Testing request body
gock.New("https://api.example.com").
 Post("/openfeature/v0/manifest/flags").
 MatchType("json").
 BodyString(`{"key":"test-flag","type":"boolean","defaultValue":true}`).
 Reply(201).
 JSON(successResponse)
```

### When to Use Gock vs Interface Mocking

**Use Gock when:**
- Testing retry logic and transient failures
- Testing HTTP-specific behaviors (status codes, headers, timeouts)
- Testing the full request/response cycle
- Verifying exact request payloads and headers

**Use Interface Mocking when:**
- Testing business logic independent of HTTP
- Need fast unit tests without network overhead
- Testing complex data transformations
- Want to avoid HTTP transport complexity

### Unit Tests

For testing business logic without HTTP calls, you can also mock the generated client interface:

```go
// client_test.go
type mockSyncClient struct {
 getManifestFunc func(ctx context.Context) (*GetOpenfeatureV0ManifestResponse, error)
}

func (m *mockSyncClient) GetOpenfeatureV0ManifestWithResponse(ctx context.Context,
 reqEditors ...RequestEditorFn) (*GetOpenfeatureV0ManifestResponse, error) {
 return m.getManifestFunc(ctx)
}

func TestPullFlags(t *testing.T) {
 mock := &mockSyncClient{
  getManifestFunc: func(ctx context.Context) (*GetOpenfeatureV0ManifestResponse, error) {
   return &GetOpenfeatureV0ManifestResponse{
    JSON200: &ManifestEnvelope{
     Flags: []ManifestFlag{{Key: "test-flag"}},
    },
   }, nil
  },
 }

 // Test with mock client...
}
```

### Integration Tests

Test against a real implementation of the API:

```bash
# Run integration tests
make test-integration

# Run specific language integration tests
make test-integration-go
make test-integration-csharp
make test-integration-nodejs
```

## Future Extensions

### Plugin System

This pattern enables a plugin architecture for supporting multiple flag providers:

```go
// Provider interface using the pattern
type Provider interface {
 PullFlags(ctx context.Context) (*flagset.Flagset, error)
 PushFlags(ctx context.Context, flags *flagset.Flagset) error
}

// Each provider implements the interface
type FlagdStudioProvider struct {
 client *sync.Client
}

type LaunchDarklyProvider struct {
 client *launchdarkly.Client
}

// Register providers
providers := map[string]Provider{
 "flagd-studio": NewFlagdStudioProvider(config),
 "launchdarkly": NewLaunchDarklyProvider(config),
}
```

### Multiple API Versions

Support multiple API versions simultaneously:

```
api/
├── v0/
│   └── sync.yaml
└── v1/
    └── sync.yaml

internal/api/client/
├── v0/
│   └── sync_client.gen.go
└── v1/
    └── sync_client.gen.go
```

### Custom Operations

Add provider-specific operations while maintaining the core pattern:

```go
// Extended client for provider-specific features
type FlagdStudioClient struct {
 *sync.Client
}

func (c *FlagdStudioClient) ListEnvironments(ctx context.Context) ([]Environment, error) {
 // Provider-specific operation
}
```

## Common Patterns

### Retry Logic

The `PushFlags` method uses the [GoRetry library](https://github.com/kriscoleman/GoRetry) to handle transient failures when creating or updating flags:

```go
import goretry "github.com/kriscoleman/GoRetry"

func (c *Client) PushFlags(ctx context.Context, localFlags *flagset.Flagset, remoteFlags *flagset.Flagset, dryRun bool) (*PushResult, error) {
 // ... flag comparison logic ...

 for _, flag := range toCreate {
  err := goretry.IfNeededWithContext(ctx, func(ctx context.Context) error {
   body, err := c.convertFlagToAPIBody(flag)
   if err != nil {
    return fmt.Errorf("failed to convert flag %s: %w", flag.Key, err)
   }

   resp, err := c.apiClient.PostOpenfeatureV0ManifestFlagsWithResponse(ctx, body)
   if err != nil {
    return fmt.Errorf("failed to create flag %s: %w", flag.Key, err)
   }

   if resp.HTTPResponse.StatusCode >= 500 {
    return &httpError{statusCode: resp.HTTPResponse.StatusCode} // Retryable
   }

   return nil
  }, goretry.WithTransientErrorFunc(isTransientHTTPError))

  // ... handle result ...
 }
}
```

Note: Retry logic is only used for write operations (POST, PUT, DELETE) where transient failures are more impactful. Read operations like `PullFlags` don't implement retry logic by default.

### Pagination

```go
func (c *Client) ListAllFlags(ctx context.Context) ([]*flagset.Flag, error) {
 var allFlags []*flagset.Flag
 cursor := ""

 for {
  // Note: Pagination is not yet implemented in the current API spec
  // This is an example of how it could work
  resp, err := c.apiClient.GetOpenfeatureV0ManifestWithResponse(ctx,
   &GetOpenfeatureV0ManifestParams{Cursor: &cursor})
  if err != nil {
   return nil, err
  }

  allFlags = append(allFlags, convertFlags(resp.JSON200.Flags)...)

  // Check for pagination cursor in response
  if resp.JSON200.NextCursor == nil {
   break
  }
  cursor = *resp.JSON200.NextCursor
 }

 return allFlags, nil
}
```

### Streaming

```go
func (c *Client) StreamFlags(ctx context.Context) (<-chan *flagset.Flag, error) {
 ch := make(chan *flagset.Flag)

 go func() {
  defer close(ch)

  // Server-sent events or WebSocket connection
  stream, err := c.apiClient.StreamManifest(ctx)
  if err != nil {
   return
  }

  for event := range stream {
   ch <- convertFlag(event)
  }
 }()

 return ch, nil
}
```

## Checklist for Adding a New Remote Operation

- [ ] Update OpenAPI specification in `api/v0/`
- [ ] Regenerate client: `make generate-api`
- [ ] Add wrapper method in `internal/api/sync/client.go`
- [ ] Add debug logging for requests/responses
- [ ] Handle all error cases (4xx, 5xx)
- [ ] Convert API models to internal models
- [ ] Add unit tests with mocked client
- [ ] Add integration tests
- [ ] Update command to use new operation
- [ ] Update documentation
- [ ] Test with `--debug` flag

## References

- [OpenAPI Specification](https://spec.openapis.org/oas/v3.0.3)
- [oapi-codegen Documentation](https://github.com/oapi-codegen/oapi-codegen)
- [Sync API Specification](../../api/v0/sync.yaml)

## Related Patterns

- **Code Generation Pattern**: Used for generating flag accessors (see [DESIGN.md](../DESIGN.md))
- **Provider Pattern**: Abstraction for different flag provider implementations
- **Repository Pattern**: Separates data access logic from business logic

---

For questions or suggestions about this pattern, please open an issue or discussion in the GitHub repository.
