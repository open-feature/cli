package java

import (
	_ "embed"
	"fmt"
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

func (g *JavaGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"OpenFeatureType":    openFeatureType,
		"FormatDefaultValue": formatDefaultValueForJava,
	}

	newParams := &generators.Params[any]{
		OutputPath: params.OutputPath,
		Custom:     params.Custom,
	}

	return g.GenerateFile(funcs, javaTmpl, newParams, "OpenFeature.java")
}

// NewGenerator creates a generator for Java.
func NewGenerator(fs *flagset.Flagset) *JavaGenerator {
	return &JavaGenerator{
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{
			flagset.ObjectType: true,
		}),
	}
}
