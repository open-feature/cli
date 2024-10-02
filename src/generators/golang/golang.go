package golang

import (
	_ "embed"
	"sort"
	"strconv"
	"text/template"

	generator "codegen/src/generators"

	"github.com/iancoleman/strcase"
)

// TmplData contains the Golang-specific data and the base data for the codegen.
type TmplData struct {
	*generator.BaseTmplData
	GoPackage string
}

type genImpl struct {
	file      string
	goPackage string
}

// BaseTmplDataInfo provides the base template data for the codegen.
func (td *TmplData) BaseTmplDataInfo() *generator.BaseTmplData {
	return td.BaseTmplData
}

// supportedFlagTypes is the flag types supported by the Go template.
var supportedFlagTypes = map[generator.FlagType]bool{
	generator.FloatType:  true,
	generator.StringType: true,
	generator.IntType:    true,
	generator.BoolType:   true,
	generator.ObjectType: false,
}

func (*genImpl) SupportedFlagTypes() map[generator.FlagType]bool {
	return supportedFlagTypes
}

//go:embed golang.tmpl
var golangTmpl string

// Go Funcs BEGIN

func flagVarName(flagName string) string {
	return strcase.ToCamel(flagName)
}

func flagInitParam(flagName string) string {
	return strconv.Quote(flagName)
}

// flagVarType returns the Go type for a flag's proto definition.
func providerType(t generator.FlagType) string {
	switch t {
	case generator.IntType:
		return "IntProvider"
	case generator.FloatType:
		return "FloatProvider"
	case generator.BoolType:
		return "BooleanProvider"
	case generator.StringType:
		return "StringProvider"
	default:
		return ""
	}
}

func flagAccessFunc(t generator.FlagType) string {
	switch t {
	case generator.IntType:
		return "IntValue"
	case generator.FloatType:
		return "FloatValue"
	case generator.BoolType:
		return "BooleanValue"
	case generator.StringType:
		return "StringValue"
	default:
		return ""
	}
}

func supportImports(flags []*generator.FlagTmplData) []string {
	var res []string
	if len(flags) > 0 {
		res = append(res, "\"context\"")
		res = append(res, "\"github.com/open-feature/go-sdk/openfeature\"")
		res = append(res, "\"codegen/src/providers\"")
	}
	sort.Strings(res)
	return res
}

func defaultValueLiteral(flag *generator.FlagTmplData) string {
	switch flag.Type {
	case generator.StringType:
		return strconv.Quote(flag.DefaultValue)
	default:
		return flag.DefaultValue
	}
}

func typeString(flagType generator.FlagType) string {
	switch flagType {
	case generator.StringType:
		return "string"
	case generator.IntType:
		return "int64"
	case generator.BoolType:
		return "bool"
	case generator.FloatType:
		return "float64"
	default:
		return ""
	}
}

// Go Funcs END

// Generate generates the Go flag accessors for OpenFeature.
func (g *genImpl) Generate(input generator.Input) error {
	funcs := template.FuncMap{
		"FlagVarName":         flagVarName,
		"FlagInitParam":       flagInitParam,
		"ProviderType":        providerType,
		"FlagAccessFunc":      flagAccessFunc,
		"SupportImports":      supportImports,
		"DefaultValueLiteral": defaultValueLiteral,
		"TypeString":          typeString,
	}
	td := TmplData{
		BaseTmplData: input.BaseData,
		GoPackage:    g.goPackage,
	}
	return generator.GenerateFile(funcs, g.file, golangTmpl, &td)
}

// Params are parameters for creating a Generator
type Params struct {
	File      string
	GoPackage string
}

// NewGenerator creates a generator for Go.
func NewGenerator(params Params) generator.Generator {
	return &genImpl{
		file:      params.File,
		goPackage: params.GoPackage,
	}
}
