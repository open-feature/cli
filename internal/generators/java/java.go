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
		return "Double" // using Double as per openfeature Java-SDK
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

// javaSchemaType converts an ObjectSchema type string to a Java type string.
// parentName is used to reference named record types for nested objects.
func javaSchemaType(schema *flagset.ObjectSchema, required bool, parentName string) string {
	if schema == nil {
		return "Object"
	}

	prefix := ""
	if !required {
		prefix = "@Nullable "
	}

	switch schema.Type {
	case "string":
		return prefix + "String"
	case "number":
		return prefix + "Double"
	case "integer":
		return prefix + "Integer"
	case "boolean":
		return prefix + "Boolean"
	case "array":
		if schema.Items != nil {
			itemType := javaSchemaType(schema.Items, true, parentName+"Item")
			return prefix + "List<" + itemType + ">"
		}
		return prefix + "List<Object>"
	case "object":
		if schema.Properties != nil {
			return prefix + parentName
		}
		return prefix + "Map<String, Object>"
	default:
		return prefix + "Object"
	}
}

// collectJavaRecordDefs recursively collects all record definitions needed for a schema,
// emitting nested records before their parents so dependencies are satisfied.
func collectJavaRecordDefs(typeName string, schema *flagset.ObjectSchema) string {
	var b strings.Builder

	// First, recurse into nested objects to emit their records
	for _, propName := range generators.SortedPropertyNames(schema.Properties) {
		propSchema := schema.Properties[propName]
		if propSchema != nil && propSchema.Type == "object" && propSchema.Properties != nil {
			nestedName := typeName + generators.ObjectTypeName(propName)
			b.WriteString(collectJavaRecordDefs(nestedName, propSchema))
			b.WriteString("\n\n")
		}
		if propSchema != nil && propSchema.Type == "array" && propSchema.Items != nil &&
			propSchema.Items.Type == "object" && propSchema.Items.Properties != nil {
			nestedName := typeName + generators.ObjectTypeName(propName) + "Item"
			b.WriteString(collectJavaRecordDefs(nestedName, propSchema.Items))
			b.WriteString("\n\n")
		}
	}

	// Emit this record
	b.WriteString(fmt.Sprintf("    public record %s(\n", typeName))
	propNames := generators.SortedPropertyNames(schema.Properties)
	for i, propName := range propNames {
		propSchema := schema.Properties[propName]
		isReq := generators.IsRequired(propName, schema.Required)
		nestedName := typeName + generators.ObjectTypeName(propName)
		javaType := javaSchemaType(propSchema, isReq, nestedName)
		fieldName := generators.ObjectTypeName(propName)
		fieldName = strings.ToLower(fieldName[:1]) + fieldName[1:]

		b.WriteString(fmt.Sprintf("        %s %s", javaType, fieldName))
		if i < len(propNames)-1 {
			b.WriteString(",")
		}
		b.WriteString("\n")
	}
	b.WriteString("    ) {}")

	return b.String()
}

// generateJavaRecordDef generates Java record definitions for a flag's object schema,
// including any nested object types.
func generateJavaRecordDef(flag flagset.Flag) string {
	if !generators.HasSchema(flag) {
		return ""
	}

	typeName := generators.ObjectTypeName(flag.Key)
	return collectJavaRecordDefs(typeName, flag.Schema)
}

// javaSafeVarName converts a path to a safe Java variable name.
func javaSafeVarName(path string) string {
	v := strings.ReplaceAll(path, ".", "_")
	v = strings.ReplaceAll(v, "[", "_")
	v = strings.ReplaceAll(v, "]", "")
	return v
}

// generateJavaHookDef generates a Hook inner class that validates the object shape for a typed flag.
func generateJavaHookDef(flag flagset.Flag) string {
	if !generators.HasSchema(flag) {
		return ""
	}

	typeName := generators.ObjectTypeName(flag.Key)
	hookName := typeName + "Hook"

	var b strings.Builder
	b.WriteString(fmt.Sprintf("    static class %s implements Hook<Object> {\n", hookName))
	b.WriteString("        @Override\n")
	b.WriteString("        public void after(HookContext<Object> ctx, FlagEvaluationDetails<Object> details, Map<String, Object> hints) {\n")
	b.WriteString("            Object value = details.getValue();\n")
	b.WriteString(generateJavaValidation(flag.Schema, "value", flag.Key, "            "))
	b.WriteString("        }\n")
	b.WriteString("    }")
	return b.String()
}

