package manifest

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/open-feature/cli/internal/filesystem"
	"github.com/spf13/afero"
)

func Load(manifestPath string) (*Manifest, error) {
	fs := filesystem.FileSystem()
	data, err := afero.ReadFile(fs, manifestPath)
	if err != nil {
		return nil, fmt.Errorf("error reading contents from file %q", manifestPath)
	}

	var raw interface{}
	err = json.Unmarshal(data, &raw)
	validationErrors, err := Validate(raw)
	if err != nil {
		return nil, err
	} else if len(validationErrors) > 0 {
		// TODO tease running manifest validate command
		return nil, fmt.Errorf("validation failed: %v", validationErrors)
	}

	return Unmarshal(data)
}

func Unmarshal(data []byte) (*Manifest, error) {
	// var raw Manifest;
	// err := json.Unmarshal(data, &raw)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to unmarshal manifest: %w", err)
	// }
	// return &raw, nil

	var rawManifest struct {
		Flags map[string]json.RawMessage `json:"flags"`
	}
	err := json.Unmarshal(data, &rawManifest)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal manifest: %w", err)
	}

	manifest := &Manifest{
		Flags: make(map[string]any),
	}

	for key, rawFlag := range rawManifest.Flags {
		var flagType struct {
			Type string `json:"flagType"`
		}
		err := json.Unmarshal(rawFlag, &flagType)
		if err != nil {
			return nil, err
		}

		switch flagType.Type {
		case "boolean":
			var booleanFlag BooleanFlag
			err := json.Unmarshal(rawFlag, &booleanFlag)
			if err != nil {
				return nil, err
			}
			manifest.Flags[key] = booleanFlag
		case "string":
			var stringFlag StringFlag
			err := json.Unmarshal(rawFlag, &stringFlag)
			if err != nil {
				return nil, err
			}
			manifest.Flags[key] = stringFlag
		case "integer":
			var integerFlag IntegerFlag
			err := json.Unmarshal(rawFlag, &integerFlag)
			if err != nil {
				return nil, err
			}
			manifest.Flags[key] = integerFlag
		case "float":
			var floatFlag FloatFlag
			err := json.Unmarshal(rawFlag, &floatFlag)
			if err != nil {
				return nil, err
			}
			manifest.Flags[key] = floatFlag
		case "object":
			var objectFlag ObjectFlag
			err := json.Unmarshal(rawFlag, &objectFlag)
			if err != nil {
				return nil, err
			}
			manifest.Flags[key] = objectFlag
		default:
			return nil, fmt.Errorf("unknown flag type: %s", flagType.Type)
		}
	}

	return manifest, nil
}

func (m *Manifest) Marshal() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}
