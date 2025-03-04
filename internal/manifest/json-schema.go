package manifest

import (
	"log"

	"github.com/invopop/jsonschema"
	"github.com/pterm/pterm"
)

// Converts the Manifest struct to a JSON schema.
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