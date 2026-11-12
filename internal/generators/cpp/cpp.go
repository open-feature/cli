package cpp

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

type CppGenerator struct {
	generators.CommonGenerator
}

type Params struct {
	Namespace    string
	CppNamespace string
}

//go:embed cpp.tmpl
var cppTmpl string

// OpenFeatureType maps FlagType to the corresponding C++ type name.
func OpenFeatureType(t flagset.FlagType) string {
	switch t {
	case flagset.BoolType:
		return "bool"
	case flagset.StringType:
		return "std::string"
	case flagset.IntType:
		return "int64_t"
	case flagset.FloatType:
		return "double"
	case flagset.ObjectType:
		return "openfeature::Value"
	default:
		return ""
	}
}

// OpenFeatureMethodType maps FlagType to the method name type segment.
func OpenFeatureMethodType(t flagset.FlagType) string {
	switch t {
	case flagset.BoolType:
		return "Boolean"
	case flagset.StringType:
		return "String"
	case flagset.IntType:
		return "Integer"
	case flagset.FloatType:
		return "Double"
	case flagset.ObjectType:
		return "Object"
	default:
		return ""
	}
}

// QuoteCppString formats s as a C++ string literal using fixed-width octal
// escapes for non-printable bytes so following characters cannot be consumed by
// greedy hexadecimal escapes.
func QuoteCppString(s string) string {
	var builder strings.Builder
	builder.Grow(len(s) + 2)
	builder.WriteByte('"')
	for i := range len(s) {
		b := s[i]
		switch b {
		case '"':
			builder.WriteString(`\"`)
		case '\\':
			builder.WriteString(`\\`)
		case '\n':
			builder.WriteString(`\n`)
		case '\r':
			builder.WriteString(`\r`)
		case '\t':
			builder.WriteString(`\t`)
		default:
			if b >= 0x20 && b <= 0x7e {
				builder.WriteByte(b)
			} else {
				fmt.Fprintf(&builder, "\\%03o", b)
			}
		}
	}
	builder.WriteByte('"')
	return builder.String()
}

// FormatComment prefixes every line of text with prefix so multiline strings
// remain safely inside C++ line comments.
func FormatComment(prefix, text string) string {
	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")
	lines := strings.Split(normalized, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

// FormatDefaultValue formats flag default values into valid C++ literals.
func FormatDefaultValue(flag flagset.Flag) string {
	switch flag.Type {
	case flagset.StringType:
		return QuoteCppString(fmt.Sprintf("%v", flag.DefaultValue))
	case flagset.BoolType:
		if flag.DefaultValue == true {
			return "true"
		}
		return "false"
	case flagset.IntType:
		switch v := flag.DefaultValue.(type) {
		case float64:
			return fmt.Sprintf("%d", int64(v))
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return fmt.Sprintf("%d", v)
		default:
			return fmt.Sprintf("%v", flag.DefaultValue)
		}
	case flagset.FloatType:
		valStr := fmt.Sprintf("%v", flag.DefaultValue)
		if !strings.Contains(valStr, ".") && !strings.ContainsAny(valStr, "eE") {
			valStr += ".0"
		}
		return valStr
	case flagset.ObjectType:
		return ToValueLiteral(flag.DefaultValue)
	default:
		return fmt.Sprintf("%v", flag.DefaultValue)
	}
}

// ToValueLiteral recursively formats values into C++ openfeature::Value initializers.
func ToValueLiteral(value any) string {
	switch val := value.(type) {
	case nil:
		return "openfeature::Value()"
	case bool:
		return fmt.Sprintf("openfeature::Value(%t)", val)
	case string:
		return fmt.Sprintf("openfeature::Value(std::string(%s))", QuoteCppString(val))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("openfeature::Value(static_cast<int64_t>(%v))", val)
	case float32, float64:
		if fmt.Sprintf("%.17g", val) == fmt.Sprintf("%.0f", val) {
			return fmt.Sprintf("openfeature::Value(static_cast<int64_t>(%.0f))", val)
		}
		return fmt.Sprintf("openfeature::Value(static_cast<double>(%v))", val)
	case map[string]any:
		return mapToValueLiteral(val)
	case []any:
		return sliceToValueLiteral(val)
	default:
		return "openfeature::Value()"
	}
}

func mapToValueLiteral(m map[string]any) string {
	keys := slices.Sorted(maps.Keys(m))
	var builder strings.Builder
	builder.WriteString("openfeature::Value(std::map<std::string, openfeature::Value>{")
	for i, k := range keys {
		if i > 0 {
			builder.WriteString(", ")
		}
		fmt.Fprintf(&builder, "{%s, %s}", QuoteCppString(k), ToValueLiteral(m[k]))
	}
	builder.WriteString("})")
	return builder.String()
}

func sliceToValueLiteral(s []any) string {
	var builder strings.Builder
	builder.WriteString("openfeature::Value(std::vector<openfeature::Value>{")
	for i, elem := range s {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(ToValueLiteral(elem))
	}
	builder.WriteString("})")
	return builder.String()
}

func (g *CppGenerator) Generate(params *generators.Params[Params]) error {
	funcs := template.FuncMap{
		"CppQuote":              QuoteCppString,
		"FormatComment":         FormatComment,
		"OpenFeatureType":       OpenFeatureType,
		"OpenFeatureMethodType": OpenFeatureMethodType,
		"FormatDefaultValue":    FormatDefaultValue,
		"ToValueLiteral":        ToValueLiteral,
	}

	custom := params.Custom
	if custom.CppNamespace == "" && custom.Namespace != "" {
		custom.CppNamespace = custom.Namespace
	} else if custom.Namespace == "" && custom.CppNamespace != "" {
		custom.Namespace = custom.CppNamespace
	}

	newParams := &generators.Params[any]{
		OutputPath:   params.OutputPath,
		TemplatePath: params.TemplatePath,
		Custom:       custom,
	}

	return g.GenerateFile(funcs, cppTmpl, newParams, "openfeature.h")
}

// NewGenerator creates a generator for C++.
func NewGenerator(fs *flagset.Flagset) *CppGenerator {
	return &CppGenerator{
		CommonGenerator: *generators.NewGenerator(fs, map[flagset.FlagType]bool{}),
	}
}
