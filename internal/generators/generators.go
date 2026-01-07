package generators

import (
	"bytes"
	"fmt"
	"maps"
	"path/filepath"
	"text/template"

	"github.com/open-feature/cli/internal/filesystem"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/logger"
)

// Represents the stability level of a generator
type Stability string

const (
	Alpha  Stability = "alpha"
	Beta   Stability = "beta"
	Stable Stability = "stable"
)

type CommonGenerator struct {
	Flagset   *flagset.Flagset
	Formatter func([]byte) ([]byte, error)
}

type Params[T any] struct {
	OutputPath   string
	TemplatePath string
	Custom       T
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

	// If a custom template path is provided, read from file instead of using embedded template
	if params.TemplatePath != "" {
		logger.Default.Debug(fmt.Sprintf("Using custom template: %s", params.TemplatePath))
		content, err := filesystem.ReadFile(params.TemplatePath)
		if err != nil {
			return fmt.Errorf("error reading custom template %s: %w", params.TemplatePath, err)
		}
		tmpl = string(content)
	}

	logger.Default.Debug(fmt.Sprintf("Generating file: %s", name))

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

	output := buf.Bytes()
	if g.Formatter != nil {
		output, err = g.Formatter(output)
		if err != nil {
			return fmt.Errorf("error executing formatter: %w", err)
		}
	}

	fullPath := filepath.Join(params.OutputPath, name)
	if err := filesystem.WriteFile(fullPath, output); err != nil {
		logger.Default.FileFailed(fullPath, err)
		return err
	}

	// Log successful file creation
	logger.Default.FileCreated(fullPath)

	return nil
}
