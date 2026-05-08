package python

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

type PythonGenerator struct {
	generators.CommonGenerator
}

type Params struct{}

//go:embed python.tmpl
var pythonTmpl string

func openFeatureType(t flagset.FlagType) string {
	switch t {
	case flagset.IntType:
		return "int"
	case flagset.FloatType:
		return "float"
	case flagset.BoolType:
		return "bool"
	case flagset.StringType:
		return "str"
	default:
		return "object"
	}
}

func methodType(flagType flagset.FlagType) string {
	switch flagType {
	case flagset.StringType:
		return "string"
	case flagset.IntType:
		return "integer"
	case flagset.BoolType:
		return "boolean"
	case flagset.FloatType:
		return "float"
	case flagset.ObjectType:
		return "object"
	default:
		panic("unsupported flag type")
	}
}

func typedGetMethodSync(flagType flagset.FlagType) string {
	return "get_" + methodType(flagType) + "_value"
}

func typedGetMethodAsync(flagType flagset.FlagType) string {
	return "get_" + methodType(flagType) + "_value_async"
}

func typedDetailsMethodSync(flagType flagset.FlagType) string {
	return "get_" + methodType(flagType) + "_details"
}

func typedDetailsMethodAsync(flagType flagset.FlagType) string {
	return "get_" + methodType(flagType) + "_details_async"
}

func pythonBoolLiteral(value any) any {
	if v, ok := value.(bool); ok {
		if v {
			return "True"
		}
		return "False"
	}
	return value
}

