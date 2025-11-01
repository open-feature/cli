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

// PushFlags fetches remote flags, compares with local flags, and intelligently
// creates or updates flags as needed. Returns a PushResult with details of what was changed.
func (c *Client) PushFlags(ctx context.Context, localFlags *flagset.Flagset, remoteFlags *flagset.Flagset, schemaURL string) (*PushResult, error) {
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

	// Create new flags with retry logic
	for _, flag := range toCreate {
		flagKey := flag.Key // Capture for closure
		err := goretry.IfNeededWithContext(ctx, func(ctx context.Context) error {
			body, err := c.convertFlagToAPIBody(flag)
			if err != nil {
				return fmt.Errorf("failed to convert flag %s: %w", flagKey, err)
			}

			resp, err := c.apiClient.PostApiV1ManifestFlagsWithResponse(ctx, body)
			if err != nil {
				return fmt.Errorf("failed to create flag %s: %w", flagKey, err)
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

			resp, err := c.apiClient.PutApiV1ManifestFlagsKeyWithResponse(ctx, flagKey, body)
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
func (c *Client) convertFlagToAPIBody(flag flagset.Flag) (syncclient.PostApiV1ManifestFlagsJSONRequestBody, error) {
	// Convert flag type to API enum
	flagType := syncclient.PostApiV1ManifestFlagsJSONBodyType(flag.Type.String())

	// Marshal and unmarshal the defaultValue through JSON to properly set the union type
	defaultValueJSON, err := json.Marshal(flag.DefaultValue)
	if err != nil {
		return syncclient.PostApiV1ManifestFlagsJSONRequestBody{}, fmt.Errorf("failed to marshal defaultValue: %w", err)
	}

	var defaultValue syncclient.FlagDefaultValue
	if err := json.Unmarshal(defaultValueJSON, &defaultValue); err != nil {
		return syncclient.PostApiV1ManifestFlagsJSONRequestBody{}, fmt.Errorf("failed to unmarshal defaultValue: %w", err)
	}

	// Create the request body
	body := syncclient.PostApiV1ManifestFlagsJSONRequestBody{
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
func (c *Client) convertFlagToPutBody(flag flagset.Flag) (syncclient.PutApiV1ManifestFlagsKeyJSONRequestBody, error) {
	// Convert flag type to API enum
	flagType := syncclient.PutApiV1ManifestFlagsKeyJSONBodyType(flag.Type.String())

	// Marshal and unmarshal the defaultValue through JSON to properly set the union type
	defaultValueJSON, err := json.Marshal(flag.DefaultValue)
	if err != nil {
		return syncclient.PutApiV1ManifestFlagsKeyJSONRequestBody{}, fmt.Errorf("failed to marshal defaultValue: %w", err)
	}

	var defaultValue syncclient.FlagDefaultValue
	if err := json.Unmarshal(defaultValueJSON, &defaultValue); err != nil {
		return syncclient.PutApiV1ManifestFlagsKeyJSONRequestBody{}, fmt.Errorf("failed to unmarshal defaultValue: %w", err)
	}

	// Create the request body
	body := syncclient.PutApiV1ManifestFlagsKeyJSONRequestBody{
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
