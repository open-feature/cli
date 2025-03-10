package cmd

import (
	"fmt"

	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/filesystem"
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
			if (manifestExists && !override) {
				confirmMessage := fmt.Sprintf("An existing manifest was found at %s. Would you like to override it?", manifestPath)
				shouldOverride, _ := pterm.DefaultInteractiveConfirm.Show(confirmMessage)
				// Print a blank line for better readability.
				pterm.Println()
				if (!shouldOverride) {
					pterm.Info.Println("No changes were made.")
					return nil
				}
			}

			pterm.Info.Println("Initializing project...")
			err := manifest.Create(manifestPath)
			if err != nil {
				return err
			}
			pterm.Info.Printfln("Manifest created at %s", pterm.LightWhite(manifestPath))
			pterm.Success.Println("Project initialized.")
			return nil
		},
	}

	config.AddInitFlags(initCmd)

	return initCmd
}
