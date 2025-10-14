// Package sync provides a wrapper around the generated OpenAPI client for sync operations
package sync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	goretry "github.com/kriscoleman/GoRetry"
	syncclient "github.com/open-feature/cli/internal/api/client"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/logger"
)

// Client wraps the generated OpenAPI client with convenience methods
type Client struct {
	apiClient *syncclient.ClientWithResponses
	authToken string
}

// httpError wraps an HTTP response status code for retry logic
type httpError struct {
	statusCode int
	message    string
}

func (e *httpError) Error() string {
	return e.message
}

// isTransientHTTPError determines if an error should trigger a retry.
// Returns true for:
// - 5xx server errors (transient)
// - Network errors (timeouts, temporary failures)
// Returns false for:
// - 4xx client errors (permanent)
// - Successful responses (2xx, 3xx)
func isTransientHTTPError(err error) bool {
	if err == nil {
		return false
	}

	// Check if it's an HTTP error with a status code
	var httpErr *httpError
	if errors.As(err, &httpErr) {
		// Retry on 5xx server errors
		if httpErr.statusCode >= 500 && httpErr.statusCode < 600 {
			return true
		}
		// Don't retry on 4xx client errors or successful responses
		return false
	}

	// For non-HTTP errors, use default transient error detection
	// This catches network errors, timeouts, etc.
	return goretry.DefaultTransientErrorFunc(err)
}

// NewClient creates a new sync client
func NewClient(baseURL string, authToken string) (*Client, error) {
	// Create a custom HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Add authentication if provided
	var opts []syncclient.ClientOption
	opts = append(opts, syncclient.WithHTTPClient(httpClient))

	if authToken != "" {
		opts = append(opts, syncclient.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
			return nil
		}))
	}

	// Add standard headers
	opts = append(opts, syncclient.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "openfeature-cli/sync")
		return nil
	}))

	apiClient, err := syncclient.NewClientWithResponses(baseURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	return &Client{
		apiClient: apiClient,
		authToken: authToken,
	}, nil
}

// PushResult contains the results of a push operation
type PushResult struct {
	Created   []flagset.Flag
	Updated   []flagset.Flag
	Unchanged []flagset.Flag
}

// PullFlags fetches flags from the remote API
func (c *Client) PullFlags(ctx context.Context) (*flagset.Flagset, error) {
	logger.Default.Debug("Fetching flags using sync API client")

	resp, err := c.apiClient.GetOpenfeatureV0ManifestWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manifest: %w", err)
	}

	// Debug: log HTTP response details
	if resp.HTTPResponse != nil {
		logger.Default.Debug(fmt.Sprintf("Pull response: HTTP %d - %s", resp.HTTPResponse.StatusCode, resp.HTTPResponse.Status))
		if len(resp.Body) > 0 {
			logger.Default.Debug(fmt.Sprintf("Response body: %s", string(resp.Body)))
		}
	}

	// Check for successful status code
	if resp.HTTPResponse == nil {
		return nil, fmt.Errorf("received nil HTTP response")
	}

	if resp.HTTPResponse.StatusCode < 200 || resp.HTTPResponse.StatusCode >= 300 {
		// Try to parse error response
		if resp.JSON401 != nil {
			return nil, fmt.Errorf("authentication failed: %s", resp.JSON401.Error.Message)
		} else if resp.JSON403 != nil {
			return nil, fmt.Errorf("authorization failed: %s", resp.JSON403.Error.Message)
		} else if resp.JSON500 != nil {
			return nil, fmt.Errorf("server error: %s", resp.JSON500.Error.Message)
		}
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.HTTPResponse.StatusCode, string(resp.Body))
	}

	// Parse successful response
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("expected manifest data but got none")
	}

	// Convert from API model to internal flagset model
	flags := make([]flagset.Flag, 0, len(resp.JSON200.Flags))
	for _, apiFlag := range resp.JSON200.Flags {
		// Parse flag type from string
		flagType, err := flagset.ParseFlagType(string(apiFlag.Type))
		if err != nil {
			return nil, fmt.Errorf("failed to parse flag type for %s: %w", apiFlag.Key, err)
		}

		flag := flagset.Flag{
			Key:  apiFlag.Key,
			Type: flagType,
		}

		// Set optional fields
		if apiFlag.Description != nil {
			flag.Description = *apiFlag.Description
		}

		// Convert defaultValue from union type by marshaling and unmarshaling
		defaultValueJSON, err := json.Marshal(apiFlag.DefaultValue)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal defaultValue for flag %s: %w", flag.Key, err)
		}
		if err := json.Unmarshal(defaultValueJSON, &flag.DefaultValue); err != nil {
			return nil, fmt.Errorf("failed to parse defaultValue for flag %s: %w", flag.Key, err)
		}

		flags = append(flags, flag)
	}

	logger.Default.Debug(fmt.Sprintf("Successfully pulled %d flags", len(flags)))

	return &flagset.Flagset{Flags: flags}, nil
}

