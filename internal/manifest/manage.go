package manifest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

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
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return &flagset, nil
}

// Write writes a flagset to a manifest file at the given path
func Write(path string, flagset flagset.Flagset) error {
	flags := make(map[string]any)
	for _, flag := range flagset.Flags {
		flags[flag.Key] = map[string]any{
			"flagType":     flag.Type.String(),
			"description":  flag.Description,
			"defaultValue": flag.DefaultValue,
		}
	}

	m := createInitManifest(flags)
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

// writeManifest marshals and writes a manifest to the given path atomically
func writeManifest(path string, manifest *initManifest) error {
	formattedManifest, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	fs := filesystem.FileSystem()
	dir := filepath.Dir(path)

	// Create temp file in same directory as target
	tmpFile, err := afero.TempFile(fs, dir, "manifest-*.json.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Write to temp file
	if _, err := tmpFile.Write(formattedManifest); err != nil {
		tmpFile.Close()
		_ = fs.Remove(tmpPath) // Ignore cleanup error, prioritize original error
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		_ = fs.Remove(tmpPath) // Ignore cleanup error, prioritize original error
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Atomically rename temp file to target
	if err := fs.Rename(tmpPath, path); err != nil {
		_ = fs.Remove(tmpPath) // Ignore cleanup error, prioritize original error
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
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

// SaveToRemote saves flags to a remote URL using HTTP/HTTPS
func SaveToRemote(url string, flags *flagset.Flagset, authToken string, method string) error {
	// Define the SourceFlag structure for the array-based format
	type SourceFlag struct {
		Key          string `json:"key"`
		Type         string `json:"type"`
		Description  string `json:"description"`
		DefaultValue any    `json:"defaultValue"`
	}

	// Convert flags to source format (array-based)
	var sourceFlags []SourceFlag
	for _, flag := range flags.Flags {
		sourceFlags = append(sourceFlags, SourceFlag{
			Key:          flag.Key,
			Type:         flag.Type.String(),
			Description:  flag.Description,
			DefaultValue: flag.DefaultValue,
		})
	}

	// Wrap in an object with "flags" property
	sourceFlagsFormat := struct {
		Flags []SourceFlag `json:"flags"`
	}{
		Flags: sourceFlags,
	}

	// Marshal to JSON
	data, err := json.Marshal(sourceFlagsFormat)
	if err != nil {
		return fmt.Errorf("failed to marshal flags: %w", err)
	}

	// Create HTTP request with the specified method
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if authToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	}

	// Execute request with a custom client that has a timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for error reporting
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check for successful status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("received error response from destination (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
