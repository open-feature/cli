package manifest

import (
	"reflect"
	"sort"
	"testing"
)

func TestCompareDifferentManifests(t *testing.T) {
	oldManifest := &Manifest{
		Flags: map[string]any{
			"flag1": map[string]any{
				"flagType":     "string",
				"defaultValue": "value1",
			},
			"flag2": map[string]any{
				"flagType":     "string",
				"defaultValue": "value2",
			},
		},
	}

	newManifest := &Manifest{
		Flags: map[string]any{
			"flag1": map[string]any{
				"flagType":     "string",
				"defaultValue": "value1",
			},
			"flag2": map[string]any{
				"flagType":     "string",
				"defaultValue": "newValue2",
			},
			"flag3": map[string]any{
				"flagType":     "string",
				"defaultValue": "value3",
			},
		},
	}

	changes, err := Compare(oldManifest, newManifest, CompareOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedChanges := []Change{
		{Type: "change", Path: "flags.flag2", OldValue: map[string]any{
			"flagType":     "string",
			"defaultValue": "value2",
		}, NewValue: map[string]any{
			"flagType":     "string",
			"defaultValue": "newValue2",
		}},
		{Type: "add", Path: "flags.flag3", NewValue: map[string]any{
			"flagType":     "string",
			"defaultValue": "value3",
		}},
	}

	sortChanges(changes)
	sortChanges(expectedChanges)

	if !reflect.DeepEqual(changes, expectedChanges) {
		t.Errorf("expected %v, got %v", expectedChanges, changes)
	}
}

func TestCompareIdenticalManifests(t *testing.T) {
	manifest := &Manifest{
		Flags: map[string]any{
			"flag1": map[string]any{
				"flagType":     "string",
				"defaultValue": "value1",
			},
			"flag2": map[string]any{
				"flagType":     "boolean",
				"defaultValue": true,
			},
		},
	}

	changes, err := Compare(manifest, manifest, CompareOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(changes) != 0 {
		t.Errorf("expected no changes, got %v", changes)
	}
}

func sortChanges(changes []Change) {
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Path < changes[j].Path
	})
}

// Test that property order differences don't trigger changes
func TestComparePropertyOrderDifferences(t *testing.T) {
	oldManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": false,
				"description":  "Enable dark mode",
			},
		},
	}

	// Same content, different property order (maps don't have order in Go, but testing the logic)
	newManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"description":  "Enable dark mode",
				"flagType":     "boolean",
				"defaultValue": false,
			},
		},
	}

	changes, err := Compare(oldManifest, newManifest, CompareOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(changes) != 0 {
		t.Errorf("expected no changes for property order differences, got %d changes: %v", len(changes), changes)
	}
}

// Test that extra properties are silently ignored
func TestCompareExtraProperties(t *testing.T) {
	oldManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": false,
			},
		},
	}

	newManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": false,
				"extraProp":    "value",
			},
		},
	}

	changes, err := Compare(oldManifest, newManifest, CompareOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Extra properties should be silently ignored
	if len(changes) != 0 {
		t.Errorf("expected 0 changes (extra properties ignored), got %d: %v", len(changes), changes)
	}
}

// Test ignoring description with shorthand pattern
func TestCompareIgnoreDescriptionShorthand(t *testing.T) {
	oldManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": false,
				"description":  "Old description",
			},
		},
	}

	newManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": false,
				"description":  "New description",
			},
		},
	}

	// Without ignore - should detect change
	changes, err := Compare(oldManifest, newManifest, CompareOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changes) != 1 {
		t.Errorf("expected 1 change without ignore, got %d", len(changes))
	}

	// With ignore pattern - should not detect change
	changes, err = Compare(oldManifest, newManifest, CompareOptions{
		IgnorePatterns: []string{"description"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changes) != 0 {
		t.Errorf("expected 0 changes with ignore pattern, got %d: %v", len(changes), changes)
	}
}

// Test ignoring description with full path pattern
func TestCompareIgnoreDescriptionFullPath(t *testing.T) {
	oldManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": false,
				"description":  "Old description",
			},
			"featureX": map[string]any{
				"flagType":     "string",
				"defaultValue": "value",
				"description":  "Old feature description",
			},
		},
	}

	newManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": false,
				"description":  "New description",
			},
			"featureX": map[string]any{
				"flagType":     "string",
				"defaultValue": "value",
				"description":  "New feature description",
			},
		},
	}

	// With full path pattern - should not detect changes
	changes, err := Compare(oldManifest, newManifest, CompareOptions{
		IgnorePatterns: []string{"flags.*.description"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(changes) != 0 {
		t.Errorf("expected 0 changes with full path pattern, got %d: %v", len(changes), changes)
	}
}

// Test detecting functional changes
func TestCompareFunctionalChanges(t *testing.T) {
	oldManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": false,
				"description":  "Enable dark mode",
			},
		},
	}

	newManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": true, // Changed default value
				"description":  "Enable dark mode",
			},
		},
	}

	changes, err := Compare(oldManifest, newManifest, CompareOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(changes) != 1 {
		t.Errorf("expected 1 change for defaultValue change, got %d", len(changes))
	}

	if changes[0].Type != "change" {
		t.Errorf("expected change type, got %s", changes[0].Type)
	}
}

