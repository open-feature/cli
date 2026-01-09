package nodejs

import (
	_ "embed"
	"encoding/json"
	"text/template"

	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators"
)

type NodejsGenerator struct {
	generators.CommonGenerator
}

type Params struct{}

//go:embed nodejs.tmpl
var nodejsTmpl string

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
	case flagset.ObjectType:
		return "object"
	default:
		return ""
	}
}

func toJSONString(value any) string {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (g *NodejsGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"OpenFeatureType": openFeatureType,
		"ToJSONString":    toJSONString,
	}

	newParams := &generators.Params[any]{
		OutputPath: params.OutputPath,
		Custom:     Params{},
	}

	return g.GenerateFile(funcs, nodejsTmpl, newParams, "openfeature.ts")
}

// NewGenerator creates a generator for NodeJS.
func NewGenerator(fs *flagset.Flagset) *NodejsGenerator {
	return &NodejsGenerator{
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{}),
	}
}
