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
	Unknown Stability = "unknown"
	Alpha   Stability = "alpha"
	Beta    Stability = "beta"
	Stable  Stability = "stable"
)

type CommonGenerator struct {
	Stability            Stability
	UnsupportedFlagTypes map[flagset.FlagType]bool
	Flagset              *flagset.Flagset
}

type Params[T any] struct {
	OutputPath string
	Custom     T
}

type TemplateData struct {
	CommonGenerator
	Params[any]
}

type Options func(*CommonGenerator)

func WithStability(stability Stability) Options {
	return func(g *CommonGenerator) {
		g.Stability = stability
	}
}

func WithUnsupportedFlagType(flagType flagset.FlagType) Options {
	return func(g *CommonGenerator) {
		if g.UnsupportedFlagTypes == nil {
			g.UnsupportedFlagTypes = make(map[flagset.FlagType]bool)
		}
		g.UnsupportedFlagTypes[flagType] = true
	}
}

func NewCommonGenerator(flagset *flagset.Flagset, options ...Options) *CommonGenerator {
	commonGenerator := &CommonGenerator{}
	for _, option := range options {
		option(commonGenerator)
	}
	commonGenerator.Flagset = flagset.Filter(commonGenerator.UnsupportedFlagTypes)
	return commonGenerator
}

func (g *CommonGenerator) GetStability() Stability {
	if g.Stability == "" {
		return Unknown
	}
	return g.Stability
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
