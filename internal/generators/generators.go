package generators

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"maps"

	"github.com/open-feature/cli/internal/filesystem"
	"github.com/open-feature/cli/internal/flagset"
)

// Represents the stability level of a generator
type Stability string

const (
	Alpha  Stability = "alpha"
	Beta   Stability = "beta"
	Stable Stability = "stable"
)

type CommonGenerator struct {
	Flagset *flagset.Flagset
}

type Params[T any] struct {
	OutputPath string
	Custom     T
}

type TemplateData struct {
	CommonGenerator
	Params[any]
}

// NewGenerator creates a new generator
func NewGenerator(flagset *flagset.Flagset, UnsupportedFlagTypes map[flagset.FlagType]bool) *CommonGenerator {
	return &CommonGenerator{
		Flagset: flagset.Filter(UnsupportedFlagTypes),
	}
}

func (g *CommonGenerator) GenerateFile(customFunc template.FuncMap, tmpl string, params *Params[any], name string) error {
	funcs := defaultFuncs()
	maps.Copy(funcs, customFunc)

	generatorTemplate, err := template.New("generator").Funcs(funcs).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("error initializing template: %v", err)
	}

	var buf bytes.Buffer
	data := TemplateData{
		CommonGenerator: *g,
		Params:          *params,
	}
	if err := generatorTemplate.Execute(&buf, data); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	return filesystem.WriteFile(filepath.Join(params.OutputPath, name), buf.Bytes())
}
