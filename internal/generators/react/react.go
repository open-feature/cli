package react

import (
	_ "embed"
	"fmt"
	"strconv"
	"text/template"

	"github.com/open-feature/cli/internal/generators"
	"github.com/open-feature/cli/internal/manifest"

	"github.com/iancoleman/strcase"
)

// Params are parameters for creating a Generator
type Params struct {
	generators.CommonParams
	// Test string
}

type ReactGenerator struct {
	generators.CommonGenerator
}

//go:embed react.tmpl
var reactTmpl string

func flagVarName(flagName string) string {
	return strcase.ToCamel(flagName)
}

func flagInitParam(flagName string) string {
	return strconv.Quote(flagName)
}

func defaultValueLiteral(flag *generators.Flag) string {
	switch manifest.FlagType(flag.Type) {
	case manifest.StringType:
		return strconv.Quote(flag.DefaultValue.(string))
	default:
		// TODO fix this
		return fmt.Sprintf("%v", flag.DefaultValue)
	}
}

func (g *ReactGenerator) Generate(params Params) error {
	funcs := template.FuncMap{
		"FlagVarName":         flagVarName,
		"FlagInitParam":       flagInitParam,
		"DefaultValueLiteral": defaultValueLiteral,
		// "TypeString":          typeString,
	}

	return g.GenerateFile(funcs, reactTmpl, generators.CommonParams{
		OutputPath: params.OutputPath,
		Params:     params.Params,
	})
}

// NewGenerator creates a generator for React.
func NewGenerator(manifest *manifest.Manifest) *ReactGenerator {
	return &ReactGenerator{
		CommonGenerator: *generators.NewCommonGenerator(
			manifest,
			generators.WithStability(generators.Alpha),
		),
	}
}
