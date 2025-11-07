package cmd

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/filesystem"
	"github.com/open-feature/cli/internal/logger"
	"github.com/open-feature/cli/internal/manifest"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func GetInitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new project",
		Long:  "Initialize a new project for OpenFeature CLI.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "init")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := config.GetManifestPath(cmd)
			override := config.GetOverride(cmd)
			providerUrl := config.GetFlagSourceUrl(cmd)

			if err := handleManifestCreation(manifestPath, override); err != nil {
				return err
			}

			if err := handleConfigFile(providerUrl, override); err != nil {
				return err
			}

			pterm.Info.Printfln("Manifest created at %s", pterm.LightWhite(manifestPath))
			logger.Default.FileCreated(manifestPath)
			logger.Default.Success("Project initialized.")
			return nil
		},
	}

	config.AddInitFlags(initCmd)
	addStabilityInfo(initCmd)

	return initCmd
}

func confirmOverride(itemType, path string) (bool, error) {
	message := fmt.Sprintf("An existing %s was found at %s. Would you like to override it?", itemType, path)
	confirmed, err := pterm.DefaultInteractiveConfirm.Show(message)
	if err != nil {
		return false, fmt.Errorf("failed to show confirmation prompt: %w", err)
	}
	pterm.Println() // blank line for readability
	return confirmed, nil
}

func handleManifestCreation(manifestPath string, override bool) error {
	exists, err := filesystem.Exists(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to check if manifest exists: %w", err)
	}

	if exists && !override {
		logger.Default.Debug(fmt.Sprintf("Manifest file already exists at %s", manifestPath))
		shouldOverride, err := confirmOverride("manifest", manifestPath)
		if err != nil {
			return fmt.Errorf("failed to get user confirmation: %w", err)
		}
		if !shouldOverride {
			logger.Default.Info("No changes were made.")
			return nil
		}
		logger.Default.Debug("User confirmed override of existing manifest")
	}

	logger.Default.Info("Initializing project...")
	if err := manifest.Create(manifestPath); err != nil {
		logger.Default.Error(fmt.Sprintf("Failed to create manifest: %v", err))
		return err
	}
	return nil
}

func handleConfigFile(providerURL string, override bool) error {
	configPath := ".openfeature.yaml"
	configExists, err := filesystem.Exists(configPath)
	if err != nil {
		return fmt.Errorf("failed to check if config file exists: %w", err)
	}

	if !configExists {
		return writeConfigFile(providerURL, "Creating .openfeature.yaml configuration file")
	}

	if providerURL == "" {
		return nil // no config to write
	}

	if override {
		return writeConfigFile(providerURL, "Updating provider URL in .openfeature.yaml")
	}

	shouldOverride, err := confirmOverride("configuration file", configPath)
	if err != nil {
		return fmt.Errorf("failed to get user confirmation: %w", err)
	}
	if shouldOverride {
		return writeConfigFile(providerURL, "Updating provider URL in .openfeature.yaml")
	}

	logger.Default.Info("Configuration file was not modified.")
	return nil
}

func writeConfigFile(providerURL, message string) error {
	pterm.Info.Println(message, pterm.LightWhite(providerURL))
	template := getConfigTemplate(providerURL)
	return filesystem.WriteFile(".openfeature.yaml", []byte(template))
}

type configTemplateData struct {
	ProviderURL    string
	HasProviderURL bool
}

const configTemplateText = `# OpenFeature CLI Configuration
# This file configures the OpenFeature CLI for your project.
# For full documentation, visit: https://github.com/open-feature/cli#configuration

# Global Configuration
# Path to your flag manifest file (default: "flags.json")
# manifest: "flags.json"

# URL of your flag provider for the 'pull' and 'push' commands
# Supports http://, https://, and file:// protocols (file:// only for pull)
{{if .HasProviderURL}}provider: {{.ProviderURL}}{{else}}# provider: "https://your-flag-service.com/api/flags"{{end}}

# Authentication token for remote flag providers (if required)
# authToken: "your-bearer-token"

# Enable debug logging (default: false)
# debug: false

# Disable interactive prompts (default: false)
# no-input: false

# Command-Specific Configuration
# Override global settings for specific commands

# pull:
#   provider: "https://api.example.com/flags"
#   auth-token: "pull-specific-token"
#   no-prompt: false

# push:
#   provider: "https://api.example.com/flags"
#   auth-token: "push-specific-token"
#   dry-run: false

# generate:
#   output: "generated"
#
#   # Language-specific generator options
#   go:
#     output: "go/flags"
#     package-name: "openfeature"
#
#   typescript:
#     output: "ts/flags"
#
#   csharp:
#     output: "csharp/flags"
#     namespace: "OpenFeature"
#
#   java:
#     output: "java/flags"
#     package-name: "com.example.openfeature"
`

func getConfigTemplate(providerURL string) string {
	tmpl := template.Must(template.New("config").Parse(configTemplateText))

	data := configTemplateData{
		ProviderURL:    providerURL,
		HasProviderURL: providerURL != "",
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		// Fallback to a simple template if there's an error
		return fmt.Sprintf("# Error generating config template: %v\n", err)
	}

	return buf.String()
}
