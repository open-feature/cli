package python

import (
	_ "embed"
	"text/template"

	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators"
)

type PythonGenerator struct {
	generators.CommonGenerator
}

type Params struct {
}

//go:embed python.tmpl
var pythonTmpl string

func openFeatureType(t flagset.FlagType) string {
	switch t {
	case flagset.IntType:
		return "int"
	case flagset.FloatType:
		return "float"
	case flagset.BoolType:
		return "bool"
	case flagset.StringType:
		return "str"
	default:
		return "object"
	}
}

func methodType(flagType flagset.FlagType) string {
	switch flagType {
	case flagset.StringType:
		return "string"
	case flagset.IntType:
		return "integer"
	case flagset.BoolType:
		return "boolean"
	case flagset.FloatType:
		return "float"
	default:
		panic("unsupported flag type")
	}
}

func typedGetMethodSync(flagType flagset.FlagType) string {
	return "get_" + methodType(flagType) + "_value"
}

func typedGetMethodAsync(flagType flagset.FlagType) string {
	return "get_" + methodType(flagType) + "_value_async"
}

func typedDetailsMethodSync(flagType flagset.FlagType) string {
	return "get_" + methodType(flagType) + "_details"
}

func typedDetailsMethodAsync(flagType flagset.FlagType) string {
	return "get_" + methodType(flagType) + "_details_async"
}

func pythonBoolLiteral(value interface{}) interface{} {
	if v, ok := value.(bool); ok {
		if v {
			return "True"
		}
		return "False"
	}
	return value
}

func (g *PythonGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"OpenFeatureType":         openFeatureType,
		"TypedGetMethodSync":      typedGetMethodSync,
		"TypedGetMethodAsync":     typedGetMethodAsync,
		"TypedDetailsMethodSync":  typedDetailsMethodSync,
		"TypedDetailsMethodAsync": typedDetailsMethodAsync,
		"PythonBoolLiteral":       pythonBoolLiteral,
	}

	newParams := &generators.Params[any]{
		OutputPath: params.OutputPath,
		Custom:     Params{},
	}

	return g.GenerateFile(funcs, pythonTmpl, newParams, "openfeature.py")
}

// NewGenerator creates a generator for Python.
func NewGenerator(fs *flagset.Flagset) *PythonGenerator {
	return &PythonGenerator{
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{
			flagset.ObjectType: true,
		}),
	}
}
