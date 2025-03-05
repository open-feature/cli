package react

import (
	_ "embed"
	"text/template"

	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators"
)

type ReactGenerator struct {
	generators.CommonGenerator
}

//go:embed react.tmpl
var reactTmpl string

func (g *ReactGenerator) Generate() error {
	funcs := template.FuncMap{}

	return g.GenerateFile(funcs, reactTmpl, &generators.Params{
		OutputPath: "test.ts",
	})
}

// NewGenerator creates a generator for React.
func NewGenerator(fs *flagset.Flagset) *ReactGenerator {
	return &ReactGenerator{
		CommonGenerator: *generators.NewCommonGenerator(
			fs,
			generators.WithStability(generators.Alpha),
			generators.WithUnsupportedFlagType(flagset.ObjectType),
		),
	}
}
