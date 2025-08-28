package flagset

import (
	"encoding/json"
	"fmt"
	"sort"
)

// FlagType are the primitive types of flags.
type FlagType int

// Collection of the different kinds of flag types
const (
	UnknownFlagType FlagType = iota
	IntType
	FloatType
	BoolType
	StringType
	ObjectType
)

func (f FlagType) String() string {
	switch f {
	case IntType:
		return "integer"
	case FloatType:
		return "float"
	case BoolType:
		return "boolean"
	case StringType:
		return "string"
	case ObjectType:
		return "object"
	default:
		return "unknown"
	}
}

type Flag struct {
	Key          string
	Type         FlagType
	Description  string
	DefaultValue any
}

type Flagset struct {
	Flags []Flag
}

// Filter removes flags from the Flagset that are of unsupported types.
func (fs *Flagset) Filter(unsupportedFlagTypes map[FlagType]bool) *Flagset {
	var filtered Flagset
	for _, flag := range fs.Flags {
		if !unsupportedFlagTypes[flag.Type] {
			filtered.Flags = append(filtered.Flags, flag)
		}
	}
	return &filtered
}

// parseFlagType converts a string flag type to FlagType enum
func parseFlagType(typeStr string) (FlagType, error) {
	switch typeStr {
	case "integer", "Integer":
		return IntType, nil
	case "float", "Float", "Number":
		return FloatType, nil
	case "boolean", "bool", "Boolean":
		return BoolType, nil
	case "string", "String":
		return StringType, nil
	case "object", "Object", "JSON":
		return ObjectType, nil
	default:
		return UnknownFlagType, fmt.Errorf("unknown flag type: %s", typeStr)
	}
}

// UnmarshalJSON unmarshals the JSON data into a Flagset. It is used by json.Unmarshal.
func (fs *Flagset) UnmarshalJSON(data []byte) error {
	var manifest struct {
		Flags map[string]struct {
			FlagType     string `json:"flagType"`
			Description  string `json:"description"`
			DefaultValue any    `json:"defaultValue"`
		} `json:"flags"`
	}

	if err := json.Unmarshal(data, &manifest); err != nil {
		return err
	}

	for key, flag := range manifest.Flags {
		flagType, err := parseFlagType(flag.FlagType)
		if err != nil {
			return err
		}

		fs.Flags = append(fs.Flags, Flag{
			Key:          key,
			Type:         flagType,
			Description:  flag.Description,
			DefaultValue: flag.DefaultValue,
		})
	}

	// Ensure consistency of order of flag generation.
	sort.Slice(fs.Flags, func(i, j int) bool {
		return fs.Flags[i].Key < fs.Flags[j].Key
	})

	return nil
}

// MarshalJSON marshals a Flagset into JSON format compatible with the manifest structure
func (fs *Flagset) MarshalJSON() ([]byte, error) {
	manifest := struct {
		Flags map[string]struct {
			FlagType     string `json:"flagType"`
			Description  string `json:"description"`
			DefaultValue any    `json:"defaultValue"`
		} `json:"flags"`
	}{
		Flags: make(map[string]struct {
			FlagType     string `json:"flagType"`
			Description  string `json:"description"`
			DefaultValue any    `json:"defaultValue"`
		}),
	}

	for _, flag := range fs.Flags {
		manifest.Flags[flag.Key] = struct {
			FlagType     string `json:"flagType"`
			Description  string `json:"description"`
			DefaultValue any    `json:"defaultValue"`
		}{
			FlagType:     flag.Type.String(),
			Description:  flag.Description,
			DefaultValue: flag.DefaultValue,
		}
	}

	return json.Marshal(manifest)
}

func LoadFromSourceFlags(data []byte) (*[]Flag, error) {
	type SourceFlag struct {
		Key          string `json:"key"`
		Type         string `json:"type"`
		Description  string `json:"description"`
		DefaultValue any    `json:"defaultValue"`
	}

	// First try to unmarshal as an object with a "flags" property
	var sourceWithWrapper struct {
		Flags []SourceFlag `json:"flags"`
	}

	var sourceFlagsArray []SourceFlag

	if err := json.Unmarshal(data, &sourceWithWrapper); err == nil && len(sourceWithWrapper.Flags) > 0 {
		// Successfully unmarshaled as object with flags property
		sourceFlagsArray = sourceWithWrapper.Flags
	} else {
		// Try to unmarshal as a direct array of flags (for backward compatibility)
		if err := json.Unmarshal(data, &sourceFlagsArray); err != nil {
			return nil, fmt.Errorf("failed to parse flags: expected either {\"flags\": [...]} or direct array format")
		}
	}

	var flags []Flag
	for _, sf := range sourceFlagsArray {
		flagType, err := parseFlagType(sf.Type)
		if err != nil {
			return nil, err
		}

		flags = append(flags, Flag{
			Key:          sf.Key,
			Type:         flagType,
			Description:  sf.Description,
			DefaultValue: sf.DefaultValue,
		})
	}

	return &flags, nil
}
