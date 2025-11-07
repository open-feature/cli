package cmd

import (
	"fmt"
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

func TestCompareWithIgnoreFlag(t *testing.T) {
	// Test that the ignore flag is properly parsed and passed to comparison

	t.Run("single_ignore_pattern", func(t *testing.T) {
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

	t.Run("multiple_ignore_patterns", func(t *testing.T) {
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

	t.Run("ignore_with_wildcard", func(t *testing.T) {
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
}
