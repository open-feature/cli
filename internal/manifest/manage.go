package manifest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/open-feature/cli/internal/api/sync"
	"github.com/open-feature/cli/internal/filesystem"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/logger"
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

// LoadFromSyncAPI loads flags from a remote URL using the sync API client
// This should be used when the remote source implements the sync API specification
func LoadFromSyncAPI(baseURL string, authToken string) (*flagset.Flagset, error) {
	logger.Default.Debug(fmt.Sprintf("Loading flags from sync API at %s", baseURL))

	client, err := sync.NewClient(baseURL, authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create sync client: %w", err)
	}

	ctx := context.Background()
	return client.PullFlags(ctx)
}

// LoadFromRemote loads flags from a remote URL using direct HTTP requests
// This is a fallback for sources that don't implement the sync API specification
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

	// Debug: log what we got from the server
	logger.Default.Debug(fmt.Sprintf("Fetched from %s (status %d):\n%s", url, resp.StatusCode, string(body)))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("received error response from flag source: %s", string(body))
	}

	return loadFlagsFromData(body)
}

// URLLooksLikeAFile checks if the given URL string appears to point to a file
func URLLooksLikeAFile(url string) bool {
	fileExtensions := []string{".json", ".yaml", ".yml"}
	for _, ext := range fileExtensions {
		if len(url) >= len(ext) && url[len(url)-len(ext):] == ext {
			return true
		}
	}
	return false
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
// This function performs a smart push: it fetches remote flags first,
// compares them with local flags, and intelligently creates or updates
// flags as needed. Returns a PushResult with details of what was changed.
// If dryRun is true, only performs the comparison without making actual API calls.
func SaveToRemote(url string, flags *flagset.Flagset, authToken string, dryRun bool) (*sync.PushResult, error) {
	// Use the generated OpenAPI client for type-safe API calls
	client, err := sync.NewClient(url, authToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create push client: %w", err)
	}

	ctx := context.Background()

	// Fetch remote flags to compare with local flags using the sync client
	logger.Default.Debug("Fetching remote flags for comparison")
	remoteFlags, err := client.PullFlags(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch remote flags: %w", err)
	}
	logger.Default.Debug(fmt.Sprintf("Fetched %d remote flags", len(remoteFlags.Flags)))

	// Smart push: compare and intelligently create or update flags
	return client.PushFlags(ctx, flags, remoteFlags, flagManifestSchemaURL, dryRun)
}
