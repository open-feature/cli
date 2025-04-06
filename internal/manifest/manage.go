package manifest

import (
	"encoding/json"

	"github.com/open-feature/cli/internal/filesystem"
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
