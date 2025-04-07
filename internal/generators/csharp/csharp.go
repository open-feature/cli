package csharp

import (
	_ "embed"
	"fmt"
	"text/template"

	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators"
)

type CsharpGenerator struct {
	generators.CommonGenerator
}

type Params struct {
	// Add C# specific parameters here if needed
	Namespace string
}

//go:embed csharp.tmpl
var csharpTmpl string

func openFeatureType(t flagset.FlagType) string {
	switch t {
	case flagset.IntType:
		return "int"
	case flagset.FloatType:
		return "double" // .NET uses double, not float
	case flagset.BoolType:
		return "bool"
	case flagset.StringType:
		return "string"
	default:
		return ""
	}
}

func formatDefaultValue(flag flagset.Flag) string {
	switch flag.Type {
	case flagset.StringType:
		return fmt.Sprintf("\"%s\"", flag.DefaultValue)
	case flagset.BoolType:
		if flag.DefaultValue == true {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", flag.DefaultValue)
	}
}

func (g *CsharpGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"OpenFeatureType":    openFeatureType,
		"FormatDefaultValue": formatDefaultValue,
	}

	newParams := &generators.Params[any]{
		OutputPath: params.OutputPath,
		Custom:     params.Custom,
	}

	return g.GenerateFile(funcs, csharpTmpl, newParams, "OpenFeature.g.cs")
}

// NewGenerator creates a generator for C#.
func NewGenerator(fs *flagset.Flagset) *CsharpGenerator {
	return &CsharpGenerator{
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{
			flagset.ObjectType: true,
		}),
	}
}
