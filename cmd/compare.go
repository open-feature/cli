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
	"gopkg.in/yaml.v3"
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
			outputFormat, _ := cmd.Flags().GetString("output")

			// Validate flags
			if sourcePath == "" || targetPath == "" {
				return fmt.Errorf("both source (--manifest) and target (--against) paths are required")
			}

			// Validate output format
			if !manifest.IsValidOutputFormat(outputFormat) {
				return fmt.Errorf("invalid output format: %s. Valid formats are: %s",
					outputFormat, strings.Join(manifest.GetValidOutputFormats(), ", "))
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

			// Render differences based on the output format
			switch manifest.OutputFormat(outputFormat) {
			case manifest.OutputFormatFlat:
				return renderFlatDiff(changes, cmd)
			case manifest.OutputFormatJSON:
				return renderJSONDiff(changes, cmd)
			case manifest.OutputFormatYAML:
				return renderYAMLDiff(changes, cmd)
			default:
				return renderTreeDiff(changes, cmd)
			}
		},
	}

	// Add flags specific to compare command
	compareCmd.Flags().StringP("against", "a", "", "Path to the target manifest file to compare against")
	compareCmd.Flags().StringP("output", "o", string(manifest.OutputFormatTree),
		fmt.Sprintf("Output format. Valid formats: %s", strings.Join(manifest.GetValidOutputFormats(), ", ")))

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

// renderJSONDiff renders changes in JSON format
func renderJSONDiff(changes []manifest.Change, cmd *cobra.Command) error {
	// Create a structured response that can be easily consumed by tools
	type structuredOutput struct {
		TotalChanges   int               `json:"totalChanges" yaml:"totalChanges"`
		Additions      []manifest.Change `json:"additions" yaml:"additions"`
		Removals       []manifest.Change `json:"removals" yaml:"removals"`
		Modifications  []manifest.Change `json:"modifications" yaml:"modifications"`
	}

	// Group changes by type
	var output structuredOutput
	output.TotalChanges = len(changes)

	for _, change := range changes {
		switch change.Type {
		case "add":
			output.Additions = append(output.Additions, change)
		case "remove":
			output.Removals = append(output.Removals, change)
		case "change":
			output.Modifications = append(output.Modifications, change)
		}
	}

	// Convert to JSON
	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON output: %w", err)
	}

	// Print the JSON
	fmt.Println(string(jsonBytes))
	return nil
}

// renderYAMLDiff renders changes in YAML format
func renderYAMLDiff(changes []manifest.Change, cmd *cobra.Command) error {
	// Use the same structured output type as JSON but with YAML tags
	type structuredOutput struct {
		TotalChanges   int               `json:"totalChanges" yaml:"totalChanges"`
		Additions      []manifest.Change `json:"additions" yaml:"additions"`
		Removals       []manifest.Change `json:"removals" yaml:"removals"`
		Modifications  []manifest.Change `json:"modifications" yaml:"modifications"`
	}

	// Group changes by type
	var output structuredOutput
	output.TotalChanges = len(changes)

	for _, change := range changes {
		switch change.Type {
		case "add":
			output.Additions = append(output.Additions, change)
		case "remove":
			output.Removals = append(output.Removals, change)
		case "change":
			output.Modifications = append(output.Modifications, change)
		}
	}

	// Convert to YAML
	yamlBytes, err := yaml.Marshal(output)
	if err != nil {
		return fmt.Errorf("error marshaling YAML output: %w", err)
	}

	// Print the YAML
	fmt.Println(string(yamlBytes))
	return nil
}
