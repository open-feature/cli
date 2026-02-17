package angular

import (
	_ "embed"
	"encoding/json"
	"text/template"

	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators"
)

// AngularGenerator generates typesafe Angular services and directives.
type AngularGenerator struct {
	generators.CommonGenerator
}

// Params holds Angular-specific generation parameters.
type Params struct{}

//go:embed angular.tmpl
var angularTmpl string

// openFeatureType maps flagset types to TypeScript types.
func openFeatureType(t flagset.FlagType) string {
	switch t {
	case flagset.IntType:
		fallthrough
	case flagset.FloatType:
		return "number"
	case flagset.BoolType:
		return "boolean"
	case flagset.StringType:
		return "string"
	case flagset.ObjectType:
		return "object"
	default:
		return ""
	}
}

// sdkServiceMethod returns the corresponding FeatureFlagService method name for a flag type.
func sdkServiceMethod(t flagset.FlagType) string {
	switch t {
	case flagset.IntType:
		fallthrough
	case flagset.FloatType:
		return "getNumberDetails"
	case flagset.BoolType:
		return "getBooleanDetails"
	case flagset.StringType:
		return "getStringDetails"
	case flagset.ObjectType:
		return "getObjectDetails"
	default:
		return ""
	}
}

// toJSONString converts a value to its JSON string representation.
func toJSONString(value any) string {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

// Generate creates the Angular typesafe client file.
func (g *AngularGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"OpenFeatureType":  openFeatureType,
		"SdkServiceMethod": sdkServiceMethod,
		"ToJSONString":     toJSONString,
	}

	newParams := &generators.Params[any]{
		OutputPath:   params.OutputPath,
		TemplatePath: params.TemplatePath,
		Custom:       Params{},
	}

	return g.GenerateFile(funcs, angularTmpl, newParams, "openfeature.generated.ts")
}

// NewGenerator creates a generator for Angular.
func NewGenerator(fs *flagset.Flagset) *AngularGenerator {
	return &AngularGenerator{
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{}),
	}
}