// generateJavaValidation generates Java validation code for a schema at a given path.
func generateJavaValidation(schema *flagset.ObjectSchema, accessor string, path string, indent string) string {
	var b strings.Builder

	switch schema.Type {
	case "object":
		varName := javaSafeVarName(path) + "Map"
		b.WriteString(fmt.Sprintf("%sif (!(%s instanceof Map)) {\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s    throw new IllegalArgumentException(\"%s: expected object, got \" + (%s == null ? \"null\" : %s.getClass().getName()));\n", indent, path, accessor, accessor))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
		b.WriteString(fmt.Sprintf("%s@SuppressWarnings(\"unchecked\")\n", indent))
		b.WriteString(fmt.Sprintf("%sMap<String, Object> %s = (Map<String, Object>) %s;\n", indent, varName, accessor))

		for _, req := range schema.Required {
			b.WriteString(fmt.Sprintf("%sif (!%s.containsKey(\"%s\")) {\n", indent, varName, req))
			b.WriteString(fmt.Sprintf("%s    throw new IllegalArgumentException(\"%s: missing required property '%s'\");\n", indent, path, req))
			b.WriteString(fmt.Sprintf("%s}\n", indent))
		}

		// Validate each property's type and recurse into nested schemas
		for _, propName := range generators.SortedPropertyNames(schema.Properties) {
			propSchema := schema.Properties[propName]
			propPath := fmt.Sprintf("%s.%s", path, propName)
			propVar := javaSafeVarName(propPath)

			b.WriteString(fmt.Sprintf("%sif (%s.containsKey(\"%s\")) {\n", indent, varName, propName))
			b.WriteString(fmt.Sprintf("%s    Object %s = %s.get(\"%s\");\n", indent, propVar, varName, propName))
			b.WriteString(generateJavaValidation(propSchema, propVar, propPath, indent+"    "))
			b.WriteString(fmt.Sprintf("%s}\n", indent))
		}

		// Check additionalProperties: false
		if schema.AdditionalProperties != nil && !*schema.AdditionalProperties {
			allowedVar := javaSafeVarName(path) + "Allowed"
			b.WriteString(fmt.Sprintf("%sjava.util.Set<String> %s = java.util.Set.of(", indent, allowedVar))
			for i, propName := range generators.SortedPropertyNames(schema.Properties) {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(fmt.Sprintf("\"%s\"", propName))
			}
			b.WriteString(");\n")
			b.WriteString(fmt.Sprintf("%sfor (String key : %s.keySet()) {\n", indent, varName))
			b.WriteString(fmt.Sprintf("%s    if (!%s.contains(key)) {\n", indent, allowedVar))
			b.WriteString(fmt.Sprintf("%s        throw new IllegalArgumentException(\"%s: unexpected property '\" + key + \"'\");\n", indent, path))
			b.WriteString(fmt.Sprintf("%s    }\n", indent))
			b.WriteString(fmt.Sprintf("%s}\n", indent))
		}

	case "array":
		varName := javaSafeVarName(path) + "List"
		b.WriteString(fmt.Sprintf("%sif (!(%s instanceof java.util.List)) {\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s    throw new IllegalArgumentException(\"%s: expected array, got \" + (%s == null ? \"null\" : %s.getClass().getName()));\n", indent, path, accessor, accessor))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
		b.WriteString(fmt.Sprintf("%sjava.util.List<?> %s = (java.util.List<?>) %s;\n", indent, varName, accessor))

		if schema.Items != nil {
			idxVar := javaSafeVarName(path) + "Idx"
			b.WriteString(fmt.Sprintf("%sfor (int %s = 0; %s < %s.size(); %s++) {\n", indent, idxVar, idxVar, varName, idxVar))
			itemAccessor := fmt.Sprintf("%s.get(%s)", varName, idxVar)
			b.WriteString(generateJavaArrayItemValidation(schema.Items, itemAccessor, path, idxVar, indent+"    "))
			b.WriteString(fmt.Sprintf("%s}\n", indent))
		}

	case "string":
		b.WriteString(fmt.Sprintf("%sif (!(%s instanceof String)) {\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s    throw new IllegalArgumentException(\"%s: expected String, got \" + (%s == null ? \"null\" : %s.getClass().getName()));\n", indent, path, accessor, accessor))
		b.WriteString(fmt.Sprintf("%s}\n", indent))

	case "number":
		b.WriteString(fmt.Sprintf("%sif (!(%s instanceof Number)) {\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s    throw new IllegalArgumentException(\"%s: expected Number, got \" + (%s == null ? \"null\" : %s.getClass().getName()));\n", indent, path, accessor, accessor))
		b.WriteString(fmt.Sprintf("%s}\n", indent))

	case "integer":
		b.WriteString(fmt.Sprintf("%sif (!(%s instanceof Integer) && !(%s instanceof Long)) {\n", indent, accessor, accessor))
		b.WriteString(fmt.Sprintf("%s    throw new IllegalArgumentException(\"%s: expected Integer, got \" + (%s == null ? \"null\" : %s.getClass().getName()));\n", indent, path, accessor, accessor))
		b.WriteString(fmt.Sprintf("%s}\n", indent))

	case "boolean":
		b.WriteString(fmt.Sprintf("%sif (!(%s instanceof Boolean)) {\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s    throw new IllegalArgumentException(\"%s: expected Boolean, got \" + (%s == null ? \"null\" : %s.getClass().getName()));\n", indent, path, accessor, accessor))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
	}

	return b.String()
}

