package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/filesystem"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManifestAddCmd(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		existingManifest string
		expectedError  string
		validateResult func(t *testing.T, fs afero.Fs)
	}{
		{
			name: "add boolean flag to empty manifest",
			args: []string{
				"add", "new-feature",
				"--default-value", "true",
				"--description", "A new feature flag",
			},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {}
			}`,
			validateResult: func(t *testing.T, fs afero.Fs) {
				content, err := afero.ReadFile(fs, "flags.json")
				require.NoError(t, err)

				var manifest map[string]interface{}
				err = json.Unmarshal(content, &manifest)
				require.NoError(t, err)

				flags := manifest["flags"].(map[string]interface{})
				assert.Contains(t, flags, "new-feature")

				flag := flags["new-feature"].(map[string]interface{})
				assert.Equal(t, "boolean", flag["flagType"])
				assert.Equal(t, true, flag["defaultValue"])
				assert.Equal(t, "A new feature flag", flag["description"])
			},
		},
		{
			name: "add string flag",
			args: []string{
				"add", "welcome-message",
				"--type", "string",
				"--default-value", "Hello World",
			},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {}
			}`,
			validateResult: func(t *testing.T, fs afero.Fs) {
				content, err := afero.ReadFile(fs, "flags.json")
				require.NoError(t, err)

				var manifest map[string]interface{}
				err = json.Unmarshal(content, &manifest)
				require.NoError(t, err)

				flags := manifest["flags"].(map[string]interface{})
				flag := flags["welcome-message"].(map[string]interface{})
				assert.Equal(t, "string", flag["flagType"])
				assert.Equal(t, "Hello World", flag["defaultValue"])
			},
		},
		{
			name: "add integer flag",
			args: []string{
				"add", "max-retries",
				"--type", "integer",
				"--default-value", "5",
			},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {}
			}`,
			validateResult: func(t *testing.T, fs afero.Fs) {
				content, err := afero.ReadFile(fs, "flags.json")
				require.NoError(t, err)

				var manifest map[string]interface{}
				err = json.Unmarshal(content, &manifest)
				require.NoError(t, err)

				flags := manifest["flags"].(map[string]interface{})
				flag := flags["max-retries"].(map[string]interface{})
				assert.Equal(t, "integer", flag["flagType"])
				// JSON unmarshaling converts numbers to float64
				assert.Equal(t, float64(5), flag["defaultValue"])
			},
		},
		{
			name: "add float flag",
			args: []string{
				"add", "discount-rate",
				"--type", "float",
				"--default-value", "0.15",
			},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {}
			}`,
			validateResult: func(t *testing.T, fs afero.Fs) {
				content, err := afero.ReadFile(fs, "flags.json")
				require.NoError(t, err)

				var manifest map[string]interface{}
				err = json.Unmarshal(content, &manifest)
				require.NoError(t, err)

				flags := manifest["flags"].(map[string]interface{})
				flag := flags["discount-rate"].(map[string]interface{})
				assert.Equal(t, "float", flag["flagType"])
				assert.Equal(t, 0.15, flag["defaultValue"])
			},
		},
		{
			name: "add object flag",
			args: []string{
				"add", "config",
				"--type", "object",
				"--default-value", `{"key":"value","nested":{"prop":123}}`,
			},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {}
			}`,
			validateResult: func(t *testing.T, fs afero.Fs) {
				content, err := afero.ReadFile(fs, "flags.json")
				require.NoError(t, err)

				var manifest map[string]interface{}
				err = json.Unmarshal(content, &manifest)
				require.NoError(t, err)

				flags := manifest["flags"].(map[string]interface{})
				flag := flags["config"].(map[string]interface{})
				assert.Equal(t, "object", flag["flagType"])

				defaultVal := flag["defaultValue"].(map[string]interface{})
				assert.Equal(t, "value", defaultVal["key"])
				nested := defaultVal["nested"].(map[string]interface{})
				assert.Equal(t, float64(123), nested["prop"])
			},
		},
		{
			name: "error on duplicate flag name",
			args: []string{
				"add", "existing-flag",
				"--default-value", "true",
			},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {
					"existing-flag": {
						"flagType": "boolean",
						"defaultValue": false,
						"description": "An existing flag"
					}
				}
			}`,
			expectedError: "flag 'existing-flag' already exists in the manifest",
		},
		{
			name: "error on missing default value",
			args: []string{
				"add", "new-flag",
				"--no-input",
			},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {}
			}`,
			expectedError: "--default-value is required",
		},
		{
			name: "error on invalid type",
			args: []string{
				"add", "new-flag",
				"--type", "invalid",
				"--default-value", "test",
			},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {}
			}`,
			expectedError: "invalid flag type: unknown flag type: invalid",
		},
		{
			name: "error on type mismatch - boolean",
			args: []string{
				"add", "new-flag",
				"--type", "boolean",
				"--default-value", "not-a-bool",
			},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {}
			}`,
			expectedError: "invalid default value for type boolean: boolean value must be 'true' or 'false', got 'not-a-bool'",
		},
		{
			name: "error on type mismatch - integer",
			args: []string{
				"add", "new-flag",
				"--type", "integer",
				"--default-value", "not-an-int",
			},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {}
			}`,
			expectedError: "invalid default value for type integer: invalid integer value: not-an-int",
		},
		{
			name: "error on type mismatch - float",
			args: []string{
				"add", "new-flag",
				"--type", "float",
				"--default-value", "not-a-float",
			},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {}
			}`,
			expectedError: "invalid default value for type float: invalid float value: not-a-float",
		},
		{
			name: "error on type mismatch - object",
			args: []string{
				"add", "new-flag",
				"--type", "object",
				"--default-value", "not-valid-json",
			},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {}
			}`,
			expectedError: "invalid default value for type object: invalid JSON object:",
		},
		{
			name: "add flag to existing manifest with flags",
			args: []string{
				"add", "new-flag",
				"--default-value", "false",
			},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {
					"existing-flag": {
						"flagType": "string",
						"defaultValue": "test",
						"description": "An existing flag"
					}
				}
			}`,
			validateResult: func(t *testing.T, fs afero.Fs) {
				content, err := afero.ReadFile(fs, "flags.json")
				require.NoError(t, err)

				var manifest map[string]interface{}
				err = json.Unmarshal(content, &manifest)
				require.NoError(t, err)

				flags := manifest["flags"].(map[string]interface{})
				assert.Len(t, flags, 2)
				assert.Contains(t, flags, "existing-flag")
				assert.Contains(t, flags, "new-flag")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			fs := afero.NewMemMapFs()
			filesystem.SetFileSystem(fs)

			// Create existing manifest if provided
			if tt.existingManifest != "" {
				err := afero.WriteFile(fs, "flags.json", []byte(tt.existingManifest), 0644)
				require.NoError(t, err)
			}

			// Create command and execute
			cmd := GetManifestCmd()
			config.AddRootFlags(cmd)

			// Set args with manifest path
			args := append(tt.args, "-m", "flags.json")
			cmd.SetArgs(args)

			// Execute command
			err := cmd.Execute()

			// Validate
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				if tt.validateResult != nil {
					tt.validateResult(t, fs)
				}
			}
		})
	}
}

func TestManifestAddCmd_CreateNewManifest(t *testing.T) {
	// Setup
	fs := afero.NewMemMapFs()
	filesystem.SetFileSystem(fs)

	// Don't create any existing manifest

	// Create command and execute
	cmd := GetManifestCmd()
	config.AddRootFlags(cmd)

	cmd.SetArgs([]string{
		"add", "first-flag",
		"--default-value", "true",
		"--description", "The first flag in a new manifest",
		"-m", "flags.json",
	})

	// Execute command
	err := cmd.Execute()
	require.NoError(t, err)

	// Validate the new manifest was created
	exists, err := afero.Exists(fs, "flags.json")
	require.NoError(t, err)
	assert.True(t, exists)

	content, err := afero.ReadFile(fs, "flags.json")
	require.NoError(t, err)

	var manifest map[string]interface{}
	err = json.Unmarshal(content, &manifest)
	require.NoError(t, err)

	// Check schema is present
	assert.Contains(t, manifest, "$schema")

	// Check flag was added
	flags := manifest["flags"].(map[string]interface{})
	assert.Contains(t, flags, "first-flag")

	flag := flags["first-flag"].(map[string]interface{})
	assert.Equal(t, "boolean", flag["flagType"])
	assert.Equal(t, true, flag["defaultValue"])
	assert.Equal(t, "The first flag in a new manifest", flag["description"])
}

func TestManifestAddCmd_DisplaysListAfterAdd(t *testing.T) {
	// Setup
	fs := afero.NewMemMapFs()
	filesystem.SetFileSystem(fs)

	// Create existing manifest with one flag
	existingManifest := `{
		"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
		"flags": {
			"existing-flag": {
				"flagType": "string",
				"defaultValue": "test",
				"description": "An existing flag"
			}
		}
	}`
	err := afero.WriteFile(fs, "flags.json", []byte(existingManifest), 0644)
	require.NoError(t, err)

	// Enable pterm output and capture it
	pterm.EnableOutput()
	defer pterm.DisableOutput()

	buf := &bytes.Buffer{}
	oldStdout := pterm.DefaultTable.Writer
	oldSection := pterm.DefaultSection.Writer
	oldInfo := pterm.Info.Writer
	oldSuccess := pterm.Success.Writer
	pterm.DefaultTable.Writer = buf
	pterm.DefaultSection.Writer = buf
	pterm.Info.Writer = buf
	pterm.Success.Writer = buf
	defer func() {
		pterm.DefaultTable.Writer = oldStdout
		pterm.DefaultSection.Writer = oldSection
		pterm.Info.Writer = oldInfo
		pterm.Success.Writer = oldSuccess
	}()

	// Create command and execute
	cmd := GetManifestCmd()
	config.AddRootFlags(cmd)

	cmd.SetArgs([]string{
		"add", "new-flag",
		"--default-value", "true",
		"--description", "A new flag",
		"-m", "flags.json",
	})

	// Execute command
	err = cmd.Execute()
	require.NoError(t, err)

	// Validate output contains list of all flags
	output := buf.String()
	assert.Contains(t, output, "existing-flag", "Output should contain existing flag")
	assert.Contains(t, output, "new-flag", "Output should contain newly added flag")
	assert.Contains(t, output, "(2)", "Output should show total count of 2 flags")
	assert.Contains(t, output, "string", "Output should show flag types")
	assert.Contains(t, output, "boolean", "Output should show flag types")
}

