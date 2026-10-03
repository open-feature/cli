package generators

import (
	"testing"

	"github.com/open-feature/cli/internal/flagset"
)

func TestHasSchema(t *testing.T) {
	tests := []struct {
		name string
		flag flagset.Flag
		want bool
	}{
		{
			name: "object flag with schema",
			flag: flagset.Flag{
				Type:   flagset.ObjectType,
				Schema: &flagset.ObjectSchema{Type: "object"},
			},
			want: true,
		},
		{
			name: "object flag without schema",
			flag: flagset.Flag{
				Type:   flagset.ObjectType,
				Schema: nil,
			},
			want: false,
		},
		{
			name: "string flag with schema is false",
			flag: flagset.Flag{
				Type:   flagset.StringType,
				Schema: &flagset.ObjectSchema{Type: "string"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasSchema(tt.flag); got != tt.want {
				t.Errorf("HasSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObjectTypeName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"themeCustomization", "ThemeCustomization"},
		{"primary_color", "PrimaryColor"},
		{"simple", "Simple"},
		{"my-flag-key", "MyFlagKey"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ObjectTypeName(tt.input); got != tt.want {
				t.Errorf("ObjectTypeName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestHasObjectFlagsWithSchema(t *testing.T) {
	tests := []struct {
		name  string
		flags []flagset.Flag
		want  bool
	}{
		{
			name:  "empty",
			flags: nil,
			want:  false,
		},
		{
			name: "no schemas",
			flags: []flagset.Flag{
				{Type: flagset.ObjectType, Schema: nil},
				{Type: flagset.StringType},
			},
			want: false,
		},
		{
			name: "has schema",
			flags: []flagset.Flag{
				{Type: flagset.StringType},
				{Type: flagset.ObjectType, Schema: &flagset.ObjectSchema{Type: "object"}},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasObjectFlagsWithSchema(tt.flags); got != tt.want {
				t.Errorf("HasObjectFlagsWithSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}
