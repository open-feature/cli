package csharp

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
	case flagset.ObjectType:
		return "object"
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

func toCSharpDict(value any) string {
	assertedMap, ok := value.(map[string]any)
	if !ok {
		return "null"
	}

	keys := slices.Sorted(maps.Keys(assertedMap))

	var builder strings.Builder
	builder.WriteString("new Dictionary<string, object>{")

	for index, key := range keys {
		if index > 0 {
			builder.WriteString(", ")
		}
		val := assertedMap[key]

		builder.WriteString(fmt.Sprintf("%q, %s", key, formatNestedValue(val)))
	}
	builder.WriteString("}")

	return builder.String()
}

func formatNestedValue(value any) string {
	switch val := value.(type) {
	case string:
		flag := flagset.Flag{
			Type:         flagset.StringType,
			DefaultValue: val,
		}
		return formatDefaultValue(flag)
	case bool:
		flag := flagset.Flag{
			Type:         flagset.BoolType,
			DefaultValue: val,
		}
		return formatDefaultValue(flag)
	case int, int64:
		flag := flagset.Flag{
			Type:         flagset.IntType,
			DefaultValue: val,
		}
		return formatDefaultValue(flag)
	case float64:
		flag := flagset.Flag{
			Type:         flagset.FloatType,
			DefaultValue: val,
		}
		return formatDefaultValue(flag)
	case map[string]any:
		return toCSharpDict(val)
	case []any:
		var sliceBuilder strings.Builder
		sliceBuilder.WriteString("new List<object>{")
		for index, elem := range val {
			if index > 0 {
				sliceBuilder.WriteString(",")
			}

			sliceBuilder.WriteString(formatNestedValue(elem))
		}
		sliceBuilder.WriteString("}")
		return sliceBuilder.String()
	default:
		jsonBytes, err := json.Marshal(val)
		if err != nil {
			return "null"
		}
		return fmt.Sprintf("%q", string(jsonBytes))
	}
}

func (g *CsharpGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"OpenFeatureType":    openFeatureType,
		"FormatDefaultValue": formatDefaultValue,
		"ToCSharpDict":       toCSharpDict,
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
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{}),
	}
}
