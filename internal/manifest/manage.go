package manifest

import (
	"encoding/json"

	"github.com/open-feature/cli/internal/filesystem"
	"github.com/spf13/afero"
)

type initManifest struct {
	Schema string `json:"$schema,omitempty"`
	Manifest
}

// Create creates a new manifest file at the given path.
func Create(path string) error {
	m := &initManifest{
		Schema: "https://raw.githubusercontent.com/open-feature/cli/main/schema/v0/flag-manifest.json",
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

// Load loads a manifest from a JSON file, unmarshals it, and returns a Manifest object.
func Load(path string) (*Manifest, error) {
	fs := filesystem.FileSystem()
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, err
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}
