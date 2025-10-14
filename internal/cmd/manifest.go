package cmd

import (
	"github.com/spf13/cobra"
)

func GetManifestCmd() *cobra.Command {
	manifestCmd := &cobra.Command{
		Use:   "manifest",
		Short: "Manage flag manifest files",
		Long:  `Commands for managing OpenFeature flag manifest files.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceErrors:              true,
		SilenceUsage:               true,
		DisableSuggestions:         false,
		SuggestionsMinimumDistance: 2,
	}

	// Add subcommands
	manifestCmd.AddCommand(GetManifestAddCmd())
	manifestCmd.AddCommand(GetManifestListCmd())

	addStabilityInfo(manifestCmd)

	return manifestCmd
}
