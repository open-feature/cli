package generators

import (
	"slices"
	"sort"

	"github.com/iancoleman/strcase"
	"github.com/open-feature/cli/internal/flagset"
)

// HasSchema returns true if the flag has an object schema defined.
func HasSchema(flag flagset.Flag) bool {
	return flag.Schema != nil && flag.Type == flagset.ObjectType
}

// ObjectTypeName returns the PascalCase type name for a flag's object schema.
func ObjectTypeName(flagKey string) string {
	return strcase.ToCamel(flagKey)
}

// HasObjectFlagsWithSchema returns true if the flagset contains any object flags with schemas.
func HasObjectFlagsWithSchema(flags []flagset.Flag) bool {
	for _, f := range flags {
		if HasSchema(f) {
			return true
		}
	}
	return false
}

// SortedPropertyNames returns property names sorted alphabetically for deterministic output.
func SortedPropertyNames(props map[string]*flagset.ObjectSchema) []string {
	names := make([]string, 0, len(props))
	for name := range props {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// IsRequired checks if a property name is in the required list.
func IsRequired(name string, required []string) bool {
	return slices.Contains(required, name)
}
