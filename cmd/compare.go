package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/manifest"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func GetCompareCmd() *cobra.Command {
	compareCmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare two feature flag manifests",
		Long:  "Compare two OpenFeature flag manifests and display the differences in a structured format.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "compare")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get flags
			sourcePath := config.GetManifestPath(cmd)
			targetPath, _ := cmd.Flags().GetString("against")
			flatOutput, _ := cmd.Flags().GetBool("flat")

			// Validate flags
			if sourcePath == "" || targetPath == "" {
				return fmt.Errorf("both source (--manifest) and target (--against) paths are required")
			}

			// Load manifests
			sourceManifest, err := loadManifest(sourcePath)
			if err != nil {
				return fmt.Errorf("error loading source manifest: %w", err)
			}

			targetManifest, err := loadManifest(targetPath)
			if err != nil {
				return fmt.Errorf("error loading target manifest: %w", err)
			}

			// Compare manifests
			changes, err := manifest.Compare(sourceManifest, targetManifest)
			if err != nil {
				return fmt.Errorf("error comparing manifests: %w", err)
			}

			// No changes
			if len(changes) == 0 {
				pterm.Success.Println("No differences found between the manifests.")
				return nil
			}

			// Render differences based on the output mode
			if flatOutput {
				return renderFlatDiff(changes, cmd)
			}
			return renderTreeDiff(changes, cmd)
		},
	}

	// Add flags specific to compare command
	compareCmd.Flags().StringP("against", "a", "", "Path to the target manifest file to compare against")
	compareCmd.Flags().Bool("flat", false, "Display differences in a flat format")

	// Mark required flags
	_ = compareCmd.MarkFlagRequired("against")

	return compareCmd
}

// loadManifest loads and unmarshals a manifest file from the given path
func loadManifest(path string) (*manifest.Manifest, error) {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal JSON
	var m manifest.Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return &m, nil
}

// renderTreeDiff renders changes with tree-structured inline differences
func renderTreeDiff(changes []manifest.Change, cmd *cobra.Command) error {
	pterm.Info.Printf("Found %d difference(s) between manifests:\n\n", len(changes))

	// Group changes by type for easier reading
	var (
		additions     []manifest.Change
		removals      []manifest.Change
		modifications []manifest.Change
	)

	for _, change := range changes {
		switch change.Type {
		case "add":
			additions = append(additions, change)
		case "remove":
			removals = append(removals, change)
		case "change":
			modifications = append(modifications, change)
		}
	}

	// Print additions
	if len(additions) > 0 {
		pterm.FgGreen.Println("◆ Additions:")
		for _, change := range additions {
			flagName := strings.TrimPrefix(change.Path, "flags.")
			pterm.FgGreen.Printf("  + %s\n", flagName)
			valueJSON, _ := json.MarshalIndent(change.NewValue, "    ", "  ")
			fmt.Printf("    %s\n", valueJSON)
		}
		fmt.Println()
	}

	// Print removals
	if len(removals) > 0 {
		pterm.FgRed.Println("◆ Removals:")
		for _, change := range removals {
			flagName := strings.TrimPrefix(change.Path, "flags.")
			pterm.FgRed.Printf("  - %s\n", flagName)
			valueJSON, _ := json.MarshalIndent(change.OldValue, "    ", "  ")
			fmt.Printf("    %s\n", valueJSON)
		}
		fmt.Println()
	}

	// Print modifications
	if len(modifications) > 0 {
		pterm.FgYellow.Println("◆ Modifications:")
		for _, change := range modifications {
			flagName := strings.TrimPrefix(change.Path, "flags.")
			pterm.FgYellow.Printf("  ~ %s\n", flagName)

			// Marshall the values
			oldJSON, _ := json.MarshalIndent(change.OldValue, "", "  ")
			newJSON, _ := json.MarshalIndent(change.NewValue, "", "  ")

			// Print the diff
			fmt.Println("    Before:")
			for _, line := range strings.Split(string(oldJSON), "\n") {
				fmt.Printf("      %s\n", line)
			}

			fmt.Println("    After:")
			for _, line := range strings.Split(string(newJSON), "\n") {
				fmt.Printf("      %s\n", line)
			}
		}
	}

	return nil
}

// renderFlatDiff renders changes in a flat format
func renderFlatDiff(changes []manifest.Change, cmd *cobra.Command) error {
	pterm.Info.Printf("Found %d difference(s) between manifests:\n\n", len(changes))

	for _, change := range changes {
		flagName := strings.TrimPrefix(change.Path, "flags.")
		switch change.Type {
		case "add":
			pterm.FgGreen.Printf("+ %s\n", flagName)
		case "remove":
			pterm.FgRed.Printf("- %s\n", flagName)
		case "change":
			pterm.FgYellow.Printf("~ %s\n", flagName)
		}
	}

	return nil
}
