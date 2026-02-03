package manifest

import (
	"strings"
	"testing"
)

func TestValidate_DuplicateFlagKeys(t *testing.T) {
	tests := []struct {
		name           string
		manifest       string
		wantDuplicates []string
	}{
		{
			name: "no duplicates",
			manifest: `{
				"flags": {
					"flag-a": {"flagType": "boolean", "defaultValue": true},
					"flag-b": {"flagType": "string", "defaultValue": "hello"}
				}
			}`,
			wantDuplicates: nil,
		},
		{
			name: "single duplicate",
			manifest: `{
				"flags": {
					"my-flag": {"flagType": "boolean", "defaultValue": true},
					"my-flag": {"flagType": "string", "defaultValue": "hello"}
				}
			}`,
			wantDuplicates: []string{"my-flag"},
		},
		{
			name: "multiple duplicates",
			manifest: `{
				"flags": {
					"flag-a": {"flagType": "boolean", "defaultValue": true},
					"flag-b": {"flagType": "string", "defaultValue": "hello"},
					"flag-a": {"flagType": "integer", "defaultValue": 42},
					"flag-b": {"flagType": "float", "defaultValue": 3.14}
				}
			}`,
			wantDuplicates: []string{"flag-a", "flag-b"},
		},
		{
			name: "triple duplicate of same key",
			manifest: `{
				"flags": {
					"repeated": {"flagType": "boolean", "defaultValue": true},
					"repeated": {"flagType": "string", "defaultValue": "hello"},
					"repeated": {"flagType": "integer", "defaultValue": 42}
				}
			}`,
			wantDuplicates: []string{"repeated", "repeated"},
		},
		{
			name: "empty flags object",
			manifest: `{
				"flags": {}
			}`,
			wantDuplicates: nil,
		},
		{
			name: "manifest with schema field",
			manifest: `{
				"$schema": "https://example.com/schema.json",
				"flags": {
					"dup": {"flagType": "boolean", "defaultValue": true},
					"dup": {"flagType": "boolean", "defaultValue": false}
				}
			}`,
			wantDuplicates: []string{"dup"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := Validate([]byte(tt.manifest))
			if err != nil {
				t.Fatalf("Validate() error = %v", err)
			}

			var gotDuplicates []string
			for _, issue := range issues {
				if issue.Type == "duplicate_key" {
					// Extract the flag key from the path (format: "flags.key")
					parts := strings.SplitN(issue.Path, ".", 2)
					if len(parts) == 2 {
						gotDuplicates = append(gotDuplicates, parts[1])
					}
				}
			}

			if len(gotDuplicates) != len(tt.wantDuplicates) {
				t.Errorf("got %d duplicates, want %d", len(gotDuplicates), len(tt.wantDuplicates))
				t.Errorf("got duplicates: %v", gotDuplicates)
				t.Errorf("want duplicates: %v", tt.wantDuplicates)
				return
			}

			for i, want := range tt.wantDuplicates {
				if gotDuplicates[i] != want {
					t.Errorf("duplicate[%d] = %q, want %q", i, gotDuplicates[i], want)
				}
			}
		})
	}
}

func TestValidate_DuplicateKeyErrorMessage(t *testing.T) {
	manifest := `{
		"flags": {
			"my-flag": {"flagType": "boolean", "defaultValue": true},
			"my-flag": {"flagType": "string", "defaultValue": "hello"}
		}
	}`

	issues, err := Validate([]byte(manifest))
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	var found bool
	for _, issue := range issues {
		if issue.Type == "duplicate_key" {
			found = true
			if issue.Path != "flags.my-flag" {
				t.Errorf("expected path 'flags.my-flag', got %q", issue.Path)
			}
			expectedMsg := "flag 'my-flag' is defined multiple times in the manifest"
			if issue.Message != expectedMsg {
				t.Errorf("expected message %q, got %q", expectedMsg, issue.Message)
			}
		}
	}

	if !found {
		t.Error("expected to find a duplicate_key validation error")
	}
}

func TestFindDuplicateFlagKeys_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "invalid JSON",
			input:    "not valid json",
			expected: nil,
		},
		{
			name:     "array instead of object",
			input:    `["a", "b", "c"]`,
			expected: nil,
		},
		{
			name:     "no flags key",
			input:    `{"other": "value"}`,
			expected: nil,
		},
		{
			name:     "flags is not an object",
			input:    `{"flags": "string value"}`,
			expected: nil,
		},
		{
			name:     "flags is an array",
			input:    `{"flags": [1, 2, 3]}`,
			expected: nil,
		},
		{
			name:     "nested duplicates not detected in flag values",
			input:    `{"flags": {"flag1": {"nested": 1, "nested": 2}}}`,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findDuplicateFlagKeys([]byte(tt.input))
			if len(result) != len(tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

// Sample test for FormatValidationError
func TestFormatValidationError_SortsByPath(t *testing.T) {
	issues := []ValidationError{
		{Path: "zeta.flag", Type: "boolean", Message: "must not be empty"},
		{Path: "alpha.flag", Type: "string", Message: "invalid value"},
		{Path: "beta.flag", Type: "number", Message: "must be greater than zero"},
	}

	output := FormatValidationError(issues)

	// The output should mention 'alpha.flag' before 'beta.flag', and 'beta.flag' before 'zeta.flag'
	alphaIdx := strings.Index(output, "flagPath: alpha.flag")
	betaIdx := strings.Index(output, "flagPath: beta.flag")
	zetaIdx := strings.Index(output, "flagPath: zeta.flag")

	if alphaIdx >= betaIdx || betaIdx >= zetaIdx {
		t.Errorf("flag paths are not sorted: alphaIdx=%d, betaIdx=%d, zetaIdx=%d\nOutput:\n%s",
			alphaIdx, betaIdx, zetaIdx, output)
	}
}
