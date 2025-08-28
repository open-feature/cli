package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/open-feature/cli/internal/filesystem"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/spf13/afero"
)

const flagManifestSchemaURL = "https://raw.githubusercontent.com/open-feature/cli/refs/heads/main/schema/v0/flag-manifest.json"

type initManifest struct {
	Schema string `json:"$schema,omitempty"`
	Manifest
}

// Create creates a new manifest file at the given path.
func Create(path string) error {
	m := createInitManifest(map[string]any{})
	return writeManifest(path, m)
}

// LoadFlagSet loads, validates, and unmarshals the manifest file at the given path into a flagset
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
		return nil, errors.New(FormatValidationError(validationErrors))
	}

	var flagset flagset.Flagset
	if err := json.Unmarshal(data, &flagset); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", validationErrors)
	}

	return &flagset, nil
}

// Write writes a flagset to a manifest file at the given path
func Write(path string, flagset flagset.Flagset) error {
	// Marshal the flagset to get the flags structure
	flagsetJSON, err := json.Marshal(flagset)
	if err != nil {
		return err
	}

	// Unmarshal into a map to extract the flags
	var flagsetMap map[string]any
	if err := json.Unmarshal(flagsetJSON, &flagsetMap); err != nil {
		return err
	}

	m := createInitManifest(flagsetMap["flags"].(map[string]any))
	return writeManifest(path, m)
}

// LoadFromLocal loads flags from a local file path
func LoadFromLocal(filePath string) (*flagset.Flagset, error) {
	fs := filesystem.FileSystem()
	data, err := afero.ReadFile(fs, filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading local flags file: %w", err)
	}

	flags, err := loadFlagsFromData(data)
	if err != nil {
		return nil, fmt.Errorf("error loading flags from local file: %w", err)
	}

	return flags, nil
}

// LoadFromRemote loads flags from a remote URL
func LoadFromRemote(url string, authToken string) (*flagset.Flagset, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if authToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Received error response from flag source: %s", string(body))
	}

	return loadFlagsFromData(body)
}

// createInitManifest creates an initManifest with the given flags
func createInitManifest(flags map[string]any) *initManifest {
	return &initManifest{
		Schema: flagManifestSchemaURL,
		Manifest: Manifest{
			Flags: flags,
		},
	}
}

// writeManifest marshals and writes a manifest to the given path
func writeManifest(path string, manifest *initManifest) error {
	formattedManifest, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return filesystem.WriteFile(path, formattedManifest)
}

// loadFlagsFromData attempts to load flags from JSON data using multiple formats
func loadFlagsFromData(data []byte) (*flagset.Flagset, error) {
	// Try the standard manifest format first (with flags as object keys)
	var flags flagset.Flagset
	if err := json.Unmarshal(data, &flags); err == nil {
		return &flags, nil
	}

	// Fallback to source flags format (array-based)
	loadedFlags, err := flagset.LoadFromSourceFlags(data)
	if err != nil {
		return nil, err
	}

	return &flagset.Flagset{Flags: *loadedFlags}, nil
}
