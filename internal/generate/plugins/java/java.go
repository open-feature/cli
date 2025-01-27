package java

import (
	_ "embed"
	"html/template"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/open-feature/cli/internal/generate"
	"github.com/open-feature/cli/internal/generate/types"
)

// TmplData contains the Java-specific data and the base data for the codegen.
type TmplData struct {
	*types.BaseTmplData
	JavaPackage string
}

type genImpl struct {
	javaPackage string
}

// BaseTmplDataInfo provides the base template data for the codegen.
func (td *TmplData) BaseTmplDataInfo() *types.BaseTmplData {
	return td.BaseTmplData
}

// supportedFlagTypes is the flag types supported by the Java template.
var supportedFlagTypes = map[types.FlagType]bool{
	types.FloatType:  true,
	types.StringType: true,
	types.IntType:    true,
	types.BoolType:   true,
	types.ObjectType: false,
}

func (*genImpl) SupportedFlagTypes() map[types.FlagType]bool {
	return supportedFlagTypes
}

//go:embed java.tmpl
var javaTmpl string

// Java Funcs BEGIN

func flagClass(flagName string) string {
	return strcase.ToCamel(flagName)
}

func topLevelClass(outputPath string) string {
	return strings.TrimSuffix(filepath.Base(outputPath), filepath.Ext(outputPath))
}

func flagInitParam(flagName string) string {
	return strconv.Quote(flagName)
}

func javaType(t types.FlagType) string {
	switch t {
	case types.IntType:
		return "Integer"
	case types.FloatType:
		return "Double"
	case types.BoolType:
		return "Boolean"
	case types.StringType:
		return "String"
	default:
		return ""
	}
}

func supportImports(flags []*types.FlagTmplData) []string {
	var res []string
	if len(flags) > 0 {
		res = append(res, "dev.openfeature.sdk.*")
	}
	sort.Strings(res)
	return res
}

func defaultValueLiteral(flag *types.FlagTmplData) string {
	switch flag.Type {
	case types.StringType:
		return strconv.Quote(flag.DefaultValue)
	default:
		return flag.DefaultValue
	}
}

// Java Funcs END

// Generate generates the Java flag accessors for OpenFeature.
func (g *genImpl) Generate(input types.Input) error {
	funcs := template.FuncMap{
		"FlagClass":           flagClass,
		"TopLevelClass":       topLevelClass,
		"FlagInitParam":       flagInitParam,
		"JavaType":            javaType,
		"SupportImports":      supportImports,
		"DefaultValueLiteral": defaultValueLiteral,
	}
	td := TmplData{
		BaseTmplData: input.BaseData,
		JavaPackage:  g.javaPackage,
	}
	return generate.GenerateFile(funcs, javaTmpl, &td)
}

// Params are parameters for creating a Generator
type Params struct {
	JavaPackage string
}

// NewGenerator creates a generator for Java.
func NewGenerator(params Params) types.Generator {
	return &genImpl{
		javaPackage: params.JavaPackage,
	}
}
