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
	sourceFlag := cmd.Flag("source")
	assert.NotNil(t, sourceFlag)

	targetFlag := cmd.Flag("target")
	assert.NotNil(t, targetFlag)

	// Verify optional flags
	flatFlag := cmd.Flag("flat")
	assert.NotNil(t, flatFlag)
	assert.Equal(t, "false", flatFlag.DefValue)
}

func TestCompareManifests(t *testing.T) {
	// This test mainly verifies the command executes without errors
	// since stdout capture is difficult with pterm usage

	cmd := GetCompareCmd()

	// Setup command line arguments
	cmd.SetArgs([]string{
		"--source", "testdata/source_manifest.json",
		"--target", "testdata/target_manifest.json",
		"--flat",
	})

	// Execute command
	err := cmd.Execute()
	assert.NoError(t, err, "Command should execute without errors")
}