// Test detecting flagType changes
func TestCompareFlagTypeChange(t *testing.T) {
	oldManifest := &Manifest{
		Flags: map[string]any{
			"feature": map[string]any{
				"flagType":     "boolean",
				"defaultValue": false,
			},
		},
	}

	newManifest := &Manifest{
		Flags: map[string]any{
			"feature": map[string]any{
				"flagType":     "string",
				"defaultValue": "false",
			},
		},
	}

	changes, err := Compare(oldManifest, newManifest, CompareOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(changes) != 1 {
		t.Errorf("expected 1 change for flagType change, got %d", len(changes))
	}
}

// Test multiple ignore patterns
func TestCompareMultipleIgnorePatterns(t *testing.T) {
	oldManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": false,
				"description":  "Old description",
				"metadata": map[string]any{
					"author":    "Old author",
					"timestamp": "2023-01-01",
				},
			},
		},
	}

	newManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": false,
				"description":  "New description",
				"metadata": map[string]any{
					"author":    "New author",
					"timestamp": "2024-01-01",
				},
			},
		},
	}

	// Ignore both description and metadata fields
	changes, err := Compare(oldManifest, newManifest, CompareOptions{
		IgnorePatterns: []string{"description", "metadata.*"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(changes) != 0 {
		t.Errorf("expected 0 changes with multiple ignore patterns, got %d: %v", len(changes), changes)
	}
}

// Test wildcard pattern matching
func TestCompareWildcardPattern(t *testing.T) {
	oldManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": false,
				"metadata": map[string]any{
					"author": "John",
					"date":   "2023-01-01",
				},
			},
		},
	}

	newManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": false,
				"metadata": map[string]any{
					"author": "Jane",
					"date":   "2024-01-01",
				},
			},
		},
	}

	// Use wildcard to ignore all metadata fields
	changes, err := Compare(oldManifest, newManifest, CompareOptions{
		IgnorePatterns: []string{"flags.*.metadata.*"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(changes) != 0 {
		t.Errorf("expected 0 changes with wildcard pattern, got %d: %v", len(changes), changes)
	}
}

// Test that functional changes are still detected when ignoring non-functional fields
func TestCompareFunctionalChangesWithIgnore(t *testing.T) {
	oldManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": false,
				"description":  "Old description",
			},
		},
	}

	newManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": true, // Functional change
				"description":  "New description",
			},
		},
	}

	// Even when ignoring description, defaultValue change should be detected
	changes, err := Compare(oldManifest, newManifest, CompareOptions{
		IgnorePatterns: []string{"description"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(changes) != 1 {
		t.Errorf("expected 1 change for defaultValue, got %d", len(changes))
	}

	if changes[0].Type != "change" {
		t.Errorf("expected change type, got %s", changes[0].Type)
	}
}

// Test specific flag ignore pattern
func TestCompareIgnoreSpecificFlag(t *testing.T) {
	oldManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": false,
				"description":  "Dark mode description",
			},
			"featureX": map[string]any{
				"flagType":     "boolean",
				"defaultValue": true,
				"description":  "Feature X description",
			},
		},
	}

	newManifest := &Manifest{
		Flags: map[string]any{
			"darkMode": map[string]any{
				"flagType":     "boolean",
				"defaultValue": false,
				"description":  "Updated dark mode description",
			},
			"featureX": map[string]any{
				"flagType":     "boolean",
				"defaultValue": true,
				"description":  "Updated feature X description",
			},
		},
	}

	// Ignore only darkMode's description
	changes, err := Compare(oldManifest, newManifest, CompareOptions{
		IgnorePatterns: []string{"flags.darkMode.description"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should detect change only in featureX description
	if len(changes) != 1 {
		t.Errorf("expected 1 change, got %d: %v", len(changes), changes)
	}

	if changes[0].Path != "flags.featureX" {
		t.Errorf("expected change in featureX, got %s", changes[0].Path)
	}
}
