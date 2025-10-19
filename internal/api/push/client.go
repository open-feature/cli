// Package push provides a wrapper around the generated OpenAPI client for push operations
package push

import (
	"context"
	"fmt"
	"net/http"
	"time"

	pushclient "github.com/open-feature/cli/internal/api/client"
	"github.com/open-feature/cli/internal/flagset"
)

// Client wraps the generated OpenAPI client with convenience methods
type Client struct {
	apiClient *pushclient.ClientWithResponses
	authToken string
}

// NewClient creates a new push client
func NewClient(baseURL string, authToken string) (*Client, error) {
	// Create a custom HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Add authentication if provided
	var opts []pushclient.ClientOption
	opts = append(opts, pushclient.WithHTTPClient(httpClient))

	if authToken != "" {
		opts = append(opts, pushclient.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
			return nil
		}))
	}

	// Add standard headers
	opts = append(opts, pushclient.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "openfeature-cli/push")
		return nil
	}))

	apiClient, err := pushclient.NewClientWithResponses(baseURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	return &Client{
		apiClient: apiClient,
		authToken: authToken,
	}, nil
}

// PushFlags pushes flags to the remote server using POST
func (c *Client) PushFlags(ctx context.Context, flags *flagset.Flagset, schemaURL string) error {
	payload := c.convertToAPIPayload(flags, schemaURL)

	resp, err := c.apiClient.PushFlagsWithResponse(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to push flags: %w", err)
	}

	return c.handleResponse(resp.HTTPResponse, resp.Body)
}

// ReplaceFlags replaces all flags on the remote server using PUT
func (c *Client) ReplaceFlags(ctx context.Context, flags *flagset.Flagset, schemaURL string) error {
	payload := c.convertToAPIPayload(flags, schemaURL)

	resp, err := c.apiClient.ReplaceFlagsWithResponse(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to replace flags: %w", err)
	}

	return c.handleResponse(resp.HTTPResponse, resp.Body)
}

// convertToAPIPayload converts internal flagset to API payload format
func (c *Client) convertToAPIPayload(flags *flagset.Flagset, schemaURL string) pushclient.FlagPayload {
	apiFlags := make([]pushclient.Flag, 0, len(flags.Flags))

	for _, flag := range flags.Flags {
		// Convert flag type to API enum
		flagType := pushclient.FlagType(flag.Type.String())

		// Create the API flag with dynamic defaultValue support
		apiFlag := pushclient.Flag{
			Key:          flag.Key,
			Type:         flagType,
			DefaultValue: flag.DefaultValue, // Now supports any type via interface{}
		}

		// Add description if present
		if flag.Description != "" {
			apiFlag.Description = &flag.Description
		}

		apiFlags = append(apiFlags, apiFlag)
	}

	payload := pushclient.FlagPayload{
		Flags: apiFlags,
	}

	// Add schema reference if provided
	if schemaURL != "" {
		payload.Schema = &schemaURL
	}

	return payload
}

// handleResponse processes the HTTP response and returns an appropriate error
func (c *Client) handleResponse(resp *http.Response, body []byte) error {
	if resp == nil {
		return fmt.Errorf("received nil response")
	}

	// Check for successful status code
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	// Return error with status and body for debugging
	return fmt.Errorf("received error response from destination (status %d): %s", resp.StatusCode, string(body))
}
