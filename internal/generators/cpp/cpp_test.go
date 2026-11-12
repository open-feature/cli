package cpp

import (
	"encoding/json"
	"testing"

	"github.com/open-feature/cli/internal/filesystem"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators"
	"github.com/spf13/afero"
)

func TestOpenFeatureType(t *testing.T) {
	testCases := []struct {
		name     string
		flagType flagset.FlagType
		expected string
	}{
		{name: "bool", flagType: flagset.BoolType, expected: "bool"},
		{name: "string", flagType: flagset.StringType, expected: "std::string"},
		{name: "integer", flagType: flagset.IntType, expected: "int64_t"},
		{name: "float", flagType: flagset.FloatType, expected: "double"},
		{name: "object", flagType: flagset.ObjectType, expected: "openfeature::Value"},
		{name: "unknown", flagType: flagset.UnknownFlagType, expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := OpenFeatureType(tc.flagType)
			if got != tc.expected {
				t.Fatalf("OpenFeatureType(%v) = %q; want %q", tc.flagType, got, tc.expected)
			}
		})
	}
}

func TestOpenFeatureMethodType(t *testing.T) {
	testCases := []struct {
		name     string
		flagType flagset.FlagType
		expected string
	}{
		{name: "bool", flagType: flagset.BoolType, expected: "Boolean"},
		{name: "string", flagType: flagset.StringType, expected: "String"},
		{name: "integer", flagType: flagset.IntType, expected: "Integer"},
		{name: "float", flagType: flagset.FloatType, expected: "Double"},
		{name: "object", flagType: flagset.ObjectType, expected: "Object"},
		{name: "unknown", flagType: flagset.UnknownFlagType, expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := OpenFeatureMethodType(tc.flagType)
			if got != tc.expected {
				t.Fatalf("OpenFeatureMethodType(%v) = %q; want %q", tc.flagType, got, tc.expected)
			}
		})
	}
}

func TestFormatDefaultValue_Boolean(t *testing.T) {
	trueFlag := flagset.Flag{Type: flagset.BoolType, DefaultValue: true}
	if got := FormatDefaultValue(trueFlag); got != "true" {
		t.Fatalf("FormatDefaultValue(true) = %q; want %q", got, "true")
	}

	falseFlag := flagset.Flag{Type: flagset.BoolType, DefaultValue: false}
	if got := FormatDefaultValue(falseFlag); got != "false" {
		t.Fatalf("FormatDefaultValue(false) = %q; want %q", got, "false")
	}
}

func TestFormatDefaultValue_String(t *testing.T) {
	flag := flagset.Flag{Type: flagset.StringType, DefaultValue: "Hello world"}
	expected := `"Hello world"`
	if got := FormatDefaultValue(flag); got != expected {
		t.Fatalf("FormatDefaultValue(string) = %q; want %q", got, expected)
	}
}

