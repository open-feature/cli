package manifest

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"

	"strings"

	"github.com/invopop/jsonschema"
	"github.com/open-feature/cli/schema/v0"
	"github.com/pterm/pterm"
	"github.com/xeipuuv/gojsonschema"
)

func ToJSONSchema() *jsonschema.Schema {
	reflector := &jsonschema.Reflector{
		ExpandedStruct: true,
		AllowAdditionalProperties: true,
		BaseSchemaID: "openfeature-cli",
	}

	if err := reflector.AddGoComments("github.com/open-feature/cli", "./internal/manifest"); err != nil {
		pterm.Error.Printf("Error extracting comments from types.go: %v\n", err)
	}

	schema := reflector.Reflect(Manifest{})
	schema.Version = "http://json-schema.org/draft-07/schema#"
	schema.Title = "OpenFeature CLI Manifest"
	flags, ok := schema.Properties.Get("flags")
	if !ok {
		log.Fatal("flags not found")
	}
	flags.PatternProperties = map[string]*jsonschema.Schema{
		"^.{1,}$": {
			Ref: "#/$defs/flag",
		},
	}
	// We only want flags keys that matches the pattern properties
	flags.AdditionalProperties = jsonschema.FalseSchema

	schema.Definitions = jsonschema.Definitions{
		"flag": &jsonschema.Schema{
			OneOf: []*jsonschema.Schema{
				{Ref: "#/$defs/booleanFlag"},
				{Ref: "#/$defs/stringFlag"},
				{Ref: "#/$defs/integerFlag"},
				{Ref: "#/$defs/floatFlag"},
				{Ref: "#/$defs/objectFlag"},
			},
			Required: []string{"flagType", "defaultValue"},
		},
		"booleanFlag": &jsonschema.Schema{
			Type:       "object",
			Properties: reflector.Reflect(BooleanFlag{}).Properties,
		},
		"stringFlag": &jsonschema.Schema{
			Type:       "object",
			Properties: reflector.Reflect(StringFlag{}).Properties,
		},
		"integerFlag": &jsonschema.Schema{
			Type:       "object",
			Properties: reflector.Reflect(IntegerFlag{}).Properties,
		},
		"floatFlag": &jsonschema.Schema{
			Type:       "object",
			Properties: reflector.Reflect(FloatFlag{}).Properties,
		},
		"objectFlag": &jsonschema.Schema{
			Type:       "object",
			Properties: reflector.Reflect(ObjectFlag{}).Properties,
		},
	}

	return schema
}

func (m *Manifest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Manifest) Marshal() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

type ValidationError struct {
	Type    string `json:"type"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

func (m *Manifest) Validate() ([]ValidationError, error) {
	schemaLoader := gojsonschema.NewStringLoader(schema.SchemaFile)
	manifestLoader := gojsonschema.NewGoLoader(m)

	result, err := gojsonschema.Validate(schemaLoader, manifestLoader)
	if err != nil {
		return nil, fmt.Errorf("failed to validate manifest: %w", err)
	}

	var issues []ValidationError
	for _, err := range result.Errors() {
		if strings.HasPrefix(err.Field(), "flags") && err.Type() == "number_one_of" {
			issues = append(issues, ValidationError{
				Type:    err.Type(),
				Path:    err.Field(),
				Message: "flagType must be 'boolean', 'string', 'integer', 'float', or 'object'",
			})
		} else {
			issues = append(issues, ValidationError{
				Type:    err.Type(),
				Path:    err.Field(),
				Message: err.Description(),
			})
		}
	}

	return issues, nil
}

// type Change struct {
// 	Type     string `json:"type"`
// 	Path     string `json:"path"`
// 	OldValue any    `json:"oldValue,omitempty"`
// 	NewValue any    `json:"newValue,omitempty"`
// }

// func Compare(oldManifest, newManifest *Manifest) ([]Change, error) {
// 	var changes []Change
// 	oldFlags := oldManifest.Flags
// 	newFlags := newManifest.Flags

// 	for key, newFlag := range newFlags {
// 		if oldFlag, exists := oldFlags[key]; exists {
// 			if !reflect.DeepEqual(oldFlag, newFlag) {
// 				changes = append(changes, Change{
// 					Type:     "change",
// 					Path:     fmt.Sprintf("flags.%s", key),
// 					OldValue: oldFlag,
// 					NewValue: newFlag,
// 				})
// 			}
// 		} else {
// 			changes = append(changes, Change{
// 				Type:     "add",
// 				Path:     fmt.Sprintf("flags.%s", key),
// 				NewValue: newFlag,
// 			})
// 		}
// 	}

// 	for key, oldFlag := range oldFlags {
// 		if _, exists := newFlags[key]; !exists {
// 			changes = append(changes, Change{
// 				Type:     "remove",
// 				Path:     fmt.Sprintf("flags.%s", key),
// 				OldValue: oldFlag,
// 			})
// 		}
// 	}

// 	return changes, nil
// }
