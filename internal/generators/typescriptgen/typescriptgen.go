// Package typescriptgen provides shared TypeScript type generation utilities
// for all TypeScript-based generators (Node.js, React, Angular, NestJS).
package typescriptgen

import (
	"fmt"
	"strings"

	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators"
)

// tsSchemaType converts an ObjectSchema type to a TypeScript type string.
func tsSchemaType(schema *flagset.ObjectSchema, required bool) string {
	if schema == nil {
		return "JsonValue"
	}

	switch schema.Type {
	case "string":
		return "string"
	case "number", "integer":
		return "number"
	case "boolean":
		return "boolean"
	case "array":
		if schema.Items != nil {
			return tsSchemaType(schema.Items, true) + "[]"
		}
		return "unknown[]"
	case "object":
		if schema.Properties != nil {
			return generateInlineInterface(schema)
		}
		return "Record<string, unknown>"
	default:
		return "unknown"
	}
}

// generateInlineInterface generates an inline TypeScript interface body.
func generateInlineInterface(schema *flagset.ObjectSchema) string {
	var b strings.Builder
	b.WriteString("{\n")

	for _, propName := range generators.SortedPropertyNames(schema.Properties) {
		propSchema := schema.Properties[propName]
		isReq := generators.IsRequired(propName, schema.Required)
		tsType := tsSchemaType(propSchema, isReq)

		if isReq {
			b.WriteString(fmt.Sprintf("  %s: %s;\n", propName, tsType))
		} else {
			b.WriteString(fmt.Sprintf("  %s?: %s;\n", propName, tsType))
		}
	}

	b.WriteString("}")
	return b.String()
}

// GenerateInterfaceDef generates a top-level TypeScript interface definition for a flag's object schema.
func GenerateInterfaceDef(flag flagset.Flag) string {
	if !generators.HasSchema(flag) {
		return ""
	}

	typeName := generators.ObjectTypeName(flag.Key)
	var b strings.Builder
	b.WriteString(fmt.Sprintf("export interface %s ", typeName))
	b.WriteString(generateInlineInterface(flag.Schema))
	b.WriteString("\n")

	return b.String()
}

// GenerateValidationHookDef generates a TypeScript validation hook function for a flag's object schema.
func GenerateValidationHookDef(flag flagset.Flag) string {
	if !generators.HasSchema(flag) {
		return ""
	}

	hookName := ValidationHookName(flag)
	var b strings.Builder
	b.WriteString(fmt.Sprintf("function %s(): Hook {\n", hookName))
	b.WriteString("  return {\n")
	b.WriteString("    after: (_hookContext: Readonly<HookContext>, evaluationDetails: EvaluationDetails<JsonValue>) => {\n")
	b.WriteString("      const value = evaluationDetails.value;\n")
	b.WriteString(generateTSValidation(flag.Schema, "value", flag.Key, "      "))
	b.WriteString("    },\n")
	b.WriteString("  };\n")
	b.WriteString("}\n")

	return b.String()
}

// tsSafeVarName converts a path to a safe TypeScript variable name.
func tsSafeVarName(path string) string {
	v := strings.ReplaceAll(path, ".", "_")
	v = strings.ReplaceAll(v, "[", "_")
	v = strings.ReplaceAll(v, "]", "")
	return v
}

// generateTSValidation generates TypeScript validation code for a schema.
func generateTSValidation(schema *flagset.ObjectSchema, accessor string, path string, indent string) string {
	var b strings.Builder

	switch schema.Type {
	case "object":
		b.WriteString(fmt.Sprintf("%sif (typeof %s !== 'object' || %s === null || Array.isArray(%s)) {\n", indent, accessor, accessor, accessor))
		b.WriteString(fmt.Sprintf("%s  throw new Error('%s: expected object');\n", indent, path))
		b.WriteString(fmt.Sprintf("%s}\n", indent))

		castVar := tsSafeVarName(path)
		b.WriteString(fmt.Sprintf("%sconst %sObj = %s as Record<string, unknown>;\n", indent, castVar, accessor))

		for _, req := range schema.Required {
			b.WriteString(fmt.Sprintf("%sif (%sObj[%q] === undefined) {\n", indent, castVar, req))
			b.WriteString(fmt.Sprintf("%s  throw new Error('%s: missing required property %q');\n", indent, path, req))
			b.WriteString(fmt.Sprintf("%s}\n", indent))
		}

		// Validate each property's type and recurse into nested schemas
		for _, propName := range generators.SortedPropertyNames(schema.Properties) {
			propSchema := schema.Properties[propName]
			propAccessor := fmt.Sprintf("%sObj[%q]", castVar, propName)
			propPath := fmt.Sprintf("%s.%s", path, propName)

			b.WriteString(fmt.Sprintf("%sif (%s !== undefined) {\n", indent, propAccessor))
			b.WriteString(generateTSValidation(propSchema, propAccessor, propPath, indent+"  "))
			b.WriteString(fmt.Sprintf("%s}\n", indent))
		}

		// Check additionalProperties: false
		if schema.AdditionalProperties != nil && !*schema.AdditionalProperties {
			allowedVar := tsSafeVarName(path) + "Allowed"
			b.WriteString(fmt.Sprintf("%sconst %s = new Set([", indent, allowedVar))
			for i, propName := range generators.SortedPropertyNames(schema.Properties) {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(fmt.Sprintf("%q", propName))
			}
			b.WriteString("]);\n")
			b.WriteString(fmt.Sprintf("%sfor (const key of Object.keys(%sObj)) {\n", indent, castVar))
			b.WriteString(fmt.Sprintf("%s  if (!%s.has(key)) {\n", indent, allowedVar))
			b.WriteString(fmt.Sprintf("%s    throw new Error(`%s: unexpected property \"${key}\"`);\n", indent, path))
			b.WriteString(fmt.Sprintf("%s  }\n", indent))
			b.WriteString(fmt.Sprintf("%s}\n", indent))
		}

	case "array":
		b.WriteString(fmt.Sprintf("%sif (!Array.isArray(%s)) {\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s  throw new Error('%s: expected array');\n", indent, path))
		b.WriteString(fmt.Sprintf("%s}\n", indent))

		if schema.Items != nil {
			idxVar := tsSafeVarName(path) + "Idx"
			b.WriteString(fmt.Sprintf("%sfor (let %s = 0; %s < %s.length; %s++) {\n", indent, idxVar, idxVar, accessor, idxVar))
			itemAccessor := fmt.Sprintf("%s[%s]", accessor, idxVar)
			b.WriteString(generateTSArrayItemValidation(schema.Items, itemAccessor, path, idxVar, indent+"  "))
			b.WriteString(fmt.Sprintf("%s}\n", indent))
		}

	case "string":
		b.WriteString(fmt.Sprintf("%sif (typeof %s !== 'string') {\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s  throw new Error('%s: expected string');\n", indent, path))
		b.WriteString(fmt.Sprintf("%s}\n", indent))

	case "number":
		b.WriteString(fmt.Sprintf("%sif (typeof %s !== 'number') {\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s  throw new Error('%s: expected number');\n", indent, path))
		b.WriteString(fmt.Sprintf("%s}\n", indent))

	case "integer":
		b.WriteString(fmt.Sprintf("%sif (typeof %s !== 'number' || !Number.isInteger(%s)) {\n", indent, accessor, accessor))
		b.WriteString(fmt.Sprintf("%s  throw new Error('%s: expected integer');\n", indent, path))
		b.WriteString(fmt.Sprintf("%s}\n", indent))

	case "boolean":
		b.WriteString(fmt.Sprintf("%sif (typeof %s !== 'boolean') {\n", indent, accessor))
		b.WriteString(fmt.Sprintf("%s  throw new Error('%s: expected boolean');\n", indent, path))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
	}

	return b.String()
}

// generateTSArrayItemValidation generates validation for array items with runtime index in error paths.
func generateTSArrayItemValidation(schema *flagset.ObjectSchema, itemAccessor string, arrayPath string, idxVar string, indent string) string {
	var b strings.Builder

	switch schema.Type {
	case "string":
		b.WriteString(fmt.Sprintf("%sif (typeof %s !== 'string') {\n", indent, itemAccessor))
		b.WriteString(fmt.Sprintf("%s  throw new Error(`%s[${%s}]: expected string`);\n", indent, arrayPath, idxVar))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
	case "number":
		b.WriteString(fmt.Sprintf("%sif (typeof %s !== 'number') {\n", indent, itemAccessor))
		b.WriteString(fmt.Sprintf("%s  throw new Error(`%s[${%s}]: expected number`);\n", indent, arrayPath, idxVar))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
	case "integer":
		b.WriteString(fmt.Sprintf("%sif (typeof %s !== 'number' || !Number.isInteger(%s)) {\n", indent, itemAccessor, itemAccessor))
		b.WriteString(fmt.Sprintf("%s  throw new Error(`%s[${%s}]: expected integer`);\n", indent, arrayPath, idxVar))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
	case "boolean":
		b.WriteString(fmt.Sprintf("%sif (typeof %s !== 'boolean') {\n", indent, itemAccessor))
		b.WriteString(fmt.Sprintf("%s  throw new Error(`%s[${%s}]: expected boolean`);\n", indent, arrayPath, idxVar))
		b.WriteString(fmt.Sprintf("%s}\n", indent))
	case "object", "array":
		// For nested objects/arrays in arrays, use static path since nesting gets complex
		b.WriteString(generateTSValidation(schema, itemAccessor, fmt.Sprintf("%s[item]", arrayPath), indent))
	}

	return b.String()
}

// FlagReturnType returns the TypeScript return type for a flag.
func FlagReturnType(flag flagset.Flag) string {
	if generators.HasSchema(flag) {
		return generators.ObjectTypeName(flag.Key)
	}
	if flag.Type == flagset.ObjectType {
		return "JsonValue"
	}
	return "" // handled by existing template logic for non-object types
}

// ValidationHookName returns the hook function name for a typed object flag.
func ValidationHookName(flag flagset.Flag) string {
	return "create" + generators.ObjectTypeName(flag.Key) + "Hook"
}
