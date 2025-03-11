package cmd

import (
	"fmt"

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
			flagSourceUrl := config.GetFlagSourceUrl(cmd)

			manifestExists, _ := filesystem.Exists(manifestPath)
			if manifestExists && !override {
				logger.Default.Debug(fmt.Sprintf("Manifest file already exists at %s", manifestPath))
				confirmMessage := fmt.Sprintf("An existing manifest was found at %s. Would you like to override it?", manifestPath)
				shouldOverride, _ := pterm.DefaultInteractiveConfirm.Show(confirmMessage)
				// Print a blank line for better readability.
				pterm.Println()
				if !shouldOverride {
					logger.Default.Info("No changes were made.")
					return nil
				}

				logger.Default.Debug("User confirmed override of existing manifest")
			}

			logger.Default.Info("Initializing project...")
			err := manifest.Create(manifestPath)
			if err != nil {
				logger.Default.Error(fmt.Sprintf("Failed to create manifest: %v", err))
				return err
			}

			configFileExists, _ := filesystem.Exists(".openfeature.yaml")
			if !configFileExists {
				pterm.Info.Println("Creating .openfeature.yaml configuration file")
				template := getConfigTemplate(flagSourceUrl)
				err = filesystem.WriteFile(".openfeature.yaml", []byte(template))
				if err != nil {
					return err
				}
			} else if flagSourceUrl != "" {
				if !override {
					confirmMessage := "An existing .openfeature.yaml configuration file was found. Would you like to override it?"
					shouldOverride, _ := pterm.DefaultInteractiveConfirm.Show(confirmMessage)
					// Print a blank line for better readability.
					pterm.Println()
					if !shouldOverride {
						logger.Default.Info("Configuration file was not modified.")
					} else {
						pterm.Info.Println("Updating flag source URL in .openfeature.yaml", pterm.LightWhite(flagSourceUrl))
						template := getConfigTemplate(flagSourceUrl)
						err = filesystem.WriteFile(".openfeature.yaml", []byte(template))
						if err != nil {
							return err
						}
					}
				} else {
					pterm.Info.Println("Updating flag source URL in .openfeature.yaml", pterm.LightWhite(flagSourceUrl))
					template := getConfigTemplate(flagSourceUrl)
					err = filesystem.WriteFile(".openfeature.yaml", []byte(template))
					if err != nil {
						return err
					}
				}
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

func getConfigTemplate(flagSourceUrl string) string {
	template := `# OpenFeature CLI Configuration
# This file configures the OpenFeature CLI for your project.
# For full documentation, visit: https://github.com/open-feature/cli#configuration

# Global Configuration
# Path to your flag manifest file (default: "flags.json")
# manifest: "flags.json"

# URL of your flag source for the 'pull' command
# Supports http://, https://, and file:// protocols
`

	if flagSourceUrl != "" {
		template += "flagSourceUrl: " + flagSourceUrl + "\n"
	} else {
		template += "# flagSourceUrl: \"https://your-flag-service.com/api/flags\"\n"
	}

	template += `
# Authentication token for remote flag sources (if required)
# authToken: "your-bearer-token"

# Enable debug logging (default: false)
# debug: false

# Disable interactive prompts (default: false)
# no-input: false

# Command-Specific Configuration
# Override global settings for specific commands

# pull:
#   flag-source-url: "https://api.example.com/flags"
#   auth-token: "pull-specific-token"
#   no-prompt: false

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

	return template
}
