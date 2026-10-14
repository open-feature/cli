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

// ObjectSchema represents a JSON Schema subset for describing the shape of object flag values.
// It supports objects with properties, arrays with item schemas, and primitive types.
type ObjectSchema struct {
	// Type is the JSON Schema type: "object", "array", "string", "number", "integer", "boolean"
	Type string `json:"type"`
	// Properties defines the properties of an object type (only valid when Type is "object")
	Properties map[string]*ObjectSchema `json:"properties,omitempty"`
	// Required lists property names that must be present (only valid when Type is "object")
	Required []string `json:"required,omitempty"`
	// Items defines the schema for array elements (only valid when Type is "array")
	Items *ObjectSchema `json:"items,omitempty"`
	// AdditionalProperties controls whether extra properties are allowed (only valid when Type is "object")
	AdditionalProperties *bool `json:"additionalProperties,omitempty"`
}

type Flag struct {
	Key          string
	Type         FlagType
	Description  string
	DefaultValue any
	// Schema is an optional JSON Schema subset describing the shape of an object flag's value.
	// It is nil when no schema is provided (backward compatible).
	Schema *ObjectSchema
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

// ParseFlagType converts a string flag type to FlagType enum
func ParseFlagType(typeStr string) (FlagType, error) {
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
			FlagType     string        `json:"flagType"`
			Description  string        `json:"description"`
			DefaultValue any           `json:"defaultValue"`
			Schema       *ObjectSchema `json:"schema,omitempty"`
		} `json:"flags"`
	}

	if err := json.Unmarshal(data, &manifest); err != nil {
		return err
	}

	for key, flag := range manifest.Flags {
		flagType, err := ParseFlagType(flag.FlagType)
		if err != nil {
			return err
		}

		f := Flag{
			Key:          key,
			Type:         flagType,
			Description:  flag.Description,
			DefaultValue: flag.DefaultValue,
		}

		// Only attach schema to object flags
		if flagType == ObjectType && flag.Schema != nil {
			f.Schema = flag.Schema
		}

		fs.Flags = append(fs.Flags, f)
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
			FlagType     string        `json:"flagType"`
			Description  string        `json:"description"`
			DefaultValue any           `json:"defaultValue"`
			Schema       *ObjectSchema `json:"schema,omitempty"`
		} `json:"flags"`
	}{
		Flags: make(map[string]struct {
			FlagType     string        `json:"flagType"`
			Description  string        `json:"description"`
			DefaultValue any           `json:"defaultValue"`
			Schema       *ObjectSchema `json:"schema,omitempty"`
		}),
	}

	for _, flag := range fs.Flags {
		manifest.Flags[flag.Key] = struct {
			FlagType     string        `json:"flagType"`
			Description  string        `json:"description"`
			DefaultValue any           `json:"defaultValue"`
			Schema       *ObjectSchema `json:"schema,omitempty"`
		}{
			FlagType:     flag.Type.String(),
			Description:  flag.Description,
			DefaultValue: flag.DefaultValue,
			Schema:       flag.Schema,
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

	if err := json.Unmarshal(data, &sourceWithWrapper); err == nil && sourceWithWrapper.Flags != nil {
		// Successfully unmarshaled as object with flags property (even if empty)
		sourceFlagsArray = sourceWithWrapper.Flags
	} else {
		// Try to unmarshal as a direct array of flags (for backward compatibility)
		if err := json.Unmarshal(data, &sourceFlagsArray); err != nil {
			return nil, fmt.Errorf("failed to parse flags: expected either {\"flags\": [...]} or direct array format")
		}
	}

	var flags []Flag
	for _, sf := range sourceFlagsArray {
		flagType, err := ParseFlagType(sf.Type)
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
