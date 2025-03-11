package react

import (
	_ "embed"
	"text/template"

	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators"
)

type ReactGenerator struct {
	generators.CommonGenerator
}

type Params struct {
}

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
	default:
		return ""
	}
}

func (g *ReactGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"OpenFeatureType": openFeatureType,
	}

	newParams := &generators.Params[any]{
		OutputPath: params.OutputPath,
		Custom:     Params{},
	}

	return g.GenerateFile(funcs, reactTmpl, newParams, "openfeature.ts")
}

// NewGenerator creates a generator for React.
func NewGenerator(fs *flagset.Flagset) *ReactGenerator {
	return &ReactGenerator{
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{
			flagset.ObjectType: true,
		}),
	}
}
