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

func setupTest(t *testing.T) (afero.Fs) {
	fs := afero.NewMemMapFs()
	filesystem.SetFileSystem(fs)
	readOsFileAndWriteToMemMap(t, "testdata/empty_manifest.golden", "manifest/path.json", fs)
	return fs
}

func TestPull(t *testing.T) {
	t.Run("pull no flag source url", func(t *testing.T) {
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
		assert.Contains(t, err.Error(), "flagSourceUrl not set in config")
	})

	t.Run("pull with flag source url", func(t *testing.T) {
		fs := setupTest(t)
		defer gock.Off()

		flags := []map[string]any{
			{"key": "testFlag", "type": "boolean", "defaultValue": true},
			{"key": "testFlag2", "type": "string", "defaultValue": "string value"},
		}

		gock.New("https://example.com/flags").
			Get("/").
			Reply(200).
			JSON(flags)
			
		cmd := GetPullCmd()

		// global flag exists on root only.
		config.AddRootFlags(cmd)

		// Prepare command arguments
		args := []string{
			"pull",
			"--flag-source-url", "https://example.com/flags",
			"--manifest", "manifest/path.json",
		}

		cmd.SetArgs(args)


		// Run command
		err := cmd.Execute()
		assert.NoError(t, err)

		// check if the file content is correct
		content, err := afero.ReadFile(fs, "manifest/path.json")
		assert.NoError(t, err)
		
		var manifestFlags map[string]interface{}
		err = json.Unmarshal(content, &manifestFlags)
		assert.NoError(t, err)

		// Compare actual content with expected flags
		for _, flag := range flags {
			flagKey := flag["key"].(string)
			_, exists := manifestFlags["flags"].(map[string]interface{})[flagKey]
			assert.True(t, exists, "Flag %s should exist in manifest", flagKey)
		}
	})

	t.Run("error when endpoint returns error", func(t *testing.T) {
		setupTest(t)
		defer gock.Off()

		gock.New("https://example.com/flags").
			Get("/").
			Reply(404)

		cmd := GetPullCmd()

		// global flag exists on root only.
		config.AddRootFlags(cmd)

		// Prepare command arguments
		args := []string{
			"pull",
			"--flag-source-url", "https://example.com/flags",
			"--manifest", "manifest/path.json",
		}

		cmd.SetArgs(args)

		// Run command
		err := cmd.Execute()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Received error response from flag source")
	})
}
