package generators

import (
	"bytes"
	"fmt"
	"text/template"

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

type Params struct {
	OutputPath string
	Custom map[string]any
}

type TemplateData struct {
	CommonGenerator
	Params
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

func (g *CommonGenerator) GenerateFile(customFunc template.FuncMap, t string, params *Params) error {
	funcs := defaultFuncs()
	for k, v := range customFunc {
		funcs[k] = v
	}

	generatorTemplate, err := template.New("generator").Funcs(funcs).Parse(t)
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

	return filesystem.WriteFile(params.OutputPath, buf.Bytes())
}
