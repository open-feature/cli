package cmd

import (
	"bytes"
	"testing"

	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/filesystem"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManifestListCmd(t *testing.T) {
	tests := []struct {
		name             string
		manifestContent  string
		expectedError    string
		expectedInOutput []string
		notInOutput      []string
	}{
		{
			name: "list flags in manifest with multiple flags",
			manifestContent: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {
					"feature-a": {
						"flagType": "boolean",
						"defaultValue": true,
						"description": "Feature A flag"
					},
					"max-items": {
						"flagType": "integer",
						"defaultValue": 100,
						"description": "Maximum items allowed"
					},
					"welcome-msg": {
						"flagType": "string",
						"defaultValue": "Hello!",
						"description": "Welcome message"
					}
				}
			}`,
			expectedInOutput: []string{
				"feature-a",
				"max-items",
				"welcome-msg",
				"boolean",
				"integer",
				"string",
				"(3)",
			},
		},
		{
			name: "list flags with empty manifest",
			manifestContent: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {}
			}`,
			expectedInOutput: []string{
				"No flags found in manifest",
			},
			notInOutput: []string{
				"Total flags:",
			},
		},
		{
			name: "list flags with various types",
			manifestContent: `{
				"$schema": "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json",
				"flags": {
					"bool-flag": {
						"flagType": "boolean",
						"defaultValue": false,
						"description": "A boolean flag"
					},
					"string-flag": {
						"flagType": "string",
						"defaultValue": "test",
						"description": "A string flag"
					},
					"int-flag": {
						"flagType": "integer",
						"defaultValue": 42,
						"description": "An integer flag"
					},
					"float-flag": {
						"flagType": "float",
						"defaultValue": 3.14,
						"description": "A float flag"
					},
					"object-flag": {
						"flagType": "object",
						"defaultValue": {"key": "value"},
						"description": "An object flag"
					}
				}
			}`,
			expectedInOutput: []string{
				"bool-flag",
				"string-flag",
				"int-flag",
				"float-flag",
				"object-flag",
				"boolean",
				"string",
				"integer",
				"float",
				"object",
				"(5)",
			},
		},
		{
			name:          "error on missing manifest file",
			expectedError: "failed to load manifest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			fs := afero.NewMemMapFs()
			filesystem.SetFileSystem(fs)

			// Create manifest if provided
			if tt.manifestContent != "" {
				err := afero.WriteFile(fs, "flags.json", []byte(tt.manifestContent), 0o644)
				require.NoError(t, err)
			}

			// Enable pterm output and capture it
			pterm.EnableOutput()
			defer pterm.DisableOutput()

			buf := &bytes.Buffer{}
			oldStdout := pterm.DefaultTable.Writer
			oldSection := pterm.DefaultSection.Writer
			oldInfo := pterm.Info.Writer
			pterm.DefaultTable.Writer = buf
			pterm.DefaultSection.Writer = buf
			pterm.Info.Writer = buf
			defer func() {
				pterm.DefaultTable.Writer = oldStdout
				pterm.DefaultSection.Writer = oldSection
				pterm.Info.Writer = oldInfo
			}()

			// Create command and execute
			cmd := GetManifestCmd()
			config.AddRootFlags(cmd)

			cmd.SetArgs([]string{"list", "-m", "flags.json"})

			// Execute command
			err := cmd.Execute()

			// Validate
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)

				output := buf.String()
				for _, expected := range tt.expectedInOutput {
					assert.Contains(t, output, expected, "Output should contain: %s", expected)
				}
				for _, notExpected := range tt.notInOutput {
					assert.NotContains(t, output, notExpected, "Output should not contain: %s", notExpected)
				}
			}
		})
	}
}

func TestDisplayFlagList(t *testing.T) {
	tests := []struct {
		name             string
		flagset          *flagset.Flagset
		manifestPath     string
		expectedInOutput []string
	}{
		{
			name: "display multiple flags",
			flagset: &flagset.Flagset{
				Flags: []flagset.Flag{
					{
						Key:          "flag1",
						Type:         flagset.BoolType,
						Description:  "First flag",
						DefaultValue: true,
					},
					{
						Key:          "flag2",
						Type:         flagset.StringType,
						Description:  "Second flag",
						DefaultValue: "test",
					},
				},
			},
			manifestPath: "test.json",
			expectedInOutput: []string{
				"flag1",
				"flag2",
				"boolean",
				"string",
				"First flag",
				"Second flag",
				"test.json",
			},
		},
		{
			name: "display empty flagset",
			flagset: &flagset.Flagset{
				Flags: []flagset.Flag{},
			},
			manifestPath: "empty.json",
			expectedInOutput: []string{
				"No flags found in manifest",
			},
		},
		{
			name: "truncate long description",
			flagset: &flagset.Flagset{
				Flags: []flagset.Flag{
					{
						Key:          "long-desc",
						Type:         flagset.BoolType,
						Description:  "This is a very long description that should be truncated because it exceeds the maximum length",
						DefaultValue: false,
					},
				},
			},
			manifestPath: "test.json",
			expectedInOutput: []string{
				"long-desc",
				"...",
			},
		},
		{
			name: "truncate long string value",
			flagset: &flagset.Flagset{
				Flags: []flagset.Flag{
					{
						Key:          "long-string",
						Type:         flagset.StringType,
						Description:  "Long string value",
						DefaultValue: "This is a very long string value that should be truncated",
					},
				},
			},
			manifestPath: "test.json",
			expectedInOutput: []string{
				"long-string",
				"...",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Enable pterm output and capture it
			pterm.EnableOutput()
			defer pterm.DisableOutput()

			buf := &bytes.Buffer{}
			oldStdout := pterm.DefaultTable.Writer
			oldSection := pterm.DefaultSection.Writer
			oldInfo := pterm.Info.Writer
			pterm.DefaultTable.Writer = buf
			pterm.DefaultSection.Writer = buf
			pterm.Info.Writer = buf
			defer func() {
				pterm.DefaultTable.Writer = oldStdout
				pterm.DefaultSection.Writer = oldSection
				pterm.Info.Writer = oldInfo
			}()

			// Call the function
			displayFlagList(tt.flagset, tt.manifestPath)

			// Validate output
			output := buf.String()
			for _, expected := range tt.expectedInOutput {
				assert.Contains(t, output, expected, "Output should contain: %s", expected)
			}
		})
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected string
	}{
		{
			name:     "format short string",
			value:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "format long string",
			value:    "this is a very long string that exceeds thirty characters",
			expected: "...",
		},
		{
			name:     "format boolean true",
			value:    true,
			expected: "true",
		},
		{
			name:     "format boolean false",
			value:    false,
			expected: "false",
		},
		{
			name:     "format integer",
			value:    42,
			expected: "42",
		},
		{
			name:     "format float",
			value:    3.14,
			expected: "3.14",
		},
		{
			name:     "format object",
			value:    map[string]any{"key": "value"},
			expected: `{"key":"value"}`,
		},
		{
			name:     "format large object",
			value:    map[string]any{"key1": "value1", "key2": "value2", "key3": "value3"},
			expected: "...",
		},
		{
			name:     "format array",
			value:    []any{"a", "b", "c"},
			expected: `["a","b","c"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatValue(tt.value)
			assert.Contains(t, result, tt.expected)
		})
	}
}

func TestMain(m *testing.M) {
	// Disable pterm output during tests by default
	pterm.DisableOutput()
	m.Run()
}
