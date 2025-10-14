package sync

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetryLogic(t *testing.T) {
	t.Run("retries on 5xx errors and eventually succeeds", func(t *testing.T) {
		defer gock.Off()

		// First attempt: 500 error
		gock.New("https://api.example.com").
			Post("/openfeature/v0/manifest/flags").
			Reply(500).
			JSON(map[string]interface{}{
				"error": map[string]interface{}{
					"message": "Internal Server Error",
					"status":  500,
				},
			})

		// Second attempt: 500 error
		gock.New("https://api.example.com").
			Post("/openfeature/v0/manifest/flags").
			Reply(500).
			JSON(map[string]interface{}{
				"error": map[string]interface{}{
					"message": "Internal Server Error",
					"status":  500,
				},
			})

		// Third attempt: success
		gock.New("https://api.example.com").
			Post("/openfeature/v0/manifest/flags").
			Reply(201).
			JSON(map[string]interface{}{
				"flag": map[string]interface{}{
					"key": "test-flag",
				},
				"updatedAt": "2024-03-02T09:45:03.000Z",
			})

		client, err := NewClient("https://api.example.com", "")
		require.NoError(t, err)

		ctx := context.Background()
		localFlags := &flagset.Flagset{
			Flags: []flagset.Flag{
				{Key: "test-flag", Type: flagset.BoolType, DefaultValue: true},
			},
		}
		remoteFlags := &flagset.Flagset{Flags: []flagset.Flag{}}

		_, err = client.PushFlags(ctx, localFlags, remoteFlags, "", false)
		assert.NoError(t, err, "Should succeed after retries")
		assert.True(t, gock.IsDone(), "All expected requests should be made")
	})

	t.Run("does not retry on 4xx errors", func(t *testing.T) {
		defer gock.Off()

		// Track the number of attempts
		attemptCount := 0

		// Mock the POST request to fail with 400 (bad request)
		gock.New("https://api.example.com").
			Post("/openfeature/v0/manifest/flags").
			AddMatcher(func(req *http.Request, _ *gock.Request) (bool, error) {
				attemptCount++
				return true, nil
			}).
			Reply(400).
			JSON(map[string]interface{}{
				"error": map[string]interface{}{
					"message": "Bad Request",
					"status":  400,
				},
			})

		client, err := NewClient("https://api.example.com", "")
		require.NoError(t, err)

		ctx := context.Background()
		localFlags := &flagset.Flagset{
			Flags: []flagset.Flag{
				{Key: "test-flag", Type: flagset.BoolType, DefaultValue: true},
			},
		}
		remoteFlags := &flagset.Flagset{Flags: []flagset.Flag{}}

		_, err = client.PushFlags(ctx, localFlags, remoteFlags, "", false)
		assert.Error(t, err, "Should fail with 400 error")
		assert.Contains(t, err.Error(), "400")
		assert.Equal(t, 1, attemptCount, "Should only attempt once (no retries for 4xx)")
	})

	t.Run("exhausts retries on persistent 5xx errors", func(t *testing.T) {
		defer gock.Off()

		// Track the number of attempts
		attemptCount := 0

		// Mock the POST request to always fail with 503
		gock.New("https://api.example.com").
			Post("/openfeature/v0/manifest/flags").
			Times(3). // Default max attempts
			AddMatcher(func(req *http.Request, _ *gock.Request) (bool, error) {
				attemptCount++
				return true, nil
			}).
			Reply(503).
			JSON(map[string]interface{}{
				"error": map[string]interface{}{
					"message": "Service Unavailable",
					"status":  503,
				},
			})

		client, err := NewClient("https://api.example.com", "")
		require.NoError(t, err)

		ctx := context.Background()
		localFlags := &flagset.Flagset{
			Flags: []flagset.Flag{
				{Key: "test-flag", Type: flagset.BoolType, DefaultValue: true},
			},
		}
		remoteFlags := &flagset.Flagset{Flags: []flagset.Flag{}}

		_, err = client.PushFlags(ctx, localFlags, remoteFlags, "", false)
		assert.Error(t, err, "Should fail after exhausting retries")
		assert.Contains(t, err.Error(), "503")
		assert.Equal(t, 3, attemptCount, "Should attempt 3 times (max attempts)")
	})

	t.Run("retries PUT operations on 5xx errors", func(t *testing.T) {
		defer gock.Off()

		// First attempt: 502 error
		gock.New("https://api.example.com").
			Put("/openfeature/v0/manifest/flags/existing-flag").
			Reply(502).
			JSON(map[string]interface{}{
				"error": map[string]interface{}{
					"message": "Bad Gateway",
					"status":  502,
				},
			})

		// Second attempt: success
		gock.New("https://api.example.com").
			Put("/openfeature/v0/manifest/flags/existing-flag").
			Reply(200).
			JSON(map[string]interface{}{
				"flag": map[string]interface{}{
					"key": "existing-flag",
				},
				"updatedAt": "2024-03-02T09:45:03.000Z",
			})

		client, err := NewClient("https://api.example.com", "")
		require.NoError(t, err)

		ctx := context.Background()
		localFlags := &flagset.Flagset{
			Flags: []flagset.Flag{
				{Key: "existing-flag", Type: flagset.BoolType, DefaultValue: true},
			},
		}
		remoteFlags := &flagset.Flagset{
			Flags: []flagset.Flag{
				{Key: "existing-flag", Type: flagset.BoolType, DefaultValue: false},
			},
		}

		_, err = client.PushFlags(ctx, localFlags, remoteFlags, "", false)
		assert.NoError(t, err, "Should succeed after retry")
		assert.True(t, gock.IsDone(), "All expected requests should be made")
	})

	t.Run("does not retry PUT operations on 404 errors", func(t *testing.T) {
		defer gock.Off()

		// Track the number of attempts
		attemptCount := 0

		// Mock the PUT request to fail with 404 (not found)
		gock.New("https://api.example.com").
			Put("/openfeature/v0/manifest/flags/nonexistent-flag").
			AddMatcher(func(req *http.Request, _ *gock.Request) (bool, error) {
				attemptCount++
				return true, nil
			}).
			Reply(404).
			JSON(map[string]interface{}{
				"error": map[string]interface{}{
					"message": "Flag not found",
					"status":  404,
				},
			})

		client, err := NewClient("https://api.example.com", "")
		require.NoError(t, err)

		ctx := context.Background()
		localFlags := &flagset.Flagset{
			Flags: []flagset.Flag{
				{Key: "nonexistent-flag", Type: flagset.BoolType, DefaultValue: true},
			},
		}
		remoteFlags := &flagset.Flagset{
			Flags: []flagset.Flag{
				{Key: "nonexistent-flag", Type: flagset.BoolType, DefaultValue: false},
			},
		}

		_, err = client.PushFlags(ctx, localFlags, remoteFlags, "", false)
		assert.Error(t, err, "Should fail with 404 error")
		assert.Contains(t, err.Error(), "404")
		assert.Equal(t, 1, attemptCount, "Should only attempt once (no retries for 4xx)")
	})

	t.Run("retries multiple flags independently", func(t *testing.T) {
		defer gock.Off()

		// Mock POST for first flag - fails once, then succeeds
		gock.New("https://api.example.com").
			Post("/openfeature/v0/manifest/flags").
			AddMatcher(func(req *http.Request, _ *gock.Request) (bool, error) {
				var body map[string]interface{}
				_ = json.NewDecoder(req.Body).Decode(&body)
				return body["key"] == "flag1", nil
			}).
			Reply(500).
			JSON(map[string]interface{}{
				"error": map[string]interface{}{
					"message": "Internal Server Error",
					"status":  500,
				},
			})

		gock.New("https://api.example.com").
			Post("/openfeature/v0/manifest/flags").
			AddMatcher(func(req *http.Request, _ *gock.Request) (bool, error) {
				var body map[string]interface{}
				_ = json.NewDecoder(req.Body).Decode(&body)
				return body["key"] == "flag1", nil
			}).
			Reply(201).
			JSON(map[string]interface{}{
				"flag": map[string]interface{}{
					"key": "flag1",
				},
				"updatedAt": "2024-03-02T09:45:03.000Z",
			})

		// Mock POST for second flag - succeeds on first try
		gock.New("https://api.example.com").
			Post("/openfeature/v0/manifest/flags").
			AddMatcher(func(req *http.Request, _ *gock.Request) (bool, error) {
				var body map[string]interface{}
				_ = json.NewDecoder(req.Body).Decode(&body)
				return body["key"] == "flag2", nil
			}).
			Reply(201).
			JSON(map[string]interface{}{
				"flag": map[string]interface{}{
					"key": "flag2",
				},
				"updatedAt": "2024-03-02T09:45:03.000Z",
			})

		client, err := NewClient("https://api.example.com", "")
		require.NoError(t, err)

		ctx := context.Background()
		localFlags := &flagset.Flagset{
			Flags: []flagset.Flag{
				{Key: "flag1", Type: flagset.BoolType, DefaultValue: true},
				{Key: "flag2", Type: flagset.StringType, DefaultValue: "test"},
			},
		}
		remoteFlags := &flagset.Flagset{Flags: []flagset.Flag{}}

		_, err = client.PushFlags(ctx, localFlags, remoteFlags, "", false)
		assert.NoError(t, err, "Should succeed with both flags")
		assert.True(t, gock.IsDone(), "All expected requests should be made")
	})

	t.Run("dry run mode does not make API calls", func(t *testing.T) {
		// No gock mocks needed - dry run should not make any HTTP requests

		client, err := NewClient("https://api.example.com", "")
		require.NoError(t, err)

		ctx := context.Background()
		localFlags := &flagset.Flagset{
			Flags: []flagset.Flag{
				{Key: "new-flag", Type: flagset.BoolType, DefaultValue: true, Description: "New flag"},
				{Key: "existing-flag", Type: flagset.StringType, DefaultValue: "updated", Description: "Updated flag"},
				{Key: "unchanged-flag", Type: flagset.IntType, DefaultValue: 42, Description: "Unchanged"},
			},
		}

		remoteFlags := &flagset.Flagset{
			Flags: []flagset.Flag{
				{Key: "existing-flag", Type: flagset.StringType, DefaultValue: "old", Description: "Old flag"},
				{Key: "unchanged-flag", Type: flagset.IntType, DefaultValue: 42, Description: "Unchanged"},
			},
		}

		// Run in dry run mode
		result, err := client.PushFlags(ctx, localFlags, remoteFlags, "", true)
		assert.NoError(t, err, "Dry run should not error")
		require.NotNil(t, result)

		// Verify the result shows what would be created and updated
		assert.Len(t, result.Created, 1, "Should identify 1 flag to create")
		assert.Equal(t, "new-flag", result.Created[0].Key)

		assert.Len(t, result.Updated, 1, "Should identify 1 flag to update")
		assert.Equal(t, "existing-flag", result.Updated[0].Key)
	})
}

