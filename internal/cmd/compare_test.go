package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/open-feature/cli/internal/manifest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCompareCmd(t *testing.T) {
	cmd := GetCompareCmd()

	assert.Equal(t, "compare", cmd.Use)
	assert.Equal(t, "Compare two feature flag manifests", cmd.Short)

	// Verify flags exist
	againstFlag := cmd.Flag("against")
	assert.NotNil(t, againstFlag)

	// Verify output flag
	outputFlag := cmd.Flag("output")
	assert.NotNil(t, outputFlag)
	assert.Equal(t, "tree", outputFlag.DefValue)

	// Verify ignore flag
	ignoreFlag := cmd.Flag("ignore")
	assert.NotNil(t, ignoreFlag)

	// Verify reverse flag
	reverseFlag := cmd.Flag("reverse")
	assert.NotNil(t, reverseFlag)
	assert.Equal(t, "false", reverseFlag.DefValue)
}

func TestCompareManifests(t *testing.T) {
	// This test mainly verifies the command executes without errors
	// with each of the supported output formats

	formats := []string{"tree", "flat", "json", "yaml"}

	for _, format := range formats {
		t.Run(fmt.Sprintf("output_format_%s", format), func(t *testing.T) {
			// Need to use the root command to properly inherit the manifest flag
			rootCmd := GetRootCmd()

			// Setup command line arguments
			rootCmd.SetArgs([]string{
				"compare",
				"--manifest", "testdata/source_manifest.json",
				"--against", "testdata/target_manifest.json",
				"--output", format,
			})

			// Execute command
			err := rootCmd.Execute()
			assert.NoError(t, err, "Command should execute without errors with output format: "+format)
		})
	}
}

// captureStdout captures stdout during test execution
func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "captureStdout: error copying output: %v\n", err)
	}
	return buf.String()
}

func TestCompareWithIgnoreFlag(t *testing.T) {
	// Test that the ignore flag is properly parsed and passed to comparison

	t.Run("single_ignore_pattern", func(t *testing.T) {
		output := captureStdout(func() {
			rootCmd := GetRootCmd()

			rootCmd.SetArgs([]string{
				"compare",
				"--manifest", "testdata/source_manifest.json",
				"--against", "testdata/target_manifest.json",
				"--ignore", "description",
			})

			err := rootCmd.Execute()
			assert.NoError(t, err, "Command should execute with single ignore pattern")
		})

		// Verify that the output doesn't contain description changes in the field-level diff
		// The word "description" may still appear in JSON for additions/removals, which is fine
		assert.NotContains(t, output, "• description:",
			"Output should not show description field changes when it's ignored")
	})

	t.Run("multiple_ignore_patterns", func(t *testing.T) {
		output := captureStdout(func() {
			rootCmd := GetRootCmd()

			rootCmd.SetArgs([]string{
				"compare",
				"--manifest", "testdata/source_manifest.json",
				"--against", "testdata/target_manifest.json",
				"--ignore", "description",
				"--ignore", "metadata.*",
			})

			err := rootCmd.Execute()
			assert.NoError(t, err, "Command should execute with multiple ignore patterns")
		})

		// Verify that the output doesn't contain ignored fields in the field-level diff
		assert.NotContains(t, output, "• description:",
			"Output should not show description field changes when it's ignored")
		assert.NotContains(t, output, "• metadata",
			"Output should not show metadata field changes when it's ignored")
	})

	t.Run("ignore_with_wildcard", func(t *testing.T) {
		output := captureStdout(func() {
			rootCmd := GetRootCmd()

			rootCmd.SetArgs([]string{
				"compare",
				"--manifest", "testdata/source_manifest.json",
				"--against", "testdata/target_manifest.json",
				"--ignore", "flags.*.description",
			})

			err := rootCmd.Execute()
			assert.NoError(t, err, "Command should execute with wildcard ignore pattern")
		})

		// Verify that the output doesn't contain description changes in the field-level diff
		assert.NotContains(t, output, "• description:",
			"Output should not show description field changes when using wildcard pattern 'flags.*.description'")
	})

	t.Run("without_ignore_shows_description", func(t *testing.T) {
		output := captureStdout(func() {
			rootCmd := GetRootCmd()

			rootCmd.SetArgs([]string{
				"compare",
				"--manifest", "testdata/source_manifest.json",
				"--against", "testdata/target_manifest.json",
			})

			err := rootCmd.Execute()
			assert.NoError(t, err, "Command should execute without ignore pattern")
		})

		// Verify that the output DOES contain description field changes when not ignored
		assert.Contains(t, output, "• description:",
			"Output should show description field changes when it's not ignored")
	})
}

