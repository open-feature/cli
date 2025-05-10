package cmd

import (
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

	// Verify optional flags
	flatFlag := cmd.Flag("flat")
	assert.NotNil(t, flatFlag)
	assert.Equal(t, "false", flatFlag.DefValue)
}

func TestCompareManifests(t *testing.T) {
	// This test mainly verifies the command executes without errors
	// since stdout capture is difficult with pterm usage

	// Need to use the root command to properly inherit the manifest flag
	rootCmd := GetRootCmd()

	// Setup command line arguments
	rootCmd.SetArgs([]string{
		"compare",
		"--manifest", "testdata/source_manifest.json",
		"--against", "testdata/target_manifest.json",
		"--flat",
	})

	// Execute command
	err := rootCmd.Execute()
	assert.NoError(t, err, "Command should execute without errors")
}