// PushFlags fetches remote flags, compares with local flags, and intelligently
// creates or updates flags as needed. Returns a PushResult with details of what was changed.
// If dryRun is true, only performs the comparison without making actual API calls.
func (c *Client) PushFlags(ctx context.Context, localFlags *flagset.Flagset, remoteFlags *flagset.Flagset, schemaURL string, dryRun bool) (*PushResult, error) {
	// Build a map of remote flags for quick lookup
	remoteFlagMap := make(map[string]flagset.Flag)
	for _, flag := range remoteFlags.Flags {
		remoteFlagMap[flag.Key] = flag
	}

	var toCreate []flagset.Flag
	var toUpdate []flagset.Flag

	// Determine which flags need to be created vs updated
	for _, localFlag := range localFlags.Flags {
		if remoteFlag, exists := remoteFlagMap[localFlag.Key]; exists {
			// Only update if the flag has actually changed
			if !flagsEqual(localFlag, remoteFlag) {
				toUpdate = append(toUpdate, localFlag)
			}
		} else {
			toCreate = append(toCreate, localFlag)
		}
	}

	result := &PushResult{}

	// If dry run, skip actual API calls and just return what would be done
	if dryRun {
		result.Created = toCreate
		result.Updated = toUpdate
		return result, nil
	}

	// Create new flags with retry logic
	for _, flag := range toCreate {
		flagKey := flag.Key // Capture for closure
		err := goretry.IfNeededWithContext(ctx, func(ctx context.Context) error {
			body, err := c.convertFlagToAPIBody(flag)
			if err != nil {
				return fmt.Errorf("failed to convert flag %s: %w", flagKey, err)
			}

			// Debug: log what we're sending
			if logger.Default.IsDebugEnabled() {
				bodyJSON, _ := json.MarshalIndent(body, "", "  ")
				logger.Default.Debug(fmt.Sprintf("Sending POST for %s:\n%s", flagKey, string(bodyJSON)))
			}

			resp, err := c.apiClient.PostOpenfeatureV0ManifestFlagsWithResponse(ctx, body)
			if err != nil {
				return fmt.Errorf("failed to create flag %s: %w", flagKey, err)
			}

			// Debug: log server response
			if logger.Default.IsDebugEnabled() {
				logger.Default.Debug(fmt.Sprintf("Server response for %s:\n%s", flagKey, string(resp.Body)))
			}

			return c.handleFlagResponse(resp.HTTPResponse, resp.Body, flagKey, "create")
		}, goretry.WithTransientErrorFunc(isTransientHTTPError))

		if err != nil {
			return nil, err
		}
		result.Created = append(result.Created, flag)
	}

	// Update existing flags with retry logic
	for _, flag := range toUpdate {
		flagKey := flag.Key // Capture for closure
		err := goretry.IfNeededWithContext(ctx, func(ctx context.Context) error {
			body, err := c.convertFlagToPutBody(flag)
			if err != nil {
				return fmt.Errorf("failed to convert flag %s: %w", flagKey, err)
			}

			// Debug: log what we're sending
			if logger.Default.IsDebugEnabled() {
				bodyJSON, _ := json.MarshalIndent(body, "", "  ")
				logger.Default.Debug(fmt.Sprintf("Sending PUT for %s:\n%s", flagKey, string(bodyJSON)))
			}

			resp, err := c.apiClient.PutOpenfeatureV0ManifestFlagsKeyWithResponse(ctx, flagKey, body)
			if err != nil {
				return fmt.Errorf("failed to update flag %s: %w", flagKey, err)
			}

			// Debug: log server response
			if logger.Default.IsDebugEnabled() {
				logger.Default.Debug(fmt.Sprintf("Server response for %s:\n%s", flagKey, string(resp.Body)))
			}

			return c.handleFlagResponse(resp.HTTPResponse, resp.Body, flagKey, "update")
		}, goretry.WithTransientErrorFunc(isTransientHTTPError))

		if err != nil {
			return nil, err
		}
		result.Updated = append(result.Updated, flag)
	}

	return result, nil
}

