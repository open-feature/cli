package golang

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"go/format"
	"maps"
	"slices"
	"strings"
	"text/template"

	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators"
	"golang.org/x/tools/imports"
)

type GolangGenerator struct {
	generators.CommonGenerator
}

type Params struct {
	GoPackage  string
	CLIVersion string
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
		res = append(res, "\"fmt\"")
		res = append(res, "\"github.com/open-feature/go-sdk/openfeature\"")
	}
	slices.Sort(res)
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
		if index > 0 {
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
		return fmt.Sprintf("%t", val)
	case int, int64, float64:
		return fmt.Sprintf("%v", val)
	case map[string]any:
		return toMapLiteral(val)
	case []any:
		var sliceBuilder strings.Builder
		sliceBuilder.WriteString("[]any{")
		for index, elem := range val {
			if index > 0 {
				sliceBuilder.WriteString(", ")
			}

			sliceBuilder.WriteString(formatNestedValue(elem))
		}
		sliceBuilder.WriteString("}")
		return sliceBuilder.String()
	default:
		jsonBytes, err := json.Marshal(val)
		if err != nil {
			return "nil"
		}
		return fmt.Sprintf("%q", string(jsonBytes))
	}
}

// goSchemaType converts an ObjectSchema type string to a Go type string.
func goSchemaType(schema *flagset.ObjectSchema, required bool) string {
	if schema == nil {
		return "any"
	}

	switch schema.Type {
	case "string":
		return "string"
	case "number":
		return "float64"
	case "integer":
		return "int64"
	case "boolean":
		return "bool"
	case "array":
		if schema.Items != nil {
			return "[]" + goSchemaType(schema.Items, true)
		}
		return "[]any"
	case "object":
		if schema.Properties != nil {
			// Inline nested struct
			return generateGoStructBody(schema)
		}
		return "map[string]any"
	default:
		return "any"
	}
}

// generateGoStructBody generates a Go struct type literal (with braces) from a schema.
func generateGoStructBody(schema *flagset.ObjectSchema) string {
	var b strings.Builder
	b.WriteString("struct {\n")

	for _, propName := range generators.SortedPropertyNames(schema.Properties) {
		propSchema := schema.Properties[propName]
		isReq := generators.IsRequired(propName, schema.Required)
		goType := goSchemaType(propSchema, isReq)
		fieldName := generators.ObjectTypeName(propName)

		jsonTag := propName
		if !isReq {
			jsonTag += ",omitempty"
		}

		b.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, goType, jsonTag))
	}

	b.WriteString("}")
	return b.String()
}

// generateGoTypeDef generates a top-level Go type definition for a flag's object schema.
func generateGoTypeDef(flag flagset.Flag) string {
	if !generators.HasSchema(flag) {
		return ""
	}

	typeName := goObjectTypeName(flag.Key)
	body := generateGoStructBody(flag.Schema)

	return fmt.Sprintf("// %s is the typed object for the %q flag.\ntype %s %s\n", typeName, flag.Key, typeName, body)
}

// generateGoHookDef generates the validation hook struct and After method for a typed object flag.
func generateGoHookDef(flag flagset.Flag) string {
	if !generators.HasSchema(flag) {
		return ""
	}

	typeName := generators.ObjectTypeName(flag.Key)
	hookName := fmt.Sprintf("%sHook", strings.ToLower(typeName[:1])+typeName[1:])

	var b strings.Builder
	b.WriteString(fmt.Sprintf("type %s struct {\n\topenfeature.UnimplementedHook\n}\n\n", hookName))
	b.WriteString(fmt.Sprintf("func (h %s) After(ctx context.Context, hookCtx openfeature.HookContext, details openfeature.InterfaceEvaluationDetails, hints openfeature.HookHints) error {\n", hookName))
	b.WriteString(generateGoValidation(flag.Schema, "details.Value", flag.Key))
	b.WriteString("\treturn nil\n}\n")

	return b.String()
}

// goSafeVarName converts a path like "themeCustomization.header" to a safe Go variable name.
func goSafeVarName(path string) string {
	v := strings.ReplaceAll(path, ".", "_")
	v = strings.ReplaceAll(v, "[", "_")
	v = strings.ReplaceAll(v, "]", "")
	return v
}

// generateGoValidation generates Go validation code for a schema at a given path.
func generateGoValidation(schema *flagset.ObjectSchema, accessor string, path string) string {
	var b strings.Builder

	switch schema.Type {
	case "object":
		varName := goSafeVarName(path)
		b.WriteString(fmt.Sprintf("\t%sMap, ok := %s.(map[string]any)\n", varName, accessor))
		b.WriteString(fmt.Sprintf("\tif !ok {\n\t\treturn fmt.Errorf(\"%s: expected object, got %%T\", %s)\n\t}\n", path, accessor))

		for _, req := range schema.Required {
			b.WriteString(fmt.Sprintf("\tif _, exists := %sMap[%q]; !exists {\n", varName, req))
			b.WriteString(fmt.Sprintf("\t\treturn fmt.Errorf(\"%s: missing required property %%q\", %q)\n\t}\n", path, req))
		}

		// Validate each property's type and recurse into nested schemas
		for _, propName := range generators.SortedPropertyNames(schema.Properties) {
			propSchema := schema.Properties[propName]
			propPath := fmt.Sprintf("%s.%s", path, propName)
			propVar := goSafeVarName(propPath)

			b.WriteString(fmt.Sprintf("\tif %s, exists := %sMap[%q]; exists {\n", propVar, varName, propName))
			b.WriteString(generateGoValidation(propSchema, propVar, propPath))
			b.WriteString("\t}\n")
		}

		// Check additionalProperties: false
		if schema.AdditionalProperties != nil && !*schema.AdditionalProperties {
			allowedVar := goSafeVarName(path) + "Allowed"
			b.WriteString(fmt.Sprintf("\t%s := map[string]bool{", allowedVar))
			for i, propName := range generators.SortedPropertyNames(schema.Properties) {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(fmt.Sprintf("%q: true", propName))
			}
			b.WriteString("}\n")
			b.WriteString(fmt.Sprintf("\tfor k := range %sMap {\n", varName))
			b.WriteString(fmt.Sprintf("\t\tif !%s[k] {\n", allowedVar))
			b.WriteString(fmt.Sprintf("\t\t\treturn fmt.Errorf(\"%s: unexpected property %%q\", k)\n", path))
			b.WriteString("\t\t}\n")
			b.WriteString("\t}\n")
		}

	case "array":
		varName := goSafeVarName(path)
		b.WriteString(fmt.Sprintf("\t%sArr, ok := %s.([]any)\n", varName, accessor))
		b.WriteString(fmt.Sprintf("\tif !ok {\n\t\treturn fmt.Errorf(\"%s: expected array, got %%T\", %s)\n\t}\n", path, accessor))

		if schema.Items != nil {
			idxVar := varName + "Idx"
			itemVar := varName + "Item"
			b.WriteString(fmt.Sprintf("\tfor %s, %s := range %sArr {\n", idxVar, itemVar, varName))
			b.WriteString(generateGoArrayItemValidation(schema.Items, itemVar, path, idxVar))
			b.WriteString("\t}\n")
		}

	case "string":
		b.WriteString(fmt.Sprintf("\tif _, ok := %s.(string); !ok {\n", accessor))
		b.WriteString(fmt.Sprintf("\t\treturn fmt.Errorf(\"%s: expected string, got %%T\", %s)\n\t}\n", path, accessor))

	case "number":
		b.WriteString(fmt.Sprintf("\tif _, ok := %s.(float64); !ok {\n", accessor))
		b.WriteString(fmt.Sprintf("\t\treturn fmt.Errorf(\"%s: expected number, got %%T\", %s)\n\t}\n", path, accessor))

	case "integer":
		intVar := goSafeVarName(path) + "Float"
		b.WriteString(fmt.Sprintf("\t%s, ok := %s.(float64)\n", intVar, accessor))
		b.WriteString(fmt.Sprintf("\tif !ok {\n\t\treturn fmt.Errorf(\"%s: expected integer, got %%T\", %s)\n\t}\n", path, accessor))
		b.WriteString(fmt.Sprintf("\tif %s != float64(int64(%s)) {\n", intVar, intVar))
		b.WriteString(fmt.Sprintf("\t\treturn fmt.Errorf(\"%s: expected integer, got float\")\n\t}\n", path))

	case "boolean":
		b.WriteString(fmt.Sprintf("\tif _, ok := %s.(bool); !ok {\n", accessor))
		b.WriteString(fmt.Sprintf("\t\treturn fmt.Errorf(\"%s: expected boolean, got %%T\", %s)\n\t}\n", path, accessor))
	}

	return b.String()
}

// generateGoArrayItemValidation generates validation for array items where the error
// path includes a runtime index variable (e.g., themeCustomization.tags[%d]).
func generateGoArrayItemValidation(schema *flagset.ObjectSchema, itemVar string, arrayPath string, idxVar string) string {
	var b strings.Builder

	switch schema.Type {
	case "string":
		b.WriteString(fmt.Sprintf("\t\tif _, ok := %s.(string); !ok {\n", itemVar))
		b.WriteString(fmt.Sprintf("\t\t\treturn fmt.Errorf(\"%s[%%d]: expected string, got %%T\", %s, %s)\n\t\t}\n", arrayPath, idxVar, itemVar))
	case "number":
		b.WriteString(fmt.Sprintf("\t\tif _, ok := %s.(float64); !ok {\n", itemVar))
		b.WriteString(fmt.Sprintf("\t\t\treturn fmt.Errorf(\"%s[%%d]: expected number, got %%T\", %s, %s)\n\t\t}\n", arrayPath, idxVar, itemVar))
	case "integer":
		floatVar := itemVar + "Float"
		b.WriteString(fmt.Sprintf("\t\t%s, ok := %s.(float64)\n", floatVar, itemVar))
		b.WriteString(fmt.Sprintf("\t\tif !ok {\n\t\t\treturn fmt.Errorf(\"%s[%%d]: expected integer, got %%T\", %s, %s)\n\t\t}\n", arrayPath, idxVar, itemVar))
		b.WriteString(fmt.Sprintf("\t\tif %s != float64(int64(%s)) {\n", floatVar, floatVar))
		b.WriteString(fmt.Sprintf("\t\t\treturn fmt.Errorf(\"%s[%%d]: expected integer, got float\", %s)\n\t\t}\n", arrayPath, idxVar))
	case "boolean":
		b.WriteString(fmt.Sprintf("\t\tif _, ok := %s.(bool); !ok {\n", itemVar))
		b.WriteString(fmt.Sprintf("\t\t\treturn fmt.Errorf(\"%s[%%d]: expected boolean, got %%T\", %s, %s)\n\t\t}\n", arrayPath, idxVar, itemVar))
	case "object":
		b.WriteString(generateGoValidation(schema, itemVar, fmt.Sprintf("%s[item]", arrayPath)))
	case "array":
		b.WriteString(generateGoValidation(schema, itemVar, fmt.Sprintf("%s[item]", arrayPath)))
	}

	return b.String()
}

// goObjectTypeName returns the Go type name with a "Value" suffix to avoid conflicts
// with the generated variable name (which is PascalCase of the flag key).
func goObjectTypeName(flagKey string) string {
	return generators.ObjectTypeName(flagKey) + "Value"
}

// goFlagReturnType returns the Go return type for a flag - typed struct name if schema exists, else the generic type.
func goFlagReturnType(flag flagset.Flag) string {
	if generators.HasSchema(flag) {
		return goObjectTypeName(flag.Key)
	}
	if flag.Type == flagset.ObjectType {
		return "any"
	}
	return typeString(flag.Type)
}

// goHookName returns the hook variable name for a typed object flag.
func goHookName(flag flagset.Flag) string {
	typeName := generators.ObjectTypeName(flag.Key)
	return strings.ToLower(typeName[:1]) + typeName[1:] + "Hook"
}

func (g *GolangGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"SupportImports":           supportImports,
		"OpenFeatureType":          openFeatureType,
		"TypeString":               typeString,
		"ToMapLiteral":             toMapLiteral,
		"HasSchema":                generators.HasSchema,
		"ObjectTypeName":           generators.ObjectTypeName,
		"GoTypeDef":                generateGoTypeDef,
		"GoHookDef":                generateGoHookDef,
		"GoFlagReturnType":         goFlagReturnType,
		"GoHookName":               goHookName,
		"HasObjectFlagsWithSchema": generators.HasObjectFlagsWithSchema,
	}

	newParams := &generators.Params[any]{
		OutputPath:        params.OutputPath,
		TemplatePath:      params.TemplatePath,
		RuntimeValidation: params.RuntimeValidation,
		Custom: Params{
			GoPackage:  params.Custom.GoPackage,
			CLIVersion: params.Custom.CLIVersion,
		},
	}

	filename := params.Custom.GoPackage + "_gen.go"
	return g.GenerateFile(funcs, golangTmpl, newParams, filename)
}

// NewGenerator creates a generator for Go.
func NewGenerator(fs *flagset.Flagset) *GolangGenerator {
	g := &GolangGenerator{
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{}),
	}
	g.Formatter = func(data []byte) ([]byte, error) {
		data, err := format.Source(data)
		if err != nil {
			return nil, fmt.Errorf("failed to format go code: %w", err)
		}

		data, err = imports.Process("", data, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to format imports: %w", err)
		}
		return data, nil
	}

	return g
}