func toPythonDict(value any) string {
	assertedMap, ok := value.(map[string]any)
	if !ok {
		return "None"
	}

	// To have a determined order of the object for comparison
	keys := slices.Sorted(maps.Keys(assertedMap))

	var builder strings.Builder
	builder.WriteString("{")

	for index, key := range keys {
		if index != 0 {
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
		return pythonBoolLiteral(val).(string)
	case int, int64, float64:
		return fmt.Sprintf("%v", val)
	case map[string]any:
		return toPythonDict(val)
	case []any:
		var sliceBuilder strings.Builder
		sliceBuilder.WriteString("[")
		for index, elem := range val {
			if index > 0 {
				sliceBuilder.WriteString(", ")
			}

			sliceBuilder.WriteString(formatNestedValue(elem))
		}
		sliceBuilder.WriteString("]")
		return sliceBuilder.String()
	default:
		jsonBytes, err := json.Marshal(val)
		if err != nil {
			return "None"
		}
		return strings.ReplaceAll(string(jsonBytes), "null", "None")
	}
}

// pythonSchemaType converts an ObjectSchema type to a Python type annotation string.
// parentName is used to generate named TypedDict classes for nested objects.
func pythonSchemaType(schema *flagset.ObjectSchema, parentName string) string {
	if schema == nil {
		return "dict[str, Any]"
	}

	switch schema.Type {
	case "string":
		return "str"
	case "number":
		return "float"
	case "integer":
		return "int"
	case "boolean":
		return "bool"
	case "array":
		if schema.Items != nil {
			return "list[" + pythonSchemaType(schema.Items, parentName+"Item") + "]"
		}
		return "list[Any]"
	case "object":
		if schema.Properties != nil {
			return parentName
		}
		return "dict[str, Any]"
	default:
		return "Any"
	}
}

// collectPythonTypedDicts recursively collects all TypedDict definitions needed for a schema,
// emitting nested types before their parents so dependencies are satisfied.
func collectPythonTypedDicts(typeName string, schema *flagset.ObjectSchema) string {
	var b strings.Builder

	// First, recurse into nested objects to emit their TypedDicts
	for _, propName := range generators.SortedPropertyNames(schema.Properties) {
		propSchema := schema.Properties[propName]
		if propSchema != nil && propSchema.Type == "object" && propSchema.Properties != nil {
			nestedName := typeName + generators.ObjectTypeName(propName)
			b.WriteString(collectPythonTypedDicts(nestedName, propSchema))
			b.WriteString("\n")
		}
		// Handle arrays of objects
		if propSchema != nil && propSchema.Type == "array" && propSchema.Items != nil &&
			propSchema.Items.Type == "object" && propSchema.Items.Properties != nil {
			nestedName := typeName + generators.ObjectTypeName(propName) + "Item"
			b.WriteString(collectPythonTypedDicts(nestedName, propSchema.Items))
			b.WriteString("\n")
		}
	}

	// Now emit this TypedDict
	b.WriteString(fmt.Sprintf("class %s(TypedDict, total=False):\n", typeName))
	for _, propName := range generators.SortedPropertyNames(schema.Properties) {
		propSchema := schema.Properties[propName]
		isReq := generators.IsRequired(propName, schema.Required)

		nestedName := typeName + generators.ObjectTypeName(propName)
		pyType := pythonSchemaType(propSchema, nestedName)

		if isReq {
			b.WriteString(fmt.Sprintf("    %s: Required[%s]\n", propName, pyType))
		} else {
			b.WriteString(fmt.Sprintf("    %s: %s\n", propName, pyType))
		}
	}

	return b.String()
}

// generatePythonTypedDictDef generates Python TypedDict class definitions for a flag's object schema,
// including any nested object types.
func generatePythonTypedDictDef(flag flagset.Flag) string {
	if !generators.HasSchema(flag) {
		return ""
	}

	typeName := generators.ObjectTypeName(flag.Key)
	return collectPythonTypedDicts(typeName, flag.Schema)
}

// generatePythonHookDef generates a Python Hook class for runtime validation of a flag's object schema.
func generatePythonHookDef(flag flagset.Flag) string {
	if !generators.HasSchema(flag) {
		return ""
	}

	hookName := pythonHookName(flag)
	var b strings.Builder
	b.WriteString(fmt.Sprintf("class %s(Hook):\n", hookName))
	b.WriteString("    def after(\n")
	b.WriteString("        self, hook_context: HookContext, details: FlagEvaluationDetails, hints: dict\n")
	b.WriteString("    ):\n")
	b.WriteString("        value = details.value\n")
	b.WriteString(generatePythonValidation(flag.Schema, "value", flag.Key, "        "))
	b.WriteString("\n")

	return b.String()
}

// pySafeVarName converts a path to a safe Python variable name.
func pySafeVarName(path string) string {
	v := strings.ReplaceAll(path, ".", "_")
	v = strings.ReplaceAll(v, "[", "_")
	v = strings.ReplaceAll(v, "]", "")
	return v
}

// generatePythonValidation generates Python validation code for a schema.
func generatePythonValidation(schema *flagset.ObjectSchema, accessor string, path string, indent string) string {
	var b strings.Builder

	switch schema.Type {
	case "object":
		b.WriteString(fmt.Sprintf("%sif not isinstance(%s, dict):\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s    raise ValueError(\"%s: expected dict\")\n", indent, path))

		for _, req := range schema.Required {
			b.WriteString(fmt.Sprintf("%sif %q not in %s:\n", indent, req, accessor))
			b.WriteString(fmt.Sprintf("%s    raise ValueError(\"%s: missing required property '%s'\")\n", indent, path, req))
		}

		// Validate each property's type and recurse into nested schemas
		for _, propName := range generators.SortedPropertyNames(schema.Properties) {
			propSchema := schema.Properties[propName]
			propAccessor := fmt.Sprintf("%s[%q]", accessor, propName)
			propPath := fmt.Sprintf("%s.%s", path, propName)

			b.WriteString(fmt.Sprintf("%sif %q in %s:\n", indent, propName, accessor))
			b.WriteString(generatePythonValidation(propSchema, propAccessor, propPath, indent+"    "))
		}

		// Check additionalProperties: false
		if schema.AdditionalProperties != nil && !*schema.AdditionalProperties {
			allowedVar := pySafeVarName(path) + "_allowed"
			b.WriteString(fmt.Sprintf("%s%s = {", indent, allowedVar))
			for i, propName := range generators.SortedPropertyNames(schema.Properties) {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(fmt.Sprintf("%q", propName))
			}
			b.WriteString("}\n")
			b.WriteString(fmt.Sprintf("%sfor _key in %s:\n", indent, accessor))
			b.WriteString(fmt.Sprintf("%s    if _key not in %s:\n", indent, allowedVar))
			b.WriteString(fmt.Sprintf("%s        raise ValueError(f\"%s: unexpected property '{_key}'\")\n", indent, path))
		}

	case "array":
		b.WriteString(fmt.Sprintf("%sif not isinstance(%s, list):\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s    raise ValueError(\"%s: expected list\")\n", indent, path))

		if schema.Items != nil {
			idxVar := pySafeVarName(path) + "_idx"
			itemVar := pySafeVarName(path) + "_item"
			b.WriteString(fmt.Sprintf("%sfor %s, %s in enumerate(%s):\n", indent, idxVar, itemVar, accessor))
			b.WriteString(generatePythonArrayItemValidation(schema.Items, itemVar, path, idxVar, indent+"    "))
		}

	case "string":
		b.WriteString(fmt.Sprintf("%sif not isinstance(%s, str):\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s    raise ValueError(\"%s: expected str\")\n", indent, path))

	case "number":
		// In Python, bool is a subclass of int, so exclude it
		b.WriteString(fmt.Sprintf("%sif isinstance(%s, bool) or not isinstance(%s, (int, float)):\n", indent, accessor, accessor))
		b.WriteString(fmt.Sprintf("%s    raise ValueError(\"%s: expected number\")\n", indent, path))

	case "integer":
		b.WriteString(fmt.Sprintf("%sif not isinstance(%s, int) or isinstance(%s, bool):\n", indent, accessor, accessor))
		b.WriteString(fmt.Sprintf("%s    raise ValueError(\"%s: expected int\")\n", indent, path))

	case "boolean":
		b.WriteString(fmt.Sprintf("%sif not isinstance(%s, bool):\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s    raise ValueError(\"%s: expected bool\")\n", indent, path))
	}

	return b.String()
}

// generatePythonArrayItemValidation generates validation for array items with runtime index in error paths.
func generatePythonArrayItemValidation(schema *flagset.ObjectSchema, itemVar string, arrayPath string, idxVar string, indent string) string {
	var b strings.Builder

	switch schema.Type {
	case "string":
		b.WriteString(fmt.Sprintf("%sif not isinstance(%s, str):\n", indent, itemVar))
		b.WriteString(fmt.Sprintf("%s    raise ValueError(f\"%s[{%s}]: expected str\")\n", indent, arrayPath, idxVar))
	case "number":
		b.WriteString(fmt.Sprintf("%sif isinstance(%s, bool) or not isinstance(%s, (int, float)):\n", indent, itemVar, itemVar))
		b.WriteString(fmt.Sprintf("%s    raise ValueError(f\"%s[{%s}]: expected number\")\n", indent, arrayPath, idxVar))
	case "integer":
		b.WriteString(fmt.Sprintf("%sif not isinstance(%s, int) or isinstance(%s, bool):\n", indent, itemVar, itemVar))
		b.WriteString(fmt.Sprintf("%s    raise ValueError(f\"%s[{%s}]: expected int\")\n", indent, arrayPath, idxVar))
	case "boolean":
		b.WriteString(fmt.Sprintf("%sif not isinstance(%s, bool):\n", indent, itemVar))
		b.WriteString(fmt.Sprintf("%s    raise ValueError(f\"%s[{%s}]: expected bool\")\n", indent, arrayPath, idxVar))
	case "object", "array":
		b.WriteString(generatePythonValidation(schema, itemVar, fmt.Sprintf("%s[item]", arrayPath), indent))
	}

	return b.String()
}

// pythonFlagReturnType returns the Python return type for a flag.
// For schema-typed object flags, it returns the TypedDict class name.
// For non-schema object flags, it returns "object".
func pythonFlagReturnType(flag flagset.Flag) string {
	if generators.HasSchema(flag) {
		return generators.ObjectTypeName(flag.Key)
	}
	if flag.Type == flagset.ObjectType {
		return "object"
	}
	return openFeatureType(flag.Type)
}

// pythonHookName returns the hook class name for a typed object flag.
func pythonHookName(flag flagset.Flag) string {
	return generators.ObjectTypeName(flag.Key) + "Hook"
}

// pythonHookInjection generates the Python code block that prepends the validation hook
// to flag_evaluation_options. This is used in the template to avoid repeating the same
// 5-line block across all 4 method variants (sync/async × value/details).
func pythonHookInjection(flag flagset.Flag) string {
	if !generators.HasSchema(flag) {
		return ""
	}
	hookName := pythonHookName(flag)
	return fmt.Sprintf(`        if flag_evaluation_options is None:
            flag_evaluation_options = FlagEvaluationOptions(hooks=[%s()])
        else:
            flag_evaluation_options = FlagEvaluationOptions(
                hooks=[*(flag_evaluation_options.hooks or []), %s()],
            )`, hookName, hookName)
}

func (g *PythonGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"OpenFeatureType":          openFeatureType,
		"TypedGetMethodSync":       typedGetMethodSync,
		"TypedGetMethodAsync":      typedGetMethodAsync,
		"TypedDetailsMethodSync":   typedDetailsMethodSync,
		"TypedDetailsMethodAsync":  typedDetailsMethodAsync,
		"PythonBoolLiteral":        pythonBoolLiteral,
		"ToPythonDict":             toPythonDict,
		"HasSchema":                generators.HasSchema,
		"PythonTypedDictDef":       generatePythonTypedDictDef,
		"PythonHookDef":            generatePythonHookDef,
		"PythonFlagReturnType":     pythonFlagReturnType,
		"PythonHookName":           pythonHookName,
		"PythonHookInjection":      pythonHookInjection,
		"HasObjectFlagsWithSchema": generators.HasObjectFlagsWithSchema,
	}

	newParams := &generators.Params[any]{
		OutputPath:        params.OutputPath,
		TemplatePath:      params.TemplatePath,
		RuntimeValidation: params.RuntimeValidation,
		Custom:            Params{},
	}

	return g.GenerateFile(funcs, pythonTmpl, newParams, "openfeature.py")
}

// NewGenerator creates a generator for Python.
func NewGenerator(fs *flagset.Flagset) *PythonGenerator {
	return &PythonGenerator{
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{}),
	}
}
