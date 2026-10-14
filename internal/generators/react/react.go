package react

import (
	_ "embed"
	"encoding/json"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators"
)

type ReactGenerator struct {
	generators.CommonGenerator
}

type Params struct{}

//go:embed react.tmpl
var reactTmpl string

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

func toJSONString(value any) string {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

// reservedNames are symbols exported by the React generator itself. Flag keys
// whose generated hook name (use + ToPascal) matches one of these will be
// excluded from the generated output and a warning will be emitted.
var reservedNames = map[string]bool{
	"useOpenFeatureClient": true,
}

func (g *ReactGenerator) Generate(params *generators.Params[Params]) error {
	g.Flagset = generators.FilterReservedFlags(g.Flagset, "React", reservedNames, func(key string) string {
		return "use" + strcase.ToCamel(key)
	})

	funcs := template.FuncMap{
		"OpenFeatureType": openFeatureType,
		"ToJSONString":    toJSONString,
	}

	newParams := &generators.Params[any]{
		OutputPath:   params.OutputPath,
		TemplatePath: params.TemplatePath,
		Custom:       Params{},
	}

	return g.GenerateFile(funcs, reactTmpl, newParams, "openfeature.ts")
}

// NewGenerator creates a generator for React.
func NewGenerator(fs *flagset.Flagset) *ReactGenerator {
	return &ReactGenerator{
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{}),
	}
}
