package manifest

import (
	"reflect"
	"sort"
	"testing"
)

func TestCompareDifferentManifests(t *testing.T) {
	oldManifest := &Manifest{
		Flags: map[string]any{
			"flag1": "value1",
			"flag2": "value2",
		},
	}

	newManifest := &Manifest{
		Flags: map[string]any{
			"flag1": "value1",
			"flag2": "newValue2",
			"flag3": "value3",
		},
	}

	changes, err := Compare(oldManifest, newManifest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedChanges := []Change{
		{Type: "change", Path: "flags.flag2", OldValue: "value2", NewValue: "newValue2"},
		{Type: "add", Path: "flags.flag3", NewValue: "value3"},
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
			"flag1": "value1",
			"flag2": "value2",
		},
	}

	changes, err := Compare(manifest, manifest)
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