func TestIsTransientHTTPError(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		shouldRetry bool
	}{
		{
			name:        "nil error is not transient",
			err:         nil,
			shouldRetry: false,
		},
		{
			name: "500 error is transient",
			err: &httpError{
				statusCode: 500,
				message:    "Internal Server Error",
			},
			shouldRetry: true,
		},
		{
			name: "502 error is transient",
			err: &httpError{
				statusCode: 502,
				message:    "Bad Gateway",
			},
			shouldRetry: true,
		},
		{
			name: "503 error is transient",
			err: &httpError{
				statusCode: 503,
				message:    "Service Unavailable",
			},
			shouldRetry: true,
		},
		{
			name: "504 error is transient",
			err: &httpError{
				statusCode: 504,
				message:    "Gateway Timeout",
			},
			shouldRetry: true,
		},
		{
			name: "400 error is not transient",
			err: &httpError{
				statusCode: 400,
				message:    "Bad Request",
			},
			shouldRetry: false,
		},
		{
			name: "401 error is not transient",
			err: &httpError{
				statusCode: 401,
				message:    "Unauthorized",
			},
			shouldRetry: false,
		},
		{
			name: "404 error is not transient",
			err: &httpError{
				statusCode: 404,
				message:    "Not Found",
			},
			shouldRetry: false,
		},
		{
			name: "409 error is not transient",
			err: &httpError{
				statusCode: 409,
				message:    "Conflict",
			},
			shouldRetry: false,
		},
		{
			name: "200 success is not transient",
			err: &httpError{
				statusCode: 200,
				message:    "OK",
			},
			shouldRetry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTransientHTTPError(tt.err)
			assert.Equal(t, tt.shouldRetry, result)
		})
	}
}
