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

			logger.Default.FileCreated(manifestPath)
			logger.Default.Success("Project initialized.")
			return nil
		},
	}

	config.AddInitFlags(initCmd)

	addStabilityInfo(initCmd)

	return initCmd
}
