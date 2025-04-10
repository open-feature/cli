package nestjs

import (
	_ "embed"
	"text/template"

	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators"
)

type NestJsGenerator struct {
	generators.CommonGenerator
}

type Params struct {
}

//go:embed nestjs.tmpl
var nestJsTmpl string

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

func (g *NestJsGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"OpenFeatureType": openFeatureType,
	}

	newParams := &generators.Params[any]{
		OutputPath: params.OutputPath,
		Custom:     Params{},
	}

	return g.GenerateFile(funcs, nestJsTmpl, newParams, "openfeature-decorators.ts")
}

// NewGenerator creates a generator for NestJS.
func NewGenerator(fs *flagset.Flagset) *NestJsGenerator {
	return &NestJsGenerator{
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{
			flagset.ObjectType: true,
		}),
	}
}