// convertFlagToAPIBody converts internal flag to POST API body format
func (c *Client) convertFlagToAPIBody(flag flagset.Flag) (syncclient.PostOpenfeatureV0ManifestFlagsJSONRequestBody, error) {
	// Convert flag type to API enum
	flagType := syncclient.PostOpenfeatureV0ManifestFlagsJSONBodyType(flag.Type.String())

	// Marshal and unmarshal the defaultValue through JSON to properly set the union type
	defaultValueJSON, err := json.Marshal(flag.DefaultValue)
	if err != nil {
		return syncclient.PostOpenfeatureV0ManifestFlagsJSONRequestBody{}, fmt.Errorf("failed to marshal defaultValue: %w", err)
	}

	var defaultValue syncclient.FlagDefaultValue
	if err := json.Unmarshal(defaultValueJSON, &defaultValue); err != nil {
		return syncclient.PostOpenfeatureV0ManifestFlagsJSONRequestBody{}, fmt.Errorf("failed to unmarshal defaultValue: %w", err)
	}

	// Create the request body
	body := syncclient.PostOpenfeatureV0ManifestFlagsJSONRequestBody{
		Key:          flag.Key,
		Type:         flagType,
		DefaultValue: defaultValue,
	}

	// Add description if present
	if flag.Description != "" {
		body.Description = &flag.Description
	}

	return body, nil
}

// convertFlagToPutBody converts internal flag to PUT API body format
func (c *Client) convertFlagToPutBody(flag flagset.Flag) (syncclient.PutOpenfeatureV0ManifestFlagsKeyJSONRequestBody, error) {
	// Convert flag type to API enum
	flagType := syncclient.PutOpenfeatureV0ManifestFlagsKeyJSONBodyType(flag.Type.String())

	// Marshal and unmarshal the defaultValue through JSON to properly set the union type
	defaultValueJSON, err := json.Marshal(flag.DefaultValue)
	if err != nil {
		return syncclient.PutOpenfeatureV0ManifestFlagsKeyJSONRequestBody{}, fmt.Errorf("failed to marshal defaultValue: %w", err)
	}

	var defaultValue syncclient.FlagDefaultValue
	if err := json.Unmarshal(defaultValueJSON, &defaultValue); err != nil {
		return syncclient.PutOpenfeatureV0ManifestFlagsKeyJSONRequestBody{}, fmt.Errorf("failed to unmarshal defaultValue: %w", err)
	}

	// Create the request body
	body := syncclient.PutOpenfeatureV0ManifestFlagsKeyJSONRequestBody{
		Key:          flag.Key,
		Type:         flagType,
		DefaultValue: &defaultValue,
	}

	// Add description if present
	if flag.Description != "" {
		body.Description = &flag.Description
	}

	return body, nil
}

// handleFlagResponse processes the HTTP response for individual flag operations
func (c *Client) handleFlagResponse(resp *http.Response, body []byte, flagKey string, operation string) error {
	if resp == nil {
		return fmt.Errorf("received nil response for flag %s", flagKey)
	}

	// Debug: log HTTP response details
	logger.Default.Debug(fmt.Sprintf("%s flag %s: HTTP %d - %s", operation, flagKey, resp.StatusCode, resp.Status))
	if len(body) > 0 {
		logger.Default.Debug(fmt.Sprintf("Response body: %s", string(body)))
	}

	// Check for successful status code
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	// Build error message
	var message string
	// Try to parse error response for better error messages
	var errorResp syncclient.ErrorResponse
	if err := json.Unmarshal(body, &errorResp); err == nil {
		message = fmt.Sprintf("failed to %s flag %s (status %d): %s", operation, flagKey, resp.StatusCode, errorResp.Error.Message)
	} else {
		// Fallback to raw response
		message = fmt.Sprintf("failed to %s flag %s (status %d): %s", operation, flagKey, resp.StatusCode, string(body))
	}

	// Return httpError so retry logic can determine if it's transient
	return &httpError{
		statusCode: resp.StatusCode,
		message:    message,
	}
}

// flagsEqual compares two flags to determine if they are effectively identical
func flagsEqual(a, b flagset.Flag) bool {
	// Compare key, type, and defaultValue
	if a.Key != b.Key || a.Type != b.Type {
		return false
	}

	// Marshal both defaultValues to JSON for comparison
	aJSON, err := json.Marshal(a.DefaultValue)
	if err != nil {
		return false
	}

	bJSON, err := json.Marshal(b.DefaultValue)
	if err != nil {
		return false
	}

	// Compare JSON representations
	if string(aJSON) != string(bJSON) {
		logger.Default.Debug(fmt.Sprintf("Flag %s differs:\n  Local: %s\n  Remote: %s", a.Key, string(aJSON), string(bJSON)))
		return false
	}

	// Compare descriptions (both empty or identical)
	if a.Description != b.Description {
		logger.Default.Debug(fmt.Sprintf("Flag %s description differs:\n  Local: %q\n  Remote: %q", a.Key, a.Description, b.Description))
		return false
	}

	return true
}
