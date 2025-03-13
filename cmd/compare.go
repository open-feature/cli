package cmd

import (
	"fmt"

	"github.com/open-feature/cli/internal/manifest"
	"github.com/spf13/cobra"
)

func GetCompareCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "compare",
		Short: "Compare two manifest files",
		Long:  `Compare two manifest files and list the changes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("please provide two manifest files to compare")
			}

			oldManifest, err := manifest.Load(args[0])
			if err != nil {
				return fmt.Errorf("failed to load old manifest: %w", err)
			}

			newManifest, err := manifest.Load(args[1])
			if err != nil {
				return fmt.Errorf("failed to load new manifest: %w", err)
			}

			changes, err := manifest.Compare(oldManifest, newManifest)
			if err != nil {
				return fmt.Errorf("failed to compare manifests: %w", err)
			}

			for _, change := range changes {
				fmt.Printf("%s: %s\n", change.Type, change.Path)
			}

			return nil
		},
	}
}
