package golang

import (
	_ "embed"
	"sort"
	"text/template"

	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators"
)

type GolangGenerator struct {
	generators.CommonGenerator
}

type Params struct {
	GoPackage string
}

//go:embed golang.tmpl
var golangTmpl string

func openFeatureType(t flagset.FlagType) string {
	switch t {
	case flagset.IntType:
		return "Int"
	case flagset.FloatType:
		return "Float"
	case flagset.BoolType:
		return "Boolean"
	case flagset.StringType:
		return "String"
	default:
		return ""
	}
}

func typeString(flagType flagset.FlagType) string {
	switch flagType {
	case flagset.StringType:
		return "string"
	case flagset.IntType:
		return "int64"
	case flagset.BoolType:
		return "bool"
	case flagset.FloatType:
		return "float64"
	default:
		return ""
	}
}

func supportImports(flags []flagset.Flag) []string {
	var res []string
	if len(flags) > 0 {
		res = append(res, "\"context\"")
		res = append(res, "\"github.com/open-feature/go-sdk/openfeature\"")
	}
	sort.Strings(res)
	return res
}

func (g *GolangGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"SupportImports":  supportImports,
		"OpenFeatureType": openFeatureType,
		"TypeString":      typeString,
	}

	newParams := &generators.Params[any]{
		OutputPath: params.OutputPath,
		Custom: Params{
			GoPackage: params.Custom.GoPackage,
		},
	}

	return g.GenerateFile(funcs, golangTmpl, newParams, params.Custom.GoPackage+".go")
}

// NewGenerator creates a generator for Go.
func NewGenerator(fs *flagset.Flagset) *GolangGenerator {
	return &GolangGenerator{
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{
			flagset.ObjectType: true,
		}),
	}
}
