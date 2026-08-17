package generators

import (
	"testing"

	"github.com/iancoleman/strcase"
	"github.com/open-feature/cli/internal/flagset"
)

func keys(fs *flagset.Flagset) []string {
	out := make([]string, 0, len(fs.Flags))
	for _, f := range fs.Flags {
		out = append(out, f.Key)
	}
	return out
}

func TestFilterReservedFlags(t *testing.T) {
	tests := []struct {
		name      string
		generator string
		reserved  map[string]bool
		transform func(string) string
		input     []string
		want      []string
	}{
		{
			name:      "drops colliding flag, keeps the rest",
			generator: "Go",
			reserved:  map[string]bool{"Client": true},
			transform: strcase.ToCamel,
			input:     []string{"client", "discountPercentage"},
			want:      []string{"discountPercentage"},
		},
		{
			name:      "no collision keeps every flag",
			generator: "Go",
			reserved:  map[string]bool{"Client": true},
			transform: strcase.ToCamel,
			input:     []string{"discountPercentage", "enableFeatureA"},
			want:      []string{"discountPercentage", "enableFeatureA"},
		},
		{
			name:      "react hook name collision",
			generator: "React",
			reserved:  map[string]bool{"useOpenFeatureClient": true},
			transform: func(key string) string { return "use" + strcase.ToCamel(key) },
			input:     []string{"openFeatureClient", "greetingMessage"},
			want:      []string{"greetingMessage"},
		},
		{
			name:      "node lower-camel collision",
			generator: "Node.js",
			reserved:  map[string]bool{"client": true},
			transform: strcase.ToLowerCamel,
			input:     []string{"client", "greetingMessage"},
			want:      []string{"greetingMessage"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &flagset.Flagset{}
			for _, k := range tt.input {
				fs.Flags = append(fs.Flags, flagset.Flag{Key: k, Type: flagset.StringType})
			}

			got := keys(FilterReservedFlags(fs, tt.generator, tt.reserved, tt.transform))

			if len(got) != len(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("got %v, want %v", got, tt.want)
				}
			}
		})
	}
}
