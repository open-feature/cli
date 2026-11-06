package java

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"
	"text/template"
	"unicode/utf16"

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

// javaReturnType is the Java type exposed by the generated accessors. Object
// flags are evaluated through the SDK's Value type.
func javaReturnType(t flagset.FlagType) string {
	if t == flagset.ObjectType {
		return "Value"
	}
	return openFeatureType(t)
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

		fmt.Fprintf(&builder, "%q, %s", key, formatNestedValue(val))
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

// javaString renders s as a Java string literal.
func javaString(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch {
		case r == '"':
			b.WriteString(`\"`)
		case r == '\\':
			b.WriteString(`\\`)
		case r == '\n':
			b.WriteString(`\n`)
		case r == '\r':
			b.WriteString(`\r`)
		case r == '\t':
			b.WriteString(`\t`)
		case r < 0x20 || r == 0x7f:
			fmt.Fprintf(&b, `\u%04x`, r)
		case r > 0xffff:
			r1, r2 := utf16.EncodeRune(r)
			fmt.Fprintf(&b, `\u%04x\u%04x`, r1, r2)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}

// javaNumber renders a JSON number as a Java int or double literal, as those
// are the numeric types supported by dev.openfeature.sdk.Value.
func javaNumber(f float64) string {
	if f == math.Trunc(f) && f >= math.MinInt32 && f <= math.MaxInt32 {
		return strconv.FormatInt(int64(f), 10)
	}
	s := strconv.FormatFloat(f, 'f', -1, 64)
	if !strings.Contains(s, ".") {
		s += ".0"
	}
	return s
}

// toValueLiteral renders value as a Java expression of type
// dev.openfeature.sdk.Value. Structures use MutableStructure rather than
// Map.of, which is limited to 10 entries and rejects nulls.
func toValueLiteral(value any) string {
	switch val := value.(type) {
	case nil:
		return "new Value()"
	case string:
		return "new Value(" + javaString(val) + ")"
	case bool:
		return fmt.Sprintf("new Value(%t)", val)
	case int:
		return "new Value(" + javaNumber(float64(val)) + ")"
	case int64:
		return "new Value(" + javaNumber(float64(val)) + ")"
	case float64:
		return "new Value(" + javaNumber(val) + ")"
	case map[string]any:
		var b strings.Builder
		b.WriteString("new Value(new MutableStructure()")
		for _, key := range slices.Sorted(maps.Keys(val)) {
			fmt.Fprintf(&b, ".add(%s, %s)", javaString(key), toValueLiteral(val[key]))
		}
		b.WriteString(")")
		return b.String()
	case []any:
		elems := make([]string, len(val))
		for i, elem := range val {
			elems[i] = toValueLiteral(elem)
		}
		return "new Value(List.<Value>of(" + strings.Join(elems, ", ") + "))"
	default:
		return "new Value()"
	}
}

// toJavadocJSON renders value as JSON that is safe inside a Java block comment.
func toJavadocJSON(value any) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(value); err != nil {
		return "null"
	}
	s := strings.TrimSpace(buf.String())
	s = strings.ReplaceAll(s, "*/", "*\\/")
	// javac processes unicode escapes even inside comments.
	return strings.ReplaceAll(s, `\u`, `\\u`)
}

func (g *JavaGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"OpenFeatureType":    openFeatureType,
		"FormatDefaultValue": formatDefaultValueForJava,
		"ToMapLiteral":       toMapLiteral,
		"ToValueLiteral":     toValueLiteral,
		"ToJavadocJSON":      toJavadocJSON,
		"JavaReturnType":     javaReturnType,
	}

	newParams := &generators.Params[any]{
		OutputPath:   params.OutputPath,
		TemplatePath: params.TemplatePath,
		Custom:       params.Custom,
	}

	return g.GenerateFile(funcs, javaTmpl, newParams, "OpenFeature.java")
}

// NewGenerator creates a generator for Java.
func NewGenerator(fs *flagset.Flagset) *JavaGenerator {
	return &JavaGenerator{
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{}),
	}
}
