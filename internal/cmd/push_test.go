package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/open-feature/cli/internal/filesystem"

	"github.com/spf13/afero"

	"github.com/stretchr/testify/assert"
)

func setupPushTest(t *testing.T) afero.Fs {
	fs := afero.NewMemMapFs()
	filesystem.SetFileSystem(fs)
	// Copy test manifest to the filesystem
	readOsFileAndWriteToMemMap(t, "testdata/success_manifest.golden", "flags.json", fs)
	return fs
}

func TestPush(t *testing.T) {
	t.Run("push without destination URL", func(t *testing.T) {
		setupPushTest(t)
		cmd := GetPushCmd()

		args := []string{
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "flag source URL is required")
	})

	t.Run("smart push creates new flags", func(t *testing.T) {
		setupPushTest(t)
		defer gock.Off()

		// Mock GET request to fetch remote flags (returns empty list)
		emptyFlags := []map[string]interface{}{}
		gock.New("http://localhost:8080").
			Get("/openfeature/v0/manifest").
			Reply(200).
			JSON(map[string]interface{}{
				"flags": emptyFlags,
			})

		// Mock individual POST requests for each flag in the manifest
		// Since remote has no flags, all local flags will be created
		flagKeys := []string{"enableFeatureA", "usernameMaxLength", "greetingMessage", "discountPercentage", "themeCustomization"}
		for _, flagKey := range flagKeys {
			gock.New("http://localhost:8080").
				Post("/openfeature/v0/manifest/flags").
				MatchType("application/json").
				MatchHeader("Content-Type", "application/json").
				Reply(201).
				JSON(map[string]interface{}{
					"flag": map[string]interface{}{
						"key": flagKey,
					},
					"updatedAt": "2024-03-02T09:45:03.000Z",
				})
		}

		cmd := GetPushCmd()

		args := []string{
			"--flag-source-url", "http://localhost:8080/openfeature/v0/manifest",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.NoError(t, err)

		// Verify that all mocked requests were made
		assert.True(t, gock.IsDone(), "Not all expected HTTP requests were made")
	})

	t.Run("smart push updates existing flags", func(t *testing.T) {
		setupPushTest(t)
		defer gock.Off()

		// Mock GET request to fetch remote flags (all flags already exist)
		flagKeys := []string{"enableFeatureA", "usernameMaxLength", "greetingMessage", "discountPercentage", "themeCustomization"}
		remoteFlags := make([]map[string]interface{}, 0)
		for _, flagKey := range flagKeys {
			remoteFlags = append(remoteFlags, map[string]interface{}{
				"key":          flagKey,
				"type":         "boolean",
				"defaultValue": false,
			})
		}
		gock.New("https://api.example.com").
			Get("/openfeature/v0/manifest").
			Reply(200).
			JSON(map[string]interface{}{
				"flags": remoteFlags,
			})

		// Mock individual PUT requests for each flag
		for _, flagKey := range flagKeys {
			gock.New("https://api.example.com").
				Put("/openfeature/v0/manifest/flags/"+flagKey).
				MatchType("application/json").
				MatchHeader("Content-Type", "application/json").
				Reply(200).
				JSON(map[string]interface{}{
					"flag": map[string]interface{}{
						"key": flagKey,
					},
					"updatedAt": "2024-03-02T09:45:03.000Z",
				})
		}

		cmd := GetPushCmd()

		args := []string{
			"--flag-source-url", "https://api.example.com/openfeature/v0/manifest",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.NoError(t, err)

		assert.True(t, gock.IsDone(), "Not all expected HTTP requests were made")
	})

	t.Run("smart push with mixed create and update", func(t *testing.T) {
		setupPushTest(t)
		defer gock.Off()

		// Mock GET request - some flags exist, some don't
		gock.New("https://api.example.com").
			Get("/openfeature/v0/manifest").
			Reply(200).
			JSON(map[string]interface{}{
				"flags": []map[string]interface{}{
					{
						"key":          "enableFeatureA",
						"type":         "boolean",
						"defaultValue": false,
					},
					{
						"key":          "usernameMaxLength",
						"type":         "integer",
						"defaultValue": 10,
					},
				},
			})

		// Mock PUT requests for existing flags
		existingFlags := []string{"enableFeatureA", "usernameMaxLength"}
		for _, flagKey := range existingFlags {
			gock.New("https://api.example.com").
				Put("/openfeature/v0/manifest/flags/" + flagKey).
				MatchType("application/json").
				Reply(200).
				JSON(map[string]interface{}{
					"flag": map[string]interface{}{
						"key": flagKey,
					},
					"updatedAt": "2024-03-02T09:45:03.000Z",
				})
		}

		// Mock POST requests for new flags
		newFlags := []string{"greetingMessage", "discountPercentage", "themeCustomization"}
		for _, flagKey := range newFlags {
			gock.New("https://api.example.com").
				Post("/openfeature/v0/manifest/flags").
				MatchType("application/json").
				Reply(201).
				JSON(map[string]interface{}{
					"flag": map[string]interface{}{
						"key": flagKey,
					},
					"updatedAt": "2024-03-02T09:45:03.000Z",
				})
		}

		cmd := GetPushCmd()

		args := []string{
			"--flag-source-url", "https://api.example.com/openfeature/v0/manifest",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.NoError(t, err)

		assert.True(t, gock.IsDone(), "Not all expected HTTP requests were made")
	})

	t.Run("push with authentication token", func(t *testing.T) {
		setupPushTest(t)
		defer gock.Off()

		// Mock GET request with auth header
		emptyFlags := []map[string]interface{}{}
		gock.New("https://api.example.com").
			Get("/openfeature/v0/manifest").
			MatchHeader("Authorization", "Bearer secret-token").
			Reply(200).
			JSON(map[string]interface{}{
				"flags": emptyFlags,
			})

		// Mock individual POST requests with auth header for each flag
		flagKeys := []string{"enableFeatureA", "usernameMaxLength", "greetingMessage", "discountPercentage", "themeCustomization"}
		for _, flagKey := range flagKeys {
			gock.New("https://api.example.com").
				Post("/openfeature/v0/manifest/flags").
				MatchType("application/json").
				MatchHeader("Authorization", "Bearer secret-token").
				Reply(201).
				JSON(map[string]interface{}{
					"flag": map[string]interface{}{
						"key": flagKey,
					},
					"updatedAt": "2024-03-02T09:45:03.000Z",
				})
		}

		cmd := GetPushCmd()

		args := []string{
			"--flag-source-url", "https://api.example.com/openfeature/v0/manifest",
			"--auth-token", "secret-token",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.NoError(t, err)

		assert.True(t, gock.IsDone(), "Not all expected HTTP requests were made")
	})

	t.Run("push with dry run", func(t *testing.T) {
		setupPushTest(t)
		defer gock.Off()

		// Mock GET request to fetch remote flags (returns some existing flags)
		gock.New("https://api.example.com").
			Get("/openfeature/v0/manifest").
			Reply(200).
			JSON(map[string]any{
				"flags": []map[string]any{
					{
						"key":          "enableFeatureA",
						"type":         "boolean",
						"defaultValue": false, // Different from local
						"description":  "Old description",
					},
					{
						"key":          "usernameMaxLength",
						"type":         "integer",
						"defaultValue": 10,
						"description":  "Max username length",
					},
				},
			})

		// Dry run should NOT make any POST or PUT requests
		// If any POST/PUT requests are made, gock will fail to match and the test will fail

		cmd := GetPushCmd()

		args := []string{
			"--flag-source-url", "https://api.example.com",
			"--dry-run",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.NoError(t, err)

		// Verify that only the GET request was made (no POST/PUT)
		// gock.IsDone() returns true if all mocked requests were consumed
		assert.True(t, gock.IsDone(), "Should only make GET request, not POST/PUT")
	})

	t.Run("push with file scheme returns error", func(t *testing.T) {
		setupPushTest(t)

		cmd := GetPushCmd()

		args := []string{
			"--flag-source-url", "file:///local/path/flags.json",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file:// scheme is not supported for push")
		assert.Contains(t, err.Error(), "Use standard shell commands")
	})

	t.Run("push with unsupported scheme returns error", func(t *testing.T) {
		setupPushTest(t)

		cmd := GetPushCmd()

		args := []string{
			"--flag-source-url", "ftp://example.com/flags",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported URL scheme: ftp")
		assert.Contains(t, err.Error(), "Supported schemes are http:// and https://")
	})

	t.Run("error when fetch returns 404", func(t *testing.T) {
		setupPushTest(t)
		defer gock.Off()

		// Mock GET request returning 404
		gock.New("https://api.example.com").
			Get("/openfeature/v0/manifest").
			Reply(404).
			JSON(map[string]interface{}{
				"error": map[string]interface{}{
					"message": "Not Found",
					"status":  404,
				},
			})

		cmd := GetPushCmd()

		args := []string{
			"--flag-source-url", "https://api.example.com/openfeature/v0/manifest",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch remote flags")
	})

	t.Run("error when create endpoint returns 404", func(t *testing.T) {
		setupPushTest(t)
		defer gock.Off()

		// Mock GET request (empty flags)
		emptyFlags := []map[string]interface{}{}
		gock.New("https://api.example.com").
			Get("/openfeature/v0/manifest").
			Reply(200).
			JSON(map[string]interface{}{
				"flags": emptyFlags,
			})

		// Mock a 404 error for all flag creation requests
		gock.New("https://api.example.com").
			Post("/openfeature/v0/manifest/flags").
			Persist(). // Apply to all requests
			Reply(404).
			JSON(map[string]interface{}{
				"error": map[string]interface{}{
					"message": "Not Found",
					"status":  404,
				},
			})

		cmd := GetPushCmd()

		args := []string{
			"--flag-source-url", "https://api.example.com/openfeature/v0/manifest",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create flag")
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("error when endpoint returns 500", func(t *testing.T) {
		setupPushTest(t)
		defer gock.Off()

		// Mock GET request (empty flags)
		emptyFlags := []map[string]interface{}{}
		gock.New("https://api.example.com").
			Get("/openfeature/v0/manifest").
			Reply(200).
			JSON(map[string]interface{}{
				"flags": emptyFlags,
			})

		// Mock a 500 error for all flag creation requests
		gock.New("https://api.example.com").
			Post("/openfeature/v0/manifest/flags").
			Persist(). // Apply to all requests
			Reply(500).
			JSON(map[string]interface{}{
				"error": map[string]interface{}{
					"message": "Internal Server Error",
					"status":  500,
				},
			})

		cmd := GetPushCmd()

		args := []string{
			"--flag-source-url", "https://api.example.com/openfeature/v0/manifest",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create flag")
		assert.Contains(t, err.Error(), "500")
	})

	t.Run("push validates request body format", func(t *testing.T) {
		setupPushTest(t)
		defer gock.Off()

		// Mock GET request (empty flags)
		emptyFlags := []map[string]interface{}{}
		gock.New("https://api.example.com").
			Get("/openfeature/v0/manifest").
			Reply(200).
			JSON(map[string]interface{}{
				"flags": emptyFlags,
			})

		// Mock all POST requests and validate that the right fields are present
		// We can't predict the order flags will be processed, so we validate structure not content
		requestCount := 0
		gock.New("https://api.example.com").
			Post("/openfeature/v0/manifest/flags").
			MatchType("application/json").
			SetMatcher(gock.NewMatcher()).
			AddMatcher(func(req *http.Request, ereq *gock.Request) (bool, error) {
				// Verify the request has the required fields: key, type, defaultValue
				var body map[string]interface{}
				decoder := json.NewDecoder(req.Body)
				if err := decoder.Decode(&body); err != nil {
					return false, err
				}
				// Check required fields exist
				if _, ok := body["key"]; !ok {
					return false, fmt.Errorf("missing key field")
				}
				if _, ok := body["type"]; !ok {
					return false, fmt.Errorf("missing type field")
				}
				if _, ok := body["defaultValue"]; !ok {
					return false, fmt.Errorf("missing defaultValue field")
				}
				requestCount++
				return true, nil
			}).
			Persist().
			Reply(201).
			JSON(map[string]interface{}{
				"flag": map[string]interface{}{
					"key": "test",
				},
				"updatedAt": "2024-03-02T09:45:03.000Z",
			})

		cmd := GetPushCmd()

		args := []string{
			"--flag-source-url", "https://api.example.com/openfeature/v0/manifest",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.NoError(t, err)
		assert.Equal(t, 5, requestCount, "Expected 5 flag creation requests")
	})

	t.Run("error when manifest file does not exist", func(t *testing.T) {
		setupPushTest(t)

		cmd := GetPushCmd()

		args := []string{
			"--flag-source-url", "https://api.example.com/flags",
			"--manifest", "nonexistent.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error loading manifest")
	})

	t.Run("push with manifest containing flag without default value", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		filesystem.SetFileSystem(fs)

		// Create a manifest with a flag missing a default value
		invalidManifest := `{
			"flags": {
				"testFlag": {
					"flagType": "boolean",
					"description": "Test flag without default"
				}
			}
		}`
		err := afero.WriteFile(fs, "invalid.json", []byte(invalidManifest), 0644)
		assert.NoError(t, err)

		cmd := GetPushCmd()

		args := []string{
			"--flag-source-url", "https://api.example.com/flags",
			"--manifest", "invalid.json",
		}
		cmd.SetArgs(args)

		err = cmd.Execute()
		assert.Error(t, err)
		// The error message is from manifest validation
		assert.Contains(t, err.Error(), "defaultValue is required")
	})
}
