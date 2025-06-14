package golang

import (
	_ "embed"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"
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
	case flagset.ObjectType:
		return "Object"
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
	case flagset.ObjectType:
		return "map[string]any"
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

func toMapLiteral(value any) string {
	assertedMap, ok := value.(map[string]any)
	if !ok {
		return "nil"
	}

	// To have a determined order of the object for comparison
	keys := slices.Sorted(maps.Keys(assertedMap))

	var builder strings.Builder
	builder.WriteString("map[string]any{")

	for index, key := range keys {
		if index != 0 {
			builder.WriteString(",")
		}
		val := assertedMap[key]

		builder.WriteString(fmt.Sprintf(`"%s": %q`, key, val))
	}

	builder.WriteString("}")
	return builder.String()
}

func (g *GolangGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"SupportImports":  supportImports,
		"OpenFeatureType": openFeatureType,
		"TypeString":      typeString,
		"ToMapLiteral":    toMapLiteral,
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
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{}),
	}
}
