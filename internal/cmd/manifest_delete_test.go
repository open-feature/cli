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

func TestManifestDeleteCmd(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		existingManifest string
		expectedError    string
		validateResult   func(t *testing.T, fs afero.Fs)
	}{
		{
			name: "delete existing flag from manifest with multiple flags",
			args: []string{"delete", "flag-to-delete"},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {
					"flag-to-keep": {
						"flagType": "boolean",
						"defaultValue": true,
						"description": "This flag should remain"
					},
					"flag-to-delete": {
						"flagType": "string",
						"defaultValue": "remove me",
						"description": "This flag should be deleted"
					},
					"another-flag": {
						"flagType": "integer",
						"defaultValue": 42,
						"description": "Another flag to keep"
					}
				}
			}`,
			validateResult: func(t *testing.T, fs afero.Fs) {
				content, err := afero.ReadFile(fs, "flags.json")
				require.NoError(t, err)

				var manifest map[string]any
				err = json.Unmarshal(content, &manifest)
				require.NoError(t, err)

				flags := manifest["flags"].(map[string]any)
				assert.Len(t, flags, 2, "Should have 2 flags remaining")
				assert.Contains(t, flags, "flag-to-keep")
				assert.Contains(t, flags, "another-flag")
				assert.NotContains(t, flags, "flag-to-delete", "Deleted flag should not be present")
			},
		},
		{
			name: "delete the only flag in manifest",
			args: []string{"delete", "only-flag"},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {
					"only-flag": {
						"flagType": "boolean",
						"defaultValue": false,
						"description": "The only flag"
					}
				}
			}`,
			validateResult: func(t *testing.T, fs afero.Fs) {
				content, err := afero.ReadFile(fs, "flags.json")
				require.NoError(t, err)

				var manifest map[string]any
				err = json.Unmarshal(content, &manifest)
				require.NoError(t, err)

				flags := manifest["flags"].(map[string]any)
				assert.Len(t, flags, 0, "Should have no flags remaining")
			},
		},
		{
			name: "error on non-existent flag",
			args: []string{"delete", "non-existent-flag"},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {
					"existing-flag": {
						"flagType": "boolean",
						"defaultValue": true,
						"description": "An existing flag"
					}
				}
			}`,
			expectedError: "flag 'non-existent-flag' not found in manifest",
		},
		{
			name: "error on empty manifest",
			args: []string{"delete", "any-flag"},
			existingManifest: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {}
			}`,
			expectedError: "flag 'any-flag' not found in manifest",
		},
		{
			name:          "error on missing manifest file",
			args:          []string{"delete", "any-flag"},
			expectedError: "manifest file does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			fs := afero.NewMemMapFs()
			filesystem.SetFileSystem(fs)

			// Create existing manifest if provided
			if tt.existingManifest != "" {
				err := afero.WriteFile(fs, "flags.json", []byte(tt.existingManifest), 0o644)
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

func TestManifestDeleteCmd_ManifestUnchangedOnError(t *testing.T) {
	// Setup
	fs := afero.NewMemMapFs()
	filesystem.SetFileSystem(fs)

	// Create manifest with one flag
	originalManifest := `{
		"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
		"flags": {
			"existing-flag": {
				"flagType": "boolean",
				"defaultValue": true,
				"description": "An existing flag"
			}
		}
	}`
	err := afero.WriteFile(fs, "flags.json", []byte(originalManifest), 0o644)
	require.NoError(t, err)

	// Try to delete a non-existent flag
	cmd := GetManifestCmd()
	config.AddRootFlags(cmd)
	cmd.SetArgs([]string{"delete", "non-existent-flag", "-m", "flags.json"})

	// Execute command - should fail
	err = cmd.Execute()
	assert.Error(t, err)

	// Verify manifest is unchanged
	content, err := afero.ReadFile(fs, "flags.json")
	require.NoError(t, err)

	var manifest map[string]any
	err = json.Unmarshal(content, &manifest)
	require.NoError(t, err)

	flags := manifest["flags"].(map[string]any)
	assert.Len(t, flags, 1, "Manifest should still have 1 flag")
	assert.Contains(t, flags, "existing-flag", "Original flag should still exist")
}

func TestManifestDeleteCmd_DisplaysSuccessMessage(t *testing.T) {
	// Setup
	fs := afero.NewMemMapFs()
	filesystem.SetFileSystem(fs)

	// Create manifest with two flags
	existingManifest := `{
		"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
		"flags": {
			"flag-1": {
				"flagType": "boolean",
				"defaultValue": true,
				"description": "First flag"
			},
			"flag-2": {
				"flagType": "string",
				"defaultValue": "test",
				"description": "Second flag"
			}
		}
	}`
	err := afero.WriteFile(fs, "flags.json", []byte(existingManifest), 0o644)
	require.NoError(t, err)

	// Enable pterm output and capture it
	pterm.EnableOutput()
	defer pterm.DisableOutput()

	buf := &bytes.Buffer{}
	oldSuccess := pterm.Success.Writer
	pterm.Success.Writer = buf
	defer func() {
		pterm.Success.Writer = oldSuccess
	}()

	// Create command and execute
	cmd := GetManifestCmd()
	config.AddRootFlags(cmd)

	cmd.SetArgs([]string{
		"delete", "flag-1",
		"-m", "flags.json",
	})

	// Execute command
	err = cmd.Execute()
	require.NoError(t, err)

	// Validate the flag was actually deleted from the manifest
	content, err := afero.ReadFile(fs, "flags.json")
	require.NoError(t, err)

	var manifest map[string]any
	err = json.Unmarshal(content, &manifest)
	require.NoError(t, err)

	flags := manifest["flags"].(map[string]any)
	assert.Len(t, flags, 1, "Should have 1 flag remaining")
	assert.Contains(t, flags, "flag-2", "Should still contain flag-2")
	assert.NotContains(t, flags, "flag-1", "Should not contain deleted flag-1")

	// Validate success message is displayed
	output := buf.String()
	assert.Contains(t, output, "flag-1", "Output should contain the flag name")
	assert.Contains(t, output, "deleted successfully", "Output should contain success message")
	assert.Contains(t, output, "flags.json", "Output should contain the manifest path")
}

func TestManifestDeleteCmd_WithCustomManifestPath(t *testing.T) {
	// Setup
	fs := afero.NewMemMapFs()
	filesystem.SetFileSystem(fs)

	// Create manifest in a custom location
	customPath := "custom/path/manifest.json"
	err := fs.MkdirAll("custom/path", 0o755)
	require.NoError(t, err)

	existingManifest := `{
		"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
		"flags": {
			"flag-to-delete": {
				"flagType": "boolean",
				"defaultValue": true,
				"description": "Flag in custom location"
			},
			"flag-to-keep": {
				"flagType": "string",
				"defaultValue": "keep",
				"description": "Keep this one"
			}
		}
	}`
	err = afero.WriteFile(fs, customPath, []byte(existingManifest), 0o644)
	require.NoError(t, err)

	// Create command and execute
	cmd := GetManifestCmd()
	config.AddRootFlags(cmd)

	cmd.SetArgs([]string{
		"delete", "flag-to-delete",
		"-m", customPath,
	})

	// Execute command
	err = cmd.Execute()
	require.NoError(t, err)

	// Validate the flag was deleted from the custom location
	content, err := afero.ReadFile(fs, customPath)
	require.NoError(t, err)

	var manifest map[string]any
	err = json.Unmarshal(content, &manifest)
	require.NoError(t, err)

	flags := manifest["flags"].(map[string]any)
	assert.Len(t, flags, 1, "Should have 1 flag remaining")
	assert.Contains(t, flags, "flag-to-keep")
	assert.NotContains(t, flags, "flag-to-delete")
}

func TestManifestDeleteCmd_ValidatesArgs(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedError string
	}{
		{
			name:          "no flag name provided",
			args:          []string{"delete"},
			expectedError: "accepts 1 arg(s), received 0",
		},
		{
			name:          "too many arguments",
			args:          []string{"delete", "flag1", "flag2"},
			expectedError: "accepts 1 arg(s), received 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			fs := afero.NewMemMapFs()
			filesystem.SetFileSystem(fs)

			// Create command
			cmd := GetManifestCmd()
			config.AddRootFlags(cmd)

			// Set args
			cmd.SetArgs(append(tt.args, "-m", "flags.json"))

			// Execute command - should fail
			err := cmd.Execute()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestManifestDeleteCmd_DeleteFirstFlag(t *testing.T) {
	// Setup
	fs := afero.NewMemMapFs()
	filesystem.SetFileSystem(fs)

	// Create manifest with flags in alphabetical order
	existingManifest := `{
		"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
		"flags": {
			"aaa-first": {
				"flagType": "boolean",
				"defaultValue": true,
				"description": "First flag"
			},
			"bbb-second": {
				"flagType": "string",
				"defaultValue": "test",
				"description": "Second flag"
			},
			"ccc-third": {
				"flagType": "integer",
				"defaultValue": 42,
				"description": "Third flag"
			}
		}
	}`
	err := afero.WriteFile(fs, "flags.json", []byte(existingManifest), 0o644)
	require.NoError(t, err)

	// Create command and execute
	cmd := GetManifestCmd()
	config.AddRootFlags(cmd)
	cmd.SetArgs([]string{"delete", "aaa-first", "-m", "flags.json"})

	// Execute command
	err = cmd.Execute()
	require.NoError(t, err)

	// Validate
	content, err := afero.ReadFile(fs, "flags.json")
	require.NoError(t, err)

	var manifest map[string]any
	err = json.Unmarshal(content, &manifest)
	require.NoError(t, err)

	flags := manifest["flags"].(map[string]any)
	assert.Len(t, flags, 2, "Should have 2 flags remaining")
	assert.NotContains(t, flags, "aaa-first")
	assert.Contains(t, flags, "bbb-second")
	assert.Contains(t, flags, "ccc-third")
}

func TestManifestDeleteCmd_DeleteLastFlag(t *testing.T) {
	// Setup
	fs := afero.NewMemMapFs()
	filesystem.SetFileSystem(fs)

	// Create manifest with flags
	existingManifest := `{
		"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
		"flags": {
			"aaa-first": {
				"flagType": "boolean",
				"defaultValue": true,
				"description": "First flag"
			},
			"bbb-second": {
				"flagType": "string",
				"defaultValue": "test",
				"description": "Second flag"
			},
			"zzz-last": {
				"flagType": "integer",
				"defaultValue": 42,
				"description": "Last flag"
			}
		}
	}`
	err := afero.WriteFile(fs, "flags.json", []byte(existingManifest), 0o644)
	require.NoError(t, err)

	// Create command and execute
	cmd := GetManifestCmd()
	config.AddRootFlags(cmd)
	cmd.SetArgs([]string{"delete", "zzz-last", "-m", "flags.json"})

	// Execute command
	err = cmd.Execute()
	require.NoError(t, err)

	// Validate
	content, err := afero.ReadFile(fs, "flags.json")
	require.NoError(t, err)

	var manifest map[string]any
	err = json.Unmarshal(content, &manifest)
	require.NoError(t, err)

	flags := manifest["flags"].(map[string]any)
	assert.Len(t, flags, 2, "Should have 2 flags remaining")
	assert.Contains(t, flags, "aaa-first")
	assert.Contains(t, flags, "bbb-second")
	assert.NotContains(t, flags, "zzz-last")
}
