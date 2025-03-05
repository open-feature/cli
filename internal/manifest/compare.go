package manifest

import (
	"fmt"
	"reflect"
)

type Change struct {
	Type     string `json:"type"`
	Path     string `json:"path"`
	OldValue any    `json:"oldValue,omitempty"`
	NewValue any    `json:"newValue,omitempty"`
}

func Compare(oldManifest, newManifest *Manifest) ([]Change, error) {
	var changes []Change
	oldFlags := oldManifest.Flags
	newFlags := newManifest.Flags

	for key, newFlag := range newFlags {
		if oldFlag, exists := oldFlags[key]; exists {
			if !reflect.DeepEqual(oldFlag, newFlag) {
				changes = append(changes, Change{
					Type:     "change",
					Path:     fmt.Sprintf("flags.%s", key),
					OldValue: oldFlag,
					NewValue: newFlag,
				})
			}
		} else {
			changes = append(changes, Change{
				Type:     "add",
				Path:     fmt.Sprintf("flags.%s", key),
				NewValue: newFlag,
			})
		}
	}

	for key, oldFlag := range oldFlags {
		if _, exists := newFlags[key]; !exists {
			changes = append(changes, Change{
				Type:     "remove",
				Path:     fmt.Sprintf("flags.%s", key),
				OldValue: oldFlag,
			})
		}
	}

	return changes, nil
}