// generateJavaArrayItemValidation generates validation for array items with runtime index in error paths.
func generateJavaArrayItemValidation(schema *flagset.ObjectSchema, itemAccessor string, arrayPath string, idxVar string, indent string) string {
	var b strings.Builder

	switch schema.Type {
	case "string":
		b.WriteString(fmt.Sprintf("%sif (!(%s instanceof String)) {\n", indent, itemAccessor))
		b.WriteString(fmt.Sprintf("%s    throw new IllegalArgumentException(\"%s[\" + %s + \"]: expected String\");\n", indent, arrayPath, idxVar))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
	case "number":
		b.WriteString(fmt.Sprintf("%sif (!(%s instanceof Number)) {\n", indent, itemAccessor))
		b.WriteString(fmt.Sprintf("%s    throw new IllegalArgumentException(\"%s[\" + %s + \"]: expected Number\");\n", indent, arrayPath, idxVar))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
	case "integer":
		b.WriteString(fmt.Sprintf("%sif (!(%s instanceof Integer) && !(%s instanceof Long)) {\n", indent, itemAccessor, itemAccessor))
		b.WriteString(fmt.Sprintf("%s    throw new IllegalArgumentException(\"%s[\" + %s + \"]: expected Integer\");\n", indent, arrayPath, idxVar))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
	case "boolean":
		b.WriteString(fmt.Sprintf("%sif (!(%s instanceof Boolean)) {\n", indent, itemAccessor))
		b.WriteString(fmt.Sprintf("%s    throw new IllegalArgumentException(\"%s[\" + %s + \"]: expected Boolean\");\n", indent, arrayPath, idxVar))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
	case "object", "array":
		b.WriteString(generateJavaValidation(schema, itemAccessor, fmt.Sprintf("%s[item]", arrayPath), indent))
	}

	return b.String()
}

// javaFlagReturnType returns the return type for a flag - the record name if schema exists, "Object" otherwise.
func javaFlagReturnType(flag flagset.Flag) string {
	if generators.HasSchema(flag) {
		return generators.ObjectTypeName(flag.Key)
	}
	return "Object"
}

// javaHookName returns the hook class name for a typed object flag.
func javaHookName(flag flagset.Flag) string {
	return generators.ObjectTypeName(flag.Key) + "Hook"
}

func (g *JavaGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"OpenFeatureType":          openFeatureType,
		"FormatDefaultValue":       formatDefaultValueForJava,
		"ToMapLiteral":             toMapLiteral,
		"HasSchema":                generators.HasSchema,
		"JavaRecordDef":            generateJavaRecordDef,
		"JavaHookDef":              generateJavaHookDef,
		"JavaFlagReturnType":       javaFlagReturnType,
		"JavaHookName":             javaHookName,
		"HasObjectFlagsWithSchema": generators.HasObjectFlagsWithSchema,
	}

	newParams := &generators.Params[any]{
		OutputPath:        params.OutputPath,
		TemplatePath:      params.TemplatePath,
		RuntimeValidation: params.RuntimeValidation,
		Custom:            params.Custom,
	}

	return g.GenerateFile(funcs, javaTmpl, newParams, "OpenFeature.java")
}

// NewGenerator creates a generator for Java.
func NewGenerator(fs *flagset.Flagset) *JavaGenerator {
	return &JavaGenerator{
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{}),
	}
}
