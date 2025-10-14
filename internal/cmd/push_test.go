package cmd

import (
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
		assert.Contains(t, err.Error(), "flag destination URL is required")
	})

	t.Run("push with HTTP destination URL using POST", func(t *testing.T) {
		setupPushTest(t)
		defer gock.Off()

		// Mock the HTTP POST request
		gock.New("http://localhost:8080").
			Post("/flags").
			MatchType("application/json").
			MatchHeader("Content-Type", "application/json").
			Reply(200).
			JSON(map[string]string{"status": "success"})

		cmd := GetPushCmd()

		args := []string{
			"--flag-destination-url", "http://localhost:8080/flags",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.NoError(t, err)

		// Verify that all mocked requests were made
		assert.True(t, gock.IsDone(), "Not all expected HTTP requests were made")
	})

	t.Run("push with HTTPS destination URL using PUT", func(t *testing.T) {
		setupPushTest(t)
		defer gock.Off()

		// Mock the HTTPS PUT request
		gock.New("https://api.example.com").
			Put("/flags/my-app").
			MatchType("application/json").
			MatchHeader("Content-Type", "application/json").
			Reply(200).
			JSON(map[string]string{"status": "success"})

		cmd := GetPushCmd()

		args := []string{
			"--flag-destination-url", "https://api.example.com/flags/my-app",
			"--method", "PUT",
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

		// Mock HTTPS request with auth header
		gock.New("https://api.example.com").
			Post("/flags").
			MatchType("application/json").
			MatchHeader("Authorization", "Bearer secret-token").
			Reply(200).
			JSON(map[string]string{"status": "success"})

		cmd := GetPushCmd()

		args := []string{
			"--flag-destination-url", "https://api.example.com/flags",
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
		// No gock setup - we're verifying no requests are made

		cmd := GetPushCmd()

		args := []string{
			"--flag-destination-url", "https://api.example.com/flags",
			"--dry-run",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.NoError(t, err)
		// Test passes if no error - dry run should not make any HTTP requests
	})

	t.Run("push with file scheme returns error", func(t *testing.T) {
		setupPushTest(t)

		cmd := GetPushCmd()

		args := []string{
			"--flag-destination-url", "file:///local/path/flags.json",
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
			"--flag-destination-url", "ftp://example.com/flags",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported URL scheme: ftp")
		assert.Contains(t, err.Error(), "Supported schemes are http:// and https://")
	})

	t.Run("error when endpoint returns 404", func(t *testing.T) {
		setupPushTest(t)
		defer gock.Off()

		gock.New("https://api.example.com").
			Post("/flags").
			Reply(404).
			BodyString("Not Found")

		cmd := GetPushCmd()

		args := []string{
			"--flag-destination-url", "https://api.example.com/flags",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "received error response from destination")
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("error when endpoint returns 500", func(t *testing.T) {
		setupPushTest(t)
		defer gock.Off()

		gock.New("https://api.example.com").
			Post("/flags").
			Reply(500).
			BodyString("Internal Server Error")

		cmd := GetPushCmd()

		args := []string{
			"--flag-destination-url", "https://api.example.com/flags",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "received error response from destination")
		assert.Contains(t, err.Error(), "500")
	})

	t.Run("push validates request body format", func(t *testing.T) {
		setupPushTest(t)
		defer gock.Off()

		// Expected request body structure
		expectedBody := map[string]interface{}{
			"flags": []interface{}{
				map[string]interface{}{
					"key":          "discountPercentage",
					"type":         "float",
					"description":  "Discount percentage applied to purchases.",
					"defaultValue": 0.15,
				},
				map[string]interface{}{
					"key":          "enableFeatureA",
					"type":         "boolean",
					"description":  "Controls whether Feature A is enabled.",
					"defaultValue": false,
				},
				map[string]interface{}{
					"key":          "greetingMessage",
					"type":         "string",
					"description":  "The message to use for greeting users.",
					"defaultValue": "Hello there!",
				},
				map[string]interface{}{
					"key":         "themeCustomization",
					"type":        "object",
					"description": "Allows customization of theme colors.",
					"defaultValue": map[string]interface{}{
						"primaryColor":   "#007bff",
						"secondaryColor": "#6c757d",
					},
				},
				map[string]interface{}{
					"key":          "usernameMaxLength",
					"type":         "integer",
					"description":  "Maximum allowed length for usernames.",
					"defaultValue": float64(50), // JSON unmarshals numbers as float64
				},
			},
		}

		gock.New("https://api.example.com").
			Post("/flags").
			MatchType("application/json").
			JSON(expectedBody).
			Reply(200).
			JSON(map[string]string{"status": "success"})

		cmd := GetPushCmd()

		args := []string{
			"--flag-destination-url", "https://api.example.com/flags",
			"--manifest", "flags.json",
		}
		cmd.SetArgs(args)

		err := cmd.Execute()
		assert.NoError(t, err)
		assert.True(t, gock.IsDone(), "Request body did not match expected structure")
	})

	t.Run("error when manifest file does not exist", func(t *testing.T) {
		setupPushTest(t)

		cmd := GetPushCmd()

		args := []string{
			"--flag-destination-url", "https://api.example.com/flags",
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
			"--flag-destination-url", "https://api.example.com/flags",
			"--manifest", "invalid.json",
		}
		cmd.SetArgs(args)

		err = cmd.Execute()
		assert.Error(t, err)
		// The error message is from manifest validation
		assert.Contains(t, err.Error(), "defaultValue is required")
	})
}
