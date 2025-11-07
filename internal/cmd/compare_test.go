package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

		// Without --reverse: Compare(target, source) shows what HAS changed
		// source has maxItems, target doesn't → should show as addition
		// target has welcomeMessage, source doesn't → should show as removal
		assert.Contains(t, output, `"path": "flags.maxItems"`,
			"Should include maxItems in changes")
		assert.Contains(t, output, `"path": "flags.welcomeMessage"`,
			"Should include welcomeMessage in changes")
		assert.Contains(t, output, `"additions"`,
			"Should have additions section")
		assert.Contains(t, output, `"removals"`,
			"Should have removals section")

		// Verify maxItems is in additions and welcomeMessage is in removals
		// by checking that the path appears in the right section
		additionsStart := strings.Index(output, `"additions"`)
		removalsStart := strings.Index(output, `"removals"`)
		maxItemsPos := strings.Index(output, `"path": "flags.maxItems"`)
		welcomeMessagePos := strings.Index(output, `"path": "flags.welcomeMessage"`)

		assert.True(t, maxItemsPos > additionsStart && maxItemsPos < removalsStart,
			"maxItems should be in additions section")
		assert.True(t, welcomeMessagePos > removalsStart,
			"welcomeMessage should be in removals section")
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

		// With --reverse: Compare(source, target) shows what WILL change
		// source has maxItems, target doesn't → should show as removal
		// target has welcomeMessage, source doesn't → should show as addition
		assert.Contains(t, output, `"path": "flags.maxItems"`,
			"Should include maxItems in changes")
		assert.Contains(t, output, `"path": "flags.welcomeMessage"`,
			"Should include welcomeMessage in changes")

		// Verify maxItems is in removals and welcomeMessage is in additions
		// by checking that the path appears in the right section
		additionsStart := strings.Index(output, `"additions"`)
		removalsStart := strings.Index(output, `"removals"`)
		maxItemsPos := strings.Index(output, `"path": "flags.maxItems"`)
		welcomeMessagePos := strings.Index(output, `"path": "flags.welcomeMessage"`)

		assert.True(t, maxItemsPos > removalsStart,
			"maxItems should be in removals section with --reverse")
		assert.True(t, welcomeMessagePos > additionsStart && welcomeMessagePos < removalsStart,
			"welcomeMessage should be in additions section with --reverse")
	})
}