func TestFormatDefaultValue_Integer(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected string
	}{
		{name: "int", input: int(42), expected: "42"},
		{name: "int64", input: int64(100), expected: "100"},
		{name: "float64 from json", input: float64(50), expected: "50"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flag := flagset.Flag{Type: flagset.IntType, DefaultValue: tc.input}
			got := FormatDefaultValue(flag)
			if got != tc.expected {
				t.Fatalf("FormatDefaultValue(int: %v) = %q; want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestFormatDefaultValue_Float(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected string
	}{
		{name: "fractional float", input: 0.15, expected: "0.15"},
		{name: "whole number float", input: float64(1), expected: "1.0"},
		{name: "zero float", input: float64(0), expected: "0.0"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flag := flagset.Flag{Type: flagset.FloatType, DefaultValue: tc.input}
			got := FormatDefaultValue(flag)
			if got != tc.expected {
				t.Fatalf("FormatDefaultValue(float: %v) = %q; want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestFormatDefaultValue_Object(t *testing.T) {
	flag := flagset.Flag{
		Type: flagset.ObjectType,
		DefaultValue: map[string]any{
			"color": "blue",
		},
	}
	expected := `openfeature::Value(std::map<std::string, openfeature::Value>{{"color", openfeature::Value(std::string("blue"))}})`
	if got := FormatDefaultValue(flag); got != expected {
		t.Fatalf("FormatDefaultValue(object) = %q; want %q", got, expected)
	}
}

func TestToValueLiteral_Primitives(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected string
	}{
		{name: "nil", input: nil, expected: "openfeature::Value()"},
		{name: "bool true", input: true, expected: "openfeature::Value(true)"},
		{name: "bool false", input: false, expected: "openfeature::Value(false)"},
		{name: "string", input: "test", expected: `openfeature::Value(std::string("test"))`},
		{name: "string with quotes", input: `he said "hi"`, expected: `openfeature::Value(std::string("he said \"hi\""))`},
		{name: "int", input: int(42), expected: "openfeature::Value(static_cast<int64_t>(42))"},
		{name: "int64", input: int64(100), expected: "openfeature::Value(static_cast<int64_t>(100))"},
		{name: "float64", input: float64(3.14), expected: "openfeature::Value(static_cast<double>(3.14))"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := ToValueLiteral(tc.input)
			if got != tc.expected {
				t.Fatalf("ToValueLiteral(%v) = %q; want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestToValueLiteral_Vector(t *testing.T) {
	input := []any{
		"first",
		int64(2),
		true,
		nil,
	}
	expected := `openfeature::Value(std::vector<openfeature::Value>{openfeature::Value(std::string("first")), openfeature::Value(static_cast<int64_t>(2)), openfeature::Value(true), openfeature::Value()})`
	got := ToValueLiteral(input)
	if got != expected {
		t.Fatalf("ToValueLiteral(slice) = %q; want %q", got, expected)
	}
}

func TestToValueLiteral_MapSortedKeys(t *testing.T) {
	input := map[string]any{
		"zebra": "last",
		"apple": "first",
		"mango": int(5),
	}
	expected := `openfeature::Value(std::map<std::string, openfeature::Value>{{"apple", openfeature::Value(std::string("first"))}, {"mango", openfeature::Value(static_cast<int64_t>(5))}, {"zebra", openfeature::Value(std::string("last"))}})`
	got := ToValueLiteral(input)
	if got != expected {
		t.Fatalf("ToValueLiteral(map) = %q; want %q", got, expected)
	}
}

func TestToValueLiteral_NestedObject(t *testing.T) {
	input := map[string]any{
		"config": map[string]any{
			"enabled": true,
		},
		"tags": []any{
			"prod",
			nil,
		},
	}
	expected := `openfeature::Value(std::map<std::string, openfeature::Value>{{"config", openfeature::Value(std::map<std::string, openfeature::Value>{{"enabled", openfeature::Value(true)}})}, {"tags", openfeature::Value(std::vector<openfeature::Value>{openfeature::Value(std::string("prod")), openfeature::Value()})}})`
	got := ToValueLiteral(input)
	if got != expected {
		t.Fatalf("ToValueLiteral(nested) = %q; want %q", got, expected)
	}
}

func TestToValueLiteral_JsonUnmarshaledNumbers(t *testing.T) {
	rawJSON := `{"ratio": 0.75, "retries": 3}`
	var unmarshaled map[string]any
	err := json.Unmarshal([]byte(rawJSON), &unmarshaled)
	if err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	expected := `openfeature::Value(std::map<std::string, openfeature::Value>{{"ratio", openfeature::Value(static_cast<double>(0.75))}, {"retries", openfeature::Value(static_cast<int64_t>(3))}})`
	got := ToValueLiteral(unmarshaled)
	if got != expected {
		t.Fatalf("ToValueLiteral(unmarshaled) = %q; want %q", got, expected)
	}
}

func TestGenerate_Execution(t *testing.T) {
	fs := afero.NewMemMapFs()
	filesystem.SetFileSystem(fs)

	flagSet := &flagset.Flagset{
		Flags: []flagset.Flag{
			{
				Key:          "testFlag",
				Type:         flagset.BoolType,
				DefaultValue: true,
				Description:  "A test flag",
			},
		},
	}

	generator := NewGenerator(flagSet)
	params := &generators.Params[Params]{
		OutputPath: "output",
		Custom: Params{
			Namespace:    "my_namespace",
			CppNamespace: "my_namespace",
		},
	}

	err := generator.Generate(params)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	data, err := afero.ReadFile(fs, "output/openfeature.h")
	if err != nil {
		t.Fatalf("Failed to read generated output: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Generated file is empty")
	}
}
