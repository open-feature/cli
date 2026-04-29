package manifest

import (
	"testing"

	"github.com/open-feature/cli/internal/filesystem"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestURLLooksLikeAFile(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "URL with .json extension",
			url:      "https://example.com/flags.json",
			expected: true,
		},
		{
			name:     "URL with .yaml extension",
			url:      "https://example.com/flags.yaml",
			expected: true,
		},
		{
			name:     "URL with .yml extension",
			url:      "https://example.com/flags.yml",
			expected: true,
		},
		{
			name:     "URL with path and .json extension",
			url:      "https://example.com/api/v0/flags.json",
			expected: true,
		},
		{
			name:     "URL with query params and .json extension",
			url:      "https://example.com/flags.json?version=1",
			expected: false, // Query params come after extension
		},
		{
			name:     "URL without file extension",
			url:      "https://example.com",
			expected: false,
		},
		{
			name:     "URL with path but no extension",
			url:      "https://example.com/api/v0/flags",
			expected: false,
		},
		{
			name:     "URL with different extension",
			url:      "https://example.com/flags.txt",
			expected: false,
		},
		{
			name:     "URL with .json in path but not at end",
			url:      "https://example.com/flags.json/export",
			expected: false,
		},
		{
			name:     "URL with uppercase extension",
			url:      "https://example.com/flags.JSON",
			expected: false, // Case sensitive check
		},
		{
			name:     "Empty URL",
			url:      "",
			expected: false,
		},
		{
			name:     "Short URL with .json",
			url:      ".json",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := URLLooksLikeAFile(tt.url)
			assert.Equal(t, tt.expected, result, "URLLooksLikeAFile(%q) should return %v", tt.url, tt.expected)
		})
	}
}

func TestWriteManifestEndsWithNewline(t *testing.T) {
	// Use an in-memory filesystem so we don't touch disk
	memFs := afero.NewMemMapFs()
	filesystem.SetFileSystem(memFs)
	t.Cleanup(func() { filesystem.SetFileSystem(afero.NewOsFs()) })

	manifest := createInitManifest(map[string]any{})
	path := "/flags.json"

	err := writeManifest(path, manifest)
	require.NoError(t, err)

	data, err := afero.ReadFile(memFs, path)
	require.NoError(t, err)

	assert.True(t, len(data) > 0, "manifest file should not be empty")
	assert.Equal(t, byte('\n'), data[len(data)-1], "manifest file should end with a newline")
}
