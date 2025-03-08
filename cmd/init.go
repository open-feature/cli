package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/filesystem"
	"github.com/open-feature/cli/internal/manifest"
	"github.com/spf13/cobra"
)

func GetInitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new project",
		Long:  "Initialize a new project for OpenFeature CLI.",
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := config.GetString(config.ManifestFlag)

			manifestExists, _ := filesystem.Exists(manifestPath)
			fmt.Println("Manifest exists:", manifestExists)

			override, _ := cmd.Flags().GetBool(config.OverrideFlag)
			if !noInput && manifestExists && !override {
				prompt := huh.
					NewConfirm().
					Title("An existing configuration was found. Would you like to override it?").
					Value(&override)
				
				if err := prompt.Run(); err != nil {
					fmt.Println("Error:", err)
					return err
				}

				if !override {
					fmt.Println("Exiting.")
					return nil
				}
			}
			fmt.Println("Initializing project...")


			err := manifest.Create(manifestPath)
			if err != nil {
				return err
			}
			fmt.Println("Project initialized.")
			return nil
		},
	}

	initCmd.Flags().BoolP(config.OverrideFlag, "f", false, "Override existing configuration")

	return initCmd
}