func TestCompareWithReverseFlag(t *testing.T) {
	// Test that the --reverse flag properly reverses the comparison direction

	type compareResult struct {
		TotalChanges  int               `json:"totalChanges"`
		Additions     []manifest.Change `json:"additions"`
		Removals      []manifest.Change `json:"removals"`
		Modifications []manifest.Change `json:"modifications"`
	}

	// Helper function to unmarshal compare JSON output
	unmarshalCompareOutput := func(output string) compareResult {
		var result compareResult
		err := json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Should be valid JSON output")
		return result
	}

	// Helper function to find a change by path in a slice
	hasChangePath := func(changes []manifest.Change, path string) bool {
		for _, change := range changes {
			if change.Path == path {
				return true
			}
		}
		return false
	}

	t.Run("without_reverse_flag", func(t *testing.T) {
		output := captureStdout(func() {
			rootCmd := GetRootCmd()

			rootCmd.SetArgs([]string{
				"compare",
				"--manifest", "testdata/source_manifest.json",
				"--against", "testdata/target_manifest.json",
				"--output", "json",
			})

			err := rootCmd.Execute()
			assert.NoError(t, err, "Command should execute without reverse flag")
		})

		result := unmarshalCompareOutput(output)

		// Without --reverse: Compare(target, source) shows what HAS changed
		// source has maxItems, target doesn't → should show as addition
		// target has welcomeMessage, source doesn't → should show as removal
		assert.True(t, hasChangePath(result.Additions, "flags.maxItems"),
			"maxItems should be in additions (exists in source, not in target)")
		assert.True(t, hasChangePath(result.Removals, "flags.welcomeMessage"),
			"welcomeMessage should be in removals (exists in target, not in source)")

		// Verify it's NOT in the wrong sections
		assert.False(t, hasChangePath(result.Removals, "flags.maxItems"),
			"maxItems should NOT be in removals")
		assert.False(t, hasChangePath(result.Additions, "flags.welcomeMessage"),
			"welcomeMessage should NOT be in additions")
	})

	t.Run("with_reverse_flag", func(t *testing.T) {
		output := captureStdout(func() {
			rootCmd := GetRootCmd()

			rootCmd.SetArgs([]string{
				"compare",
				"--manifest", "testdata/source_manifest.json",
				"--against", "testdata/target_manifest.json",
				"--output", "json",
				"--reverse",
			})

			err := rootCmd.Execute()
			assert.NoError(t, err, "Command should execute with reverse flag")
		})

		result := unmarshalCompareOutput(output)

		// With --reverse: Compare(source, target) shows what WILL change
		// source has maxItems, target doesn't → should show as removal
		// target has welcomeMessage, source doesn't → should show as addition
		assert.True(t, hasChangePath(result.Removals, "flags.maxItems"),
			"maxItems should be in removals with --reverse (exists in source, not in target)")
		assert.True(t, hasChangePath(result.Additions, "flags.welcomeMessage"),
			"welcomeMessage should be in additions with --reverse (exists in target, not in source)")

		// Verify it's NOT in the wrong sections
		assert.False(t, hasChangePath(result.Additions, "flags.maxItems"),
			"maxItems should NOT be in additions")
		assert.False(t, hasChangePath(result.Removals, "flags.welcomeMessage"),
			"welcomeMessage should NOT be in removals")
	})
}
