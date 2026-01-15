package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/manifest"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func GetManifestListCmd() *cobra.Command {
	manifestListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all flags in the manifest",
		Long:  `Display all flags defined in the manifest file with their configuration.`,
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "manifest.list")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := config.GetManifestPath(cmd)

			// Load existing manifest
			fs, err := manifest.LoadFlagSet(manifestPath)
			if err != nil {
				return fmt.Errorf("failed to load manifest: %w", err)
			}

			displayFlagList(fs, manifestPath)
			return nil
		},
	}

	// Add command-specific flags
	config.AddManifestListFlags(manifestListCmd)
	addStabilityInfo(manifestListCmd)

	return manifestListCmd
}

// displayFlagList prints a formatted table of all flags in the flagset
func displayFlagList(fs *flagset.Flagset, manifestPath string) {
	if len(fs.Flags) == 0 {
		pterm.Info.Println("No flags found in manifest")
		return
	}

	// Print header
	pterm.DefaultSection.Println(fmt.Sprintf("Flags in %s (%d)", manifestPath, len(fs.Flags)))

	// Create table data
	tableData := pterm.TableData{
		{"Key", "Type", "Default Value", "Description"},
	}

	for _, flag := range fs.Flags {
		// Format default value for display
		defaultValueStr := formatValue(flag.DefaultValue)

		// Truncate description if too long
		description := flag.Description
		const maxDescriptionLength = 50

		if len(description) > maxDescriptionLength {
			description = description[:maxDescriptionLength-3] + "..."
		}

		tableData = append(tableData, []string{
			flag.Key,
			flag.Type.String(),
			defaultValueStr,
			description,
		})
	}

	// Render table
	_ = pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

// formatValue converts a value to a string representation suitable for display
func formatValue(value any) string {
	switch v := value.(type) {
	case string:
		if len(v) > 30 {
			return fmt.Sprintf("\"%s...\"", v[:27])
		}
		return fmt.Sprintf("\"%s\"", v)
	case bool, int, float64:
		return fmt.Sprintf("%v", v)
	case map[string]any, []any:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		jsonStr := string(jsonBytes)
		if len(jsonStr) > 30 {
			return jsonStr[:27] + "..."
		}
		return jsonStr
	default:
		return fmt.Sprintf("%v", v)
	}
}
