package manifest

import (
	"encoding/json"
	"fmt"

	"github.com/open-feature/cli/internal/filesystem"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/spf13/afero"
)

type initManifest struct {
	Schema string `json:"$schema,omitempty"`
	Manifest
}

// Create creates a new manifest file at the given path.
func Create(path string) error {
	m := &initManifest{
		Schema:  "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag_manifest.json",
		Manifest: Manifest{
			Flags: map[string]any{},
		},
	}
	formattedInitManifest, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return filesystem.WriteFile(path, formattedInitManifest)
}

// Loads, validates, and unmarshals the manifest file at the given path into a flagset
func LoadFlagSet(manifestPath string) (*flagset.Flagset, error) {
	fs := filesystem.FileSystem()
	data, err := afero.ReadFile(fs, manifestPath)
	if err != nil {
		return nil, fmt.Errorf("error reading contents from file %q", manifestPath)
	}

	validationErrors, err := Validate(data)
	if err != nil {
		return nil, err
	} else if len(validationErrors) > 0 {
		// TODO tease running manifest validate command
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	var flagset flagset.Flagset
	if err := json.Unmarshal(data, &flagset); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", validationErrors)
	}

	return &flagset, nil
}

func Write(path string, flagset flagset.Flagset) error {
	m := &initManifest{
		Schema:  "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag_manifest.json",
		Manifest: Manifest{
			Flags: map[string]any{},
		},
	}
	for _, flag := range flagset.Flags {
		m.Manifest.Flags[flag.Key] = map[string]any{
			"flagType": flag.Type.String(),
			"description": flag.Description,
			"defaultValue": flag.DefaultValue,
		}
	}
	formattedInitManifest, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return filesystem.WriteFile(path, formattedInitManifest)
}