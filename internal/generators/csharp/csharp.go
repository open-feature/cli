package csharp

import (
	_ "embed"
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
	builder.WriteString("new Value(Structure.Builder()")

	for _, key := range keys {
		val := assertedMap[key]

		builder.WriteString(fmt.Sprintf(".Set(%q, %s)", key, formatNestedValue(val)))
	}
	builder.WriteString(".Build())")

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
		sliceBuilder.WriteString("new Value(new List<Value>{")
		for index, elem := range val {
			if index > 0 {
				sliceBuilder.WriteString(", ")
			}

			sliceBuilder.WriteString(formatNestedValue(elem))
		}
		sliceBuilder.WriteString("})")
		return sliceBuilder.String()
	default:
		return fmt.Sprintf("new Value(%s)", val)
	}
}

// csharpSchemaType converts an ObjectSchema type string to a C# type string.
// parentName is used to reference named record types for nested objects.
func csharpSchemaType(schema *flagset.ObjectSchema, required bool, parentName string) string {
	if schema == nil {
		return "object"
	}

	baseType := ""
	switch schema.Type {
	case "string":
		baseType = "string"
	case "number":
		baseType = "double"
	case "integer":
		baseType = "int"
	case "boolean":
		baseType = "bool"
	case "array":
		if schema.Items != nil {
			baseType = "List<" + csharpSchemaType(schema.Items, true, parentName+"Item") + ">"
		} else {
			baseType = "List<object>"
		}
	case "object":
		if schema.Properties != nil {
			baseType = parentName
		} else {
			baseType = "Dictionary<string, object>"
		}
	default:
		baseType = "object"
	}

	if !required {
		return baseType + "?"
	}
	return baseType
}

// collectCSharpRecordDefs recursively collects all record definitions needed for a schema,
// emitting nested records before their parents so dependencies are satisfied.
func collectCSharpRecordDefs(typeName string, schema *flagset.ObjectSchema, flagKey string) string {
	var b strings.Builder

	// First, recurse into nested objects to emit their records
	for _, propName := range generators.SortedPropertyNames(schema.Properties) {
		propSchema := schema.Properties[propName]
		if propSchema != nil && propSchema.Type == "object" && propSchema.Properties != nil {
			nestedName := typeName + generators.ObjectTypeName(propName)
			b.WriteString(collectCSharpRecordDefs(nestedName, propSchema, flagKey))
			b.WriteString("\n")
		}
		if propSchema != nil && propSchema.Type == "array" && propSchema.Items != nil &&
			propSchema.Items.Type == "object" && propSchema.Items.Properties != nil {
			nestedName := typeName + generators.ObjectTypeName(propName) + "Item"
			b.WriteString(collectCSharpRecordDefs(nestedName, propSchema.Items, flagKey))
			b.WriteString("\n")
		}
	}

	// Emit this record
	b.WriteString(fmt.Sprintf("    /// <summary>\n    /// Typed object for the %q flag.\n    /// </summary>\n", flagKey))
	b.WriteString(fmt.Sprintf("    public record %s(\n", typeName))
	propNames := generators.SortedPropertyNames(schema.Properties)
	for i, propName := range propNames {
		propSchema := schema.Properties[propName]
		isReq := generators.IsRequired(propName, schema.Required)
		nestedName := typeName + generators.ObjectTypeName(propName)
		csType := csharpSchemaType(propSchema, isReq, nestedName)
		fieldName := generators.ObjectTypeName(propName)

		b.WriteString(fmt.Sprintf("        %s %s", csType, fieldName))
		if i < len(propNames)-1 {
			b.WriteString(",\n")
		} else {
			b.WriteString("\n")
		}
	}
	b.WriteString("    );\n")

	return b.String()
}

// generateCSharpRecordDef generates C# record definitions for a flag's object schema,
// including any nested object types.
func generateCSharpRecordDef(flag flagset.Flag) string {
	if !generators.HasSchema(flag) {
		return ""
	}

	typeName := generators.ObjectTypeName(flag.Key)
	return collectCSharpRecordDefs(typeName, flag.Schema, flag.Key)
}

// csSafeVarName converts a path to a safe C# variable name.
func csSafeVarName(path string) string {
	v := strings.ReplaceAll(path, ".", "_")
	v = strings.ReplaceAll(v, "[", "_")
	v = strings.ReplaceAll(v, "]", "")
	return v
}

// generateCSharpHookDef generates a C# Hook class for runtime validation of a typed object flag.
func generateCSharpHookDef(flag flagset.Flag) string {
	if !generators.HasSchema(flag) {
		return ""
	}

	typeName := generators.ObjectTypeName(flag.Key)
	hookName := csharpHookName(flag)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("    /// <summary>\n    /// Validation hook for the %q flag.\n    /// </summary>\n", flag.Key))
	b.WriteString(fmt.Sprintf("    public class %s : Hook\n", hookName))
	b.WriteString("    {\n")
	b.WriteString(fmt.Sprintf("        public override ValueTask<EvaluationContext> BeforeAsync<T>(HookContext<T> context, IReadOnlyDictionary<string, object>? hints = null)\n"))
	b.WriteString("        {\n")
	b.WriteString("            return new ValueTask<EvaluationContext>(EvaluationContext.Empty);\n")
	b.WriteString("        }\n\n")
	b.WriteString(fmt.Sprintf("        public override ValueTask AfterAsync<T>(HookContext<T> context, FlagEvaluationDetails<T> details, IReadOnlyDictionary<string, object>? hints = null)\n"))
	b.WriteString("        {\n")
	b.WriteString("            if (details.Value is not Value value)\n")
	b.WriteString("            {\n")
	b.WriteString(fmt.Sprintf("                throw new InvalidOperationException(\"%s: expected Value type\");\n", flag.Key))
	b.WriteString("            }\n")
	b.WriteString(generateCSharpValidation(flag.Schema, "value", flag.Key, "            "))

	_ = typeName // typeName used in the doc comment above
	b.WriteString("            return new ValueTask();\n")
	b.WriteString("        }\n")
	b.WriteString("    }\n")
	return b.String()
}

// generateCSharpValidation generates C# validation code for a schema at a given path.
// For objects, the accessor is a Value and we use .AsStructure. For primitives, we use typed accessors.
func generateCSharpValidation(schema *flagset.ObjectSchema, accessor string, path string, indent string) string {
	var b strings.Builder

	switch schema.Type {
	case "object":
		varName := csSafeVarName(path) + "Struct"
		b.WriteString(fmt.Sprintf("%svar %s = %s.AsStructure;\n", indent, varName, accessor))
		b.WriteString(fmt.Sprintf("%sif (%s == null)\n", indent, varName))
		b.WriteString(fmt.Sprintf("%s{\n", indent))
		b.WriteString(fmt.Sprintf("%s    throw new InvalidOperationException(\"%s: expected object structure\");\n", indent, path))
		b.WriteString(fmt.Sprintf("%s}\n", indent))

		for _, req := range schema.Required {
			b.WriteString(fmt.Sprintf("%sif (!%s.ContainsKey(%q))\n", indent, varName, req))
			b.WriteString(fmt.Sprintf("%s{\n", indent))
			b.WriteString(fmt.Sprintf("%s    throw new InvalidOperationException(\"%s: missing required property '%s'\");\n", indent, path, req))
			b.WriteString(fmt.Sprintf("%s}\n", indent))
		}

		// Validate each property's type and recurse into nested schemas
		for _, propName := range generators.SortedPropertyNames(schema.Properties) {
			propSchema := schema.Properties[propName]
			propPath := fmt.Sprintf("%s.%s", path, propName)
			propVar := csSafeVarName(propPath) + "Val"

			b.WriteString(fmt.Sprintf("%sif (%s.ContainsKey(%q))\n", indent, varName, propName))
			b.WriteString(fmt.Sprintf("%s{\n", indent))
			b.WriteString(fmt.Sprintf("%s    var %s = %s[%q];\n", indent, propVar, varName, propName))
			b.WriteString(generateCSharpValidation(propSchema, propVar, propPath, indent+"    "))
			b.WriteString(fmt.Sprintf("%s}\n", indent))
		}

		// Check additionalProperties: false
		if schema.AdditionalProperties != nil && !*schema.AdditionalProperties {
			allowedVar := csSafeVarName(path) + "Allowed"
			b.WriteString(fmt.Sprintf("%svar %s = new HashSet<string> { ", indent, allowedVar))
			for i, propName := range generators.SortedPropertyNames(schema.Properties) {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(fmt.Sprintf("%q", propName))
			}
			b.WriteString(" };\n")
			b.WriteString(fmt.Sprintf("%sforeach (var key in %s.Keys)\n", indent, varName))
			b.WriteString(fmt.Sprintf("%s{\n", indent))
			b.WriteString(fmt.Sprintf("%s    if (!%s.Contains(key))\n", indent, allowedVar))
			b.WriteString(fmt.Sprintf("%s    {\n", indent))
			b.WriteString(fmt.Sprintf("%s        throw new InvalidOperationException($\"%s: unexpected property '{key}'\");\n", indent, path))
			b.WriteString(fmt.Sprintf("%s    }\n", indent))
			b.WriteString(fmt.Sprintf("%s}\n", indent))
		}

	case "array":
		varName := csSafeVarName(path) + "List"
		b.WriteString(fmt.Sprintf("%svar %s = %s.AsList;\n", indent, varName, accessor))
		b.WriteString(fmt.Sprintf("%sif (%s == null)\n", indent, varName))
		b.WriteString(fmt.Sprintf("%s{\n", indent))
		b.WriteString(fmt.Sprintf("%s    throw new InvalidOperationException(\"%s: expected array\");\n", indent, path))
		b.WriteString(fmt.Sprintf("%s}\n", indent))

		if schema.Items != nil {
			idxVar := csSafeVarName(path) + "Idx"
			b.WriteString(fmt.Sprintf("%sfor (var %s = 0; %s < %s.Count; %s++)\n", indent, idxVar, idxVar, varName, idxVar))
			b.WriteString(fmt.Sprintf("%s{\n", indent))
			itemAccessor := fmt.Sprintf("%s[%s]", varName, idxVar)
			b.WriteString(generateCSharpArrayItemValidation(schema.Items, itemAccessor, path, idxVar, indent+"    "))
			b.WriteString(fmt.Sprintf("%s}\n", indent))
		}

	case "string":
		b.WriteString(fmt.Sprintf("%sif (%s.AsString == null)\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s{\n", indent))
		b.WriteString(fmt.Sprintf("%s    throw new InvalidOperationException(\"%s: expected string\");\n", indent, path))
		b.WriteString(fmt.Sprintf("%s}\n", indent))

	case "number":
		b.WriteString(fmt.Sprintf("%sif (%s.AsDouble == null)\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s{\n", indent))
		b.WriteString(fmt.Sprintf("%s    throw new InvalidOperationException(\"%s: expected number\");\n", indent, path))
		b.WriteString(fmt.Sprintf("%s}\n", indent))

	case "integer":
		b.WriteString(fmt.Sprintf("%sif (%s.AsInteger == null)\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s{\n", indent))
		b.WriteString(fmt.Sprintf("%s    throw new InvalidOperationException(\"%s: expected integer\");\n", indent, path))
		b.WriteString(fmt.Sprintf("%s}\n", indent))

	case "boolean":
		b.WriteString(fmt.Sprintf("%sif (%s.AsBoolean == null)\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s{\n", indent))
		b.WriteString(fmt.Sprintf("%s    throw new InvalidOperationException(\"%s: expected boolean\");\n", indent, path))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
	}

	return b.String()
}

// generateCSharpArrayItemValidation generates validation for array items with runtime index in error paths.
func generateCSharpArrayItemValidation(schema *flagset.ObjectSchema, itemAccessor string, arrayPath string, idxVar string, indent string) string {
	var b strings.Builder

	switch schema.Type {
	case "string":
		b.WriteString(fmt.Sprintf("%sif (%s.AsString == null)\n", indent, itemAccessor))
		b.WriteString(fmt.Sprintf("%s{\n", indent))
		b.WriteString(fmt.Sprintf("%s    throw new InvalidOperationException($\"%s[{%s}]: expected string\");\n", indent, arrayPath, idxVar))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
	case "number":
		b.WriteString(fmt.Sprintf("%sif (%s.AsDouble == null)\n", indent, itemAccessor))
		b.WriteString(fmt.Sprintf("%s{\n", indent))
		b.WriteString(fmt.Sprintf("%s    throw new InvalidOperationException($\"%s[{%s}]: expected number\");\n", indent, arrayPath, idxVar))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
	case "integer":
		b.WriteString(fmt.Sprintf("%sif (%s.AsInteger == null)\n", indent, itemAccessor))
		b.WriteString(fmt.Sprintf("%s{\n", indent))
		b.WriteString(fmt.Sprintf("%s    throw new InvalidOperationException($\"%s[{%s}]: expected integer\");\n", indent, arrayPath, idxVar))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
	case "boolean":
		b.WriteString(fmt.Sprintf("%sif (%s.AsBoolean == null)\n", indent, itemAccessor))
		b.WriteString(fmt.Sprintf("%s{\n", indent))
		b.WriteString(fmt.Sprintf("%s    throw new InvalidOperationException($\"%s[{%s}]: expected boolean\");\n", indent, arrayPath, idxVar))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
	case "object", "array":
		b.WriteString(generateCSharpValidation(schema, itemAccessor, fmt.Sprintf("%s[item]", arrayPath), indent))
	}

	return b.String()
}

// csharpFlagReturnType returns the C# return type for a flag.
// For schema-typed object flags it returns the record type name; otherwise the standard type.
func csharpFlagReturnType(flag flagset.Flag) string {
	if generators.HasSchema(flag) {
		return generators.ObjectTypeName(flag.Key)
	}
	if flag.Type == flagset.ObjectType {
		return "Value"
	}
	return openFeatureType(flag.Type)
}

// csharpHookName returns the hook class name for a typed object flag.
func csharpHookName(flag flagset.Flag) string {
	return generators.ObjectTypeName(flag.Key) + "Hook"
}

func (g *CsharpGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"OpenFeatureType":          openFeatureType,
		"FormatDefaultValue":       formatDefaultValue,
		"ToCSharpDict":             toCSharpDict,
		"HasSchema":                generators.HasSchema,
		"CSharpRecordDef":          generateCSharpRecordDef,
		"CSharpHookDef":            generateCSharpHookDef,
		"CSharpFlagReturnType":     csharpFlagReturnType,
		"CSharpHookName":           csharpHookName,
		"HasObjectFlagsWithSchema": generators.HasObjectFlagsWithSchema,
	}

	newParams := &generators.Params[any]{
		OutputPath:        params.OutputPath,
		TemplatePath:      params.TemplatePath,
		RuntimeValidation: params.RuntimeValidation,
		Custom:            params.Custom,
	}

	return g.GenerateFile(funcs, csharpTmpl, newParams, "OpenFeature.g.cs")
}

// NewGenerator creates a generator for C#.
func NewGenerator(fs *flagset.Flagset) *CsharpGenerator {
	return &CsharpGenerator{
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{}),
	}
}
