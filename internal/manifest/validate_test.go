package manifest

import (
	"strings"
	"testing"

	"github.com/open-feature/cli/internal/flagset"
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

func TestValidateDefaultValues(t *testing.T) {
	boolTrue := true
	boolFalse := false

	tests := []struct {
		name       string
		flags      []flagset.Flag
		wantErrors int
	}{
		{
			name: "no schema - no validation",
			flags: []flagset.Flag{
				{Key: "obj", Type: flagset.ObjectType, DefaultValue: map[string]any{"a": 1}, Schema: nil},
			},
			wantErrors: 0,
		},
		{
			name: "valid object matches schema",
			flags: []flagset.Flag{
				{
					Key:          "theme",
					Type:         flagset.ObjectType,
					DefaultValue: map[string]any{"color": "#fff"},
					Schema: &flagset.ObjectSchema{
						Type: "object",
						Properties: map[string]*flagset.ObjectSchema{
							"color": {Type: "string"},
						},
						Required: []string{"color"},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "missing required property",
			flags: []flagset.Flag{
				{
					Key:          "theme",
					Type:         flagset.ObjectType,
					DefaultValue: map[string]any{"other": "val"},
					Schema: &flagset.ObjectSchema{
						Type: "object",
						Properties: map[string]*flagset.ObjectSchema{
							"color": {Type: "string"},
						},
						Required: []string{"color"},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "wrong property type",
			flags: []flagset.Flag{
				{
					Key:          "theme",
					Type:         flagset.ObjectType,
					DefaultValue: map[string]any{"color": 123},
					Schema: &flagset.ObjectSchema{
						Type: "object",
						Properties: map[string]*flagset.ObjectSchema{
							"color": {Type: "string"},
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "additional properties disallowed",
			flags: []flagset.Flag{
				{
					Key:          "theme",
					Type:         flagset.ObjectType,
					DefaultValue: map[string]any{"color": "#fff", "extra": "bad"},
					Schema: &flagset.ObjectSchema{
						Type: "object",
						Properties: map[string]*flagset.ObjectSchema{
							"color": {Type: "string"},
						},
						AdditionalProperties: &boolFalse,
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "additional properties allowed explicitly",
			flags: []flagset.Flag{
				{
					Key:          "theme",
					Type:         flagset.ObjectType,
					DefaultValue: map[string]any{"color": "#fff", "extra": "ok"},
					Schema: &flagset.ObjectSchema{
						Type: "object",
						Properties: map[string]*flagset.ObjectSchema{
							"color": {Type: "string"},
						},
						AdditionalProperties: &boolTrue,
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "nested object validation",
			flags: []flagset.Flag{
				{
					Key:  "config",
					Type: flagset.ObjectType,
					DefaultValue: map[string]any{
						"layout": map[string]any{"fontSize": 12.0},
					},
					Schema: &flagset.ObjectSchema{
						Type: "object",
						Properties: map[string]*flagset.ObjectSchema{
							"layout": {
								Type: "object",
								Properties: map[string]*flagset.ObjectSchema{
									"fontSize": {Type: "integer"},
								},
							},
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "nested object wrong type",
			flags: []flagset.Flag{
				{
					Key:  "config",
					Type: flagset.ObjectType,
					DefaultValue: map[string]any{
						"layout": map[string]any{"fontSize": "not-a-number"},
					},
					Schema: &flagset.ObjectSchema{
						Type: "object",
						Properties: map[string]*flagset.ObjectSchema{
							"layout": {
								Type: "object",
								Properties: map[string]*flagset.ObjectSchema{
									"fontSize": {Type: "integer"},
								},
							},
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "array validation passes",
			flags: []flagset.Flag{
				{
					Key:          "tags",
					Type:         flagset.ObjectType,
					DefaultValue: map[string]any{"items": []any{"a", "b"}},
					Schema: &flagset.ObjectSchema{
						Type: "object",
						Properties: map[string]*flagset.ObjectSchema{
							"items": {Type: "array", Items: &flagset.ObjectSchema{Type: "string"}},
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "array element wrong type",
			flags: []flagset.Flag{
				{
					Key:          "tags",
					Type:         flagset.ObjectType,
					DefaultValue: map[string]any{"items": []any{"a", 123}},
					Schema: &flagset.ObjectSchema{
						Type: "object",
						Properties: map[string]*flagset.ObjectSchema{
							"items": {Type: "array", Items: &flagset.ObjectSchema{Type: "string"}},
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "defaultValue is not an object when schema expects object",
			flags: []flagset.Flag{
				{
					Key:          "bad",
					Type:         flagset.ObjectType,
					DefaultValue: "not-an-object",
					Schema: &flagset.ObjectSchema{
						Type: "object",
						Properties: map[string]*flagset.ObjectSchema{
							"a": {Type: "string"},
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "float passes as number",
			flags: []flagset.Flag{
				{
					Key:          "config",
					Type:         flagset.ObjectType,
					DefaultValue: map[string]any{"ratio": 0.5},
					Schema: &flagset.ObjectSchema{
						Type: "object",
						Properties: map[string]*flagset.ObjectSchema{
							"ratio": {Type: "number"},
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "integer as float64 with no fraction passes integer check",
			flags: []flagset.Flag{
				{
					Key:          "config",
					Type:         flagset.ObjectType,
					DefaultValue: map[string]any{"count": float64(10)},
					Schema: &flagset.ObjectSchema{
						Type: "object",
						Properties: map[string]*flagset.ObjectSchema{
							"count": {Type: "integer"},
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "float64 with fraction fails integer check",
			flags: []flagset.Flag{
				{
					Key:          "config",
					Type:         flagset.ObjectType,
					DefaultValue: map[string]any{"count": 10.5},
					Schema: &flagset.ObjectSchema{
						Type: "object",
						Properties: map[string]*flagset.ObjectSchema{
							"count": {Type: "integer"},
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "boolean validation",
			flags: []flagset.Flag{
				{
					Key:          "config",
					Type:         flagset.ObjectType,
					DefaultValue: map[string]any{"enabled": true},
					Schema: &flagset.ObjectSchema{
						Type: "object",
						Properties: map[string]*flagset.ObjectSchema{
							"enabled": {Type: "boolean"},
						},
					},
				},
			},
			wantErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := ValidateDefaultValues(tt.flags)
			if len(issues) != tt.wantErrors {
				t.Errorf("got %d validation errors, want %d: %v", len(issues), tt.wantErrors, issues)
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
