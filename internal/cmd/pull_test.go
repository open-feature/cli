package cmd

import (
	"encoding/json"
	"testing"

	"github.com/h2non/gock"
	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/filesystem"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) afero.Fs {
	fs := afero.NewMemMapFs()
	filesystem.SetFileSystem(fs)
	readOsFileAndWriteToMemMap(t, "testdata/empty_manifest.golden", "manifest/path.json", fs)
	return fs
}

func TestPull(t *testing.T) {
	t.Run("pull no provider url", func(t *testing.T) {
		setupTest(t)
		cmd := GetPullCmd()

		// Prepare command arguments
		args := []string{
			"pull",
		}

		cmd.SetArgs(args)

		// Run command
		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "provider URL not set in config")
	})

	t.Run("pull with provider url", func(t *testing.T) {
		fs := setupTest(t)
		defer gock.Off()

		// Mock response in OpenAPI ManifestEnvelope format
		manifestResponse := map[string]any{
			"flags": []map[string]any{
				{
					"key":          "testFlag",
					"type":         "boolean",
					"defaultValue": true,
					"description":  "Test boolean flag",
				},
				{
					"key":          "testFlag2",
					"type":         "string",
					"defaultValue": "string value",
					"description":  "Test string flag",
				},
			},
		}

		// Mock the sync API endpoint
		gock.New("https://example.com").
			Get("/openfeature/v0/manifest").
			Reply(200).
			JSON(manifestResponse)

		cmd := GetPullCmd()

		// global flag exists on root only.
		config.AddRootFlags(cmd)

		// Prepare command arguments - use base URL only
		args := []string{
			"pull",
			"--provider-url", "https://example.com",
			"--manifest", "manifest/path.json",
		}

		cmd.SetArgs(args)

		// Run command
		err := cmd.Execute()
		assert.NoError(t, err)

		// check if the file content is correct
		content, err := afero.ReadFile(fs, "manifest/path.json")
		assert.NoError(t, err)

		var manifestFlags map[string]any
		err = json.Unmarshal(content, &manifestFlags)
		assert.NoError(t, err)

		// Compare actual content with expected flags
		flags := manifestResponse["flags"].([]map[string]any)
		for _, flag := range flags {
			flagKey := flag["key"].(string)
			_, exists := manifestFlags["flags"].(map[string]any)[flagKey]
			assert.True(t, exists, "Flag %s should exist in manifest", flagKey)
		}
	})

	t.Run("error when endpoint returns error", func(t *testing.T) {
		setupTest(t)
		defer gock.Off()

		// Mock error response from sync API
		gock.New("https://example.com").
			Get("/openfeature/v0/manifest").
			Reply(404).
			JSON(map[string]any{
				"error": map[string]any{
					"message": "Not found",
					"status":  404,
				},
			})

		cmd := GetPullCmd()

		// global flag exists on root only.
		config.AddRootFlags(cmd)

		// Prepare command arguments - use base URL only
		args := []string{
			"pull",
			"--provider-url", "https://example.com",
			"--manifest", "manifest/path.json",
		}

		cmd.SetArgs(args)

		// Run command
		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code 404")
	})

	t.Run("pull with .json URL uses LoadFromRemote", func(t *testing.T) {
		fs := setupTest(t)
		defer gock.Off()

		// Mock response - direct file response, not wrapped in OpenAPI format
		flagsResponse := map[string]any{
			"flags": map[string]any{
				"jsonFileFlag": map[string]any{
					"flagType":     "boolean",
					"defaultValue": true,
					"description":  "Flag from JSON file",
				},
			},
		}

		// Mock direct HTTP GET to the file URL (no /openfeature/v0/manifest suffix)
		gock.New("https://example.com").
			Get("/flags.json").
			Reply(200).
			JSON(flagsResponse)

		cmd := GetPullCmd()

		// global flag exists on root only.
		config.AddRootFlags(cmd)

		// Prepare command arguments - URL with .json extension
		args := []string{
			"pull",
			"--provider-url", "https://example.com/flags.json",
			"--manifest", "manifest/path.json",
		}

		cmd.SetArgs(args)

		// Run command
		err := cmd.Execute()
		assert.NoError(t, err)

		// Verify the manifest was written
		content, err := afero.ReadFile(fs, "manifest/path.json")
		assert.NoError(t, err)

		var manifestFlags map[string]any
		err = json.Unmarshal(content, &manifestFlags)
		assert.NoError(t, err)

		// Verify the flag exists in the manifest
		flags := manifestFlags["flags"].(map[string]any)
		_, exists := flags["jsonFileFlag"]
		assert.True(t, exists, "Flag jsonFileFlag should exist in manifest")
	})

	t.Run("pull with .yaml URL uses LoadFromRemote", func(t *testing.T) {
		fs := setupTest(t)
		defer gock.Off()

		// Mock response - direct file response
		flagsResponse := map[string]any{
			"flags": map[string]any{
				"yamlFileFlag": map[string]any{
					"flagType":     "string",
					"defaultValue": "yaml value",
					"description":  "Flag from YAML file",
				},
			},
		}

		// Mock direct HTTP GET to the file URL (no /openfeature/v0/manifest suffix)
		gock.New("https://example.com").
			Get("/flags.yaml").
			Reply(200).
			JSON(flagsResponse)

		cmd := GetPullCmd()

		// global flag exists on root only.
		config.AddRootFlags(cmd)

		// Prepare command arguments - URL with .yaml extension
		args := []string{
			"pull",
			"--provider-url", "https://example.com/flags.yaml",
			"--manifest", "manifest/path.json",
		}

		cmd.SetArgs(args)

		// Run command
		err := cmd.Execute()
		assert.NoError(t, err)

		// Verify the manifest was written
		content, err := afero.ReadFile(fs, "manifest/path.json")
		assert.NoError(t, err)

		var manifestFlags map[string]any
		err = json.Unmarshal(content, &manifestFlags)
		assert.NoError(t, err)

		// Verify the flag exists in the manifest
		flags := manifestFlags["flags"].(map[string]any)
		_, exists := flags["yamlFileFlag"]
		assert.True(t, exists, "Flag yamlFileFlag should exist in manifest")
	})

	t.Run("pull with .yml URL uses LoadFromRemote", func(t *testing.T) {
		fs := setupTest(t)
		defer gock.Off()

		// Mock response - direct file response
		flagsResponse := map[string]any{
			"flags": map[string]any{
				"ymlFileFlag": map[string]any{
					"flagType":     "integer",
					"defaultValue": 42,
					"description":  "Flag from YML file",
				},
			},
		}

		// Mock direct HTTP GET to the file URL (no /openfeature/v0/manifest suffix)
		gock.New("https://example.com").
			Get("/config.yml").
			Reply(200).
			JSON(flagsResponse)

		cmd := GetPullCmd()

		// global flag exists on root only.
		config.AddRootFlags(cmd)

		// Prepare command arguments - URL with .yml extension
		args := []string{
			"pull",
			"--provider-url", "https://example.com/config.yml",
			"--manifest", "manifest/path.json",
		}

		cmd.SetArgs(args)

		// Run command
		err := cmd.Execute()
		assert.NoError(t, err)

		// Verify the manifest was written
		content, err := afero.ReadFile(fs, "manifest/path.json")
		assert.NoError(t, err)

		var manifestFlags map[string]any
		err = json.Unmarshal(content, &manifestFlags)
		assert.NoError(t, err)

		// Verify the flag exists in the manifest
		flags := manifestFlags["flags"].(map[string]any)
		_, exists := flags["ymlFileFlag"]
		assert.True(t, exists, "Flag ymlFileFlag should exist in manifest")
	})

	t.Run("pull with non-file URL uses LoadFromSyncAPI", func(t *testing.T) {
		fs := setupTest(t)
		defer gock.Off()

		// Mock response in OpenAPI ManifestEnvelope format
		manifestResponse := map[string]any{
			"flags": []map[string]any{
				{
					"key":          "syncApiFlag",
					"type":         "boolean",
					"defaultValue": false,
					"description":  "Flag from sync API",
				},
			},
		}

		// Mock the sync API endpoint - note the /openfeature/v0/manifest suffix
		gock.New("https://api.example.com").
			Get("/openfeature/v0/manifest").
			Reply(200).
			JSON(manifestResponse)

		cmd := GetPullCmd()

		// global flag exists on root only.
		config.AddRootFlags(cmd)

		// Prepare command arguments - base URL without file extension
		args := []string{
			"pull",
			"--provider-url", "https://api.example.com",
			"--manifest", "manifest/path.json",
		}

		cmd.SetArgs(args)

		// Run command
		err := cmd.Execute()
		assert.NoError(t, err)

		// Verify the manifest was written
		content, err := afero.ReadFile(fs, "manifest/path.json")
		assert.NoError(t, err)

		var manifestFlags map[string]any
		err = json.Unmarshal(content, &manifestFlags)
		assert.NoError(t, err)

		// Verify the flag exists in the manifest
		flags := manifestFlags["flags"].(map[string]any)
		_, exists := flags["syncApiFlag"]
		assert.True(t, exists, "Flag syncApiFlag should exist in manifest")
	})

	t.Run("backward compatibility with deprecated --flag-source-url", func(t *testing.T) {
		fs := setupTest(t)
		defer gock.Off()

		// Mock response in OpenAPI ManifestEnvelope format
		manifestResponse := map[string]any{
			"flags": []map[string]any{
				{
					"key":          "backwardCompatFlag",
					"type":         "boolean",
					"defaultValue": true,
					"description":  "Test backward compatibility",
				},
			},
		}

		// Mock the sync API endpoint
		gock.New("https://example.com").
			Get("/openfeature/v0/manifest").
			Reply(200).
			JSON(manifestResponse)

		cmd := GetPullCmd()

		// global flag exists on root only.
		config.AddRootFlags(cmd)

		// Use the deprecated flag to test backward compatibility
		args := []string{
			"pull",
			"--flag-source-url", "https://example.com",
			"--manifest", "manifest/path.json",
		}

		cmd.SetArgs(args)

		// Run command - should work but show deprecation warning
		err := cmd.Execute()
		assert.NoError(t, err)

		// Verify the manifest was written correctly
		content, err := afero.ReadFile(fs, "manifest/path.json")
		assert.NoError(t, err)

		var manifestFlags map[string]any
		err = json.Unmarshal(content, &manifestFlags)
		assert.NoError(t, err)

		// Verify the flag exists in the manifest
		flags := manifestFlags["flags"].(map[string]any)
		_, exists := flags["backwardCompatFlag"]
		assert.True(t, exists, "Flag backwardCompatFlag should exist in manifest")
	})
}
