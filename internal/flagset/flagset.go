package flagset

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/open-feature/cli/internal/filesystem"
	"github.com/open-feature/cli/internal/manifest"
	"github.com/spf13/afero"
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
		return "int"
	case FloatType:
		return "float"
	case BoolType:
		return "bool"
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

// Loads, validates, and unmarshals the manifest file at the given path into a flagset
func Load(manifestPath string) (*Flagset, error) {
	fs := filesystem.FileSystem()
	data, err := afero.ReadFile(fs, manifestPath)
	if err != nil {
		return nil, fmt.Errorf("error reading contents from file %q", manifestPath)
	}

	validationErrors, err := manifest.Validate(data)
	if err != nil {
		return nil, err
	} else if len(validationErrors) > 0 {
		return nil, errors.New(FormatValidationError(validationErrors))
	}

	var flagset Flagset
	if err := json.Unmarshal(data, &flagset); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", validationErrors)
	}

	return &flagset, nil
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
		var flagType FlagType
		switch flag.FlagType {
		case "integer":
			flagType = IntType
		case "float":
			flagType = FloatType
		case "boolean":
			flagType = BoolType
		case "string":
			flagType = StringType
		case "object":
			flagType = ObjectType
		default:
			return errors.New("unknown flag type")
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
func FormatValidationError(issues []manifest.ValidationError) string {
	var sb strings.Builder
	sb.WriteString("flag manifest validation failed:\n\n")

	// Group messages by flag path
	grouped := make(map[string]struct {
		flagType string
		messages []string
	})

	for _, issue := range issues {
		entry := grouped[issue.Path]
		entry.flagType = issue.Type
		entry.messages = append(entry.messages, issue.Message)
		grouped[issue.Path] = entry
	}

	// Sort paths for consistent output
	paths := make([]string, 0, len(grouped))
	for path := range grouped {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	// Format each row
	for _, path := range paths {
		entry := grouped[path]
		flagType := entry.flagType
		if flagType == "" {
			flagType = "missing"
		}
		sb.WriteString(fmt.Sprintf(
			"- flagType: %s\n  flagPath: %s\n  errors:\n    ~ %s\n  \tSuggestions:\n      \t- flagType: boolean\n      \t- defaultValue: true\n\n",
			flagType,
			path,
			strings.Join(entry.messages, "\n    ~ "),
		))
	}
	return sb.String()
}
