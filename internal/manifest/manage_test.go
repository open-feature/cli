package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			url:      "https://example.com/api/v1/flags.json",
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
			url:      "https://example.com/api/v1/flags",
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
