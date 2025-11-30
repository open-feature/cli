package cmd

import (
	"fmt"

	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/filesystem"
	"github.com/open-feature/cli/internal/logger"
	"github.com/open-feature/cli/internal/manifest"
	"github.com/pterm/pterm"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func GetManifestDeleteCmd() *cobra.Command {
	manifestDeleteCmd := &cobra.Command{
		Use:   "delete <flag-name>",
		Short: "Delete a flag from the manifest",
		Long: `Delete a flag from the manifest file by its key.

Examples:
  # Delete a flag named 'old-feature'
  openfeature manifest delete old-feature

  # Delete a flag from a specific manifest file
  openfeature manifest delete old-feature --manifest path/to/flags.json`,
		Args: cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "manifest.delete")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			flagName := args[0]
			manifestPath := config.GetManifestPath(cmd)

			// Check if manifest exists
			exists, err := afero.Exists(filesystem.FileSystem(), manifestPath)
			if err != nil {
				return fmt.Errorf("failed to check manifest existence: %w", err)
			}

			if !exists {
				return fmt.Errorf("manifest file does not exist: %s", manifestPath)
			}

			// Load existing manifest
			fs, err := manifest.LoadFlagSet(manifestPath)
			if err != nil {
				return fmt.Errorf("failed to load manifest: %w", err)
			}

			// Check if manifest has any flags
			if len(fs.Flags) == 0 {
				return fmt.Errorf("manifest contains no flags")
			}

			// Check if flag exists
			flagIndex := -1
			for i, flag := range fs.Flags {
				if flag.Key == flagName {
					flagIndex = i
					break
				}
			}

			if flagIndex == -1 {
				return fmt.Errorf("flag '%s' not found in manifest", flagName)
			}

			// Remove the flag by creating a new slice without it
			fs.Flags = append(fs.Flags[:flagIndex], fs.Flags[flagIndex+1:]...)

			// Write updated manifest
			if err := manifest.Write(manifestPath, *fs); err != nil {
				return fmt.Errorf("failed to write manifest: %w", err)
			}

			// Success message
			pterm.Success.Printfln("Flag '%s' deleted successfully from %s", flagName, manifestPath)
			logger.Default.Debug(fmt.Sprintf("Deleted flag: name=%s, manifestPath=%s", flagName, manifestPath))

			return nil
		},
	}

	// Add command-specific flags
	config.AddManifestDeleteFlags(manifestDeleteCmd)
	addStabilityInfo(manifestDeleteCmd)

	return manifestDeleteCmd
}
