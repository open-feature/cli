package java

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

type JavaGenerator struct {
	generators.CommonGenerator
}

type Params struct {
	// Add Java parameters here if needed
	JavaPackage string
}

//go:embed java.tmpl
var javaTmpl string

func openFeatureType(t flagset.FlagType) string {
	switch t {
	case flagset.IntType:
		return "Integer"
	case flagset.FloatType:
		return "Double" //using Double as per openfeature Java-SDK
	case flagset.BoolType:
		return "Boolean"
	case flagset.StringType:
		return "String"
	case flagset.ObjectType:
		return "Object"
	default:
		return ""
	}
}

func formatDefaultValueForJava(flag flagset.Flag) string {
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

func toMapLiteral(value any) string {
	assertedMap, ok := value.(map[string]any)
	if !ok {
		return "null"
	}

	keys := slices.Sorted(maps.Keys(assertedMap))

	var builder strings.Builder
	builder.WriteString("Map.of(")

	for index, key := range keys {
		if index > 0 {
			builder.WriteString(", ")
		}
		val := assertedMap[key]

		builder.WriteString(fmt.Sprintf("%q, %s", key, formatNestedValue(val)))
	}
	builder.WriteString(")")

	return builder.String()
}

func formatNestedValue(value any) string {
	switch val := value.(type) {
	case string:
		flag := flagset.Flag{
			Type:         flagset.StringType,
			DefaultValue: val,
		}
		return formatDefaultValueForJava(flag)
	case bool:
		flag := flagset.Flag{
			Type:         flagset.BoolType,
			DefaultValue: val,
		}
		return formatDefaultValueForJava(flag)
	case int, int64:
		flag := flagset.Flag{
			Type:         flagset.IntType,
			DefaultValue: val,
		}
		return formatDefaultValueForJava(flag)
	case float64:
		flag := flagset.Flag{
			Type:         flagset.FloatType,
			DefaultValue: val,
		}
		return formatDefaultValueForJava(flag)
	case map[string]any:
		return toMapLiteral(val)
	case []any:
		var sliceBuilder strings.Builder
		sliceBuilder.WriteString("List.of(")
		for index, elem := range val {
			if index > 0 {
				sliceBuilder.WriteString(", ")
			}

			sliceBuilder.WriteString(formatNestedValue(elem))
		}
		sliceBuilder.WriteString(")")
		return sliceBuilder.String()
	default:
		jsonBytes, err := json.Marshal(val)
		if err != nil {
			return "null"
		}
		return fmt.Sprintf("%q", string(jsonBytes))
	}
}

func (g *JavaGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"OpenFeatureType":    openFeatureType,
		"FormatDefaultValue": formatDefaultValueForJava,
		"ToMapLiteral":       toMapLiteral,
	}

	newParams := &generators.Params[any]{
		OutputPath: params.OutputPath,
		Custom:     params.Custom,
	}

	return g.CommonGenerator.GenerateFile(funcs, javaTmpl, newParams, "OpenFeature.java")
}

// NewGenerator creates a generator for Java.
func NewGenerator(fs *flagset.Flagset) *JavaGenerator {
	return &JavaGenerator{
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{}),
	}
}
