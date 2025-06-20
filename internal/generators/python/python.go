package python

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"
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
	case flagset.ObjectType:
		return "object"
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

func pythonBoolLiteral(value any) any {
	if v, ok := value.(bool); ok {
		if v {
			return "True"
		}
		return "False"
	}
	return value
}

func toPythonDict(value any) string {
	assertedMap, ok := value.(map[string]any)
	if !ok {
		return "None"
	}

	// To have a determined order of the object for comparison
	keys := slices.Sorted(maps.Keys(assertedMap))

	var builder strings.Builder
	builder.WriteString("{")

	for index, key := range keys {
		if index != 0 {
			builder.WriteString(", ")
		}
		val := assertedMap[key]

		builder.WriteString(fmt.Sprintf(`%q: %s`, key, formatNestedValue(val)))
	}

	builder.WriteString("}")
	return builder.String()
}

func formatNestedValue(value any) string {
	switch val := value.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case bool:
		return fmt.Sprintf(pythonBoolLiteral(val).(string))
	case int, int64, float64:
		return fmt.Sprintf("%v", val)
	case map[string]any:
		return toPythonDict(val)
	case []any:
		var sliceBuilder strings.Builder
		sliceBuilder.WriteString("[")
		for index, elem := range val {
			if index > 0 {
				sliceBuilder.WriteString(", ")
			}

			sliceBuilder.WriteString(formatNestedValue(elem))
		}
		sliceBuilder.WriteString("]")
		return sliceBuilder.String()
	default:
		jsonBytes, err := json.Marshal(val)
		if err != nil {
			return "None"
		}
		return strings.ReplaceAll(string(jsonBytes), "null", "None")
	}

}

func (g *PythonGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"OpenFeatureType":         openFeatureType,
		"TypedGetMethodSync":      typedGetMethodSync,
		"TypedGetMethodAsync":     typedGetMethodAsync,
		"TypedDetailsMethodSync":  typedDetailsMethodSync,
		"TypedDetailsMethodAsync": typedDetailsMethodAsync,
		"PythonBoolLiteral":       pythonBoolLiteral,
		"ToPythonDict":            toPythonDict,
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
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{}),
	}
}
