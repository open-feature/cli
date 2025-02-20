package cmd

import (
	"encoding/json"
	"os"

	"github.com/open-feature/cli/internal/manifest"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Manage OpenFeature manifests",
	Long:  `Manage OpenFeature manifests with subcommands to init, validate, and compare manifests`,
}

type InitManifest struct {
	Schema string `json:"$schema,omitempty"`
	manifest.Manifest
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new manifest",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(ManifestPath); err == nil {
			overwrite, _ := cmd.Flags().GetBool("overwrite")
			if !overwrite {
				pterm.Warning.Printf("%s already exists. Use --overwrite to overwrite", ManifestPath)
				return nil
			}
		}

		m := &InitManifest{
			Schema: "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag_manifest.json",
			Manifest: manifest.Manifest{
				Flags: map[string]any{},
			},
		}
		formattedInitManifest, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			return err
		}
		err = os.WriteFile(ManifestPath, formattedInitManifest, 0644)
		if err != nil {
			pterm.Error.Println("error creating manifest:", err)
			return nil
		}
		pterm.Success.Printf("%s created successfully\n", ManifestPath)
		return nil
	},
}