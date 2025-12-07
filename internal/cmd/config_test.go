package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func setupTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "test",
	}

	// Add some test flags
	cmd.Flags().String("output", "", "output path")
	cmd.Flags().String("package-name", "default", "package name")

	return cmd
}

// setupConfigFileForTest creates a temporary directory with a config file
// and changes the working directory to it.
func setupConfigFileForTest(t *testing.T, configContent string) {
	// Create a temporary config file
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, ".openfeature.yaml")
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	// Change to the temporary directory so the config file can be found
	t.Chdir(tmpDir)
}

func TestRootCommandIgnoresUnrelatedConfig(t *testing.T) {
	configContent := `
generate:
  output: output-from-generate
`
	setupConfigFileForTest(t, configContent)

	rootCmd := setupTestCommand()
	err := initializeConfig(rootCmd, "")

	assert.NoError(t, err)
	assert.Equal(t, "", rootCmd.Flag("output").Value.String(),
		"Root command should not get output config from unrelated sections")
}

func TestGenerateCommandGetsGenerateConfig(t *testing.T) {
	configContent := `
generate:
  output: output-from-generate
`
	setupConfigFileForTest(t, configContent)

	generateCmd := setupTestCommand()
	err := initializeConfig(generateCmd, "generate")

	assert.NoError(t, err)
	assert.Equal(t, "output-from-generate", generateCmd.Flag("output").Value.String(),
		"Generate command should get generate.output value")
}

func TestSubcommandGetsSpecificConfig(t *testing.T) {
	configContent := `
generate:
  output: output-from-generate
  go:
    output: output-from-go
    package-name: fromconfig
`
	setupConfigFileForTest(t, configContent)

	goCmd := setupTestCommand()
	err := initializeConfig(goCmd, "generate.go")

	assert.NoError(t, err)
	assert.Equal(t, "output-from-go", goCmd.Flag("output").Value.String(),
		"Go command should get generate.go.output, not generate.output")
	assert.Equal(t, "fromconfig", goCmd.Flag("package-name").Value.String(),
		"Go command should get generate.go.package-name")
}

func TestSubcommandInheritsFromParent(t *testing.T) {
	configContent := `
generate:
  output: output-from-generate
`
	setupConfigFileForTest(t, configContent)

	otherCmd := setupTestCommand()
	err := initializeConfig(otherCmd, "generate.other")

	assert.NoError(t, err)
	assert.Equal(t, "output-from-generate", otherCmd.Flag("output").Value.String(),
		"Other command should inherit generate.output when no specific config exists")
}

func TestCommandLineOverridesConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, ".openfeature.yaml")
	configContent := `
generate:
  output: output-from-config
`
	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	// Change to the temporary directory so the config file can be found
	t.Chdir(tmpDir)

	// Set up a command with a flag value already set via command line
	cmd := setupTestCommand()
	_ = cmd.Flags().Set("output", "output-from-cmdline")

	// Initialize config
	err = initializeConfig(cmd, "generate")
	assert.NoError(t, err)

	// Command line value should take precedence
	assert.Equal(t, "output-from-cmdline", cmd.Flag("output").Value.String(),
		"Command line value should override config file")
}
