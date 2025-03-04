package generators

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/open-feature/cli/internal/filesystem"
	"github.com/open-feature/cli/internal/manifest"
)

type Flag struct {
	Key string
	Type string
	Description string
	DefaultValue any
}

type CommonGenerator struct {
	Stability Stability
	Flags []Flag
}

type CommonParams struct {
	OutputPath string
	Params map[string]any
}

type Options func(*CommonGenerator)

func (g *CommonGenerator) GetStability() Stability {
	if (g.Stability == "") {
		return Unknown
	}
	return g.Stability
}

func (g *CommonGenerator) SupportedFlagTypes() map[manifest.FlagType]bool {
	supportedTypes := map[manifest.FlagType]bool{}
	for _, flagType := range manifest.AllFlagTypes {
		supportedTypes[flagType] = true
	}
	return supportedTypes
}

func (g *CommonGenerator) GenerateFile(funcs template.FuncMap, t string, params CommonParams) error {
	// outputPath := params.OutputPath
	outputPath := "./blah.ts"
	generatorTemplate, err := template.New("generator").Funcs(funcs).Parse(t)
	if err != nil {
		return fmt.Errorf("error initializing template: %v", err)
	}

	var buf bytes.Buffer
	if err := generatorTemplate.Execute(&buf, g); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	fs := filesystem.FileSystem()
	if err := fs.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		return err
	}
	f, err := fs.Create(path.Join(outputPath))
	if err != nil {
		return fmt.Errorf("error creating file %q: %v", outputPath, err)
	}
	defer f.Close()
	writtenBytes, err := f.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("error writing contents to file %q: %v", outputPath, err)
	}
	if writtenBytes != buf.Len() {
		return fmt.Errorf("error writing entire file %v: writtenBytes != expectedWrittenBytes", outputPath)
	}

	return nil
}

func WithStability(stability Stability) Options {
	return func(g *CommonGenerator) {
		g.Stability = stability
	}
}

func NewCommonGenerator(manifest *manifest.Manifest, options ...Options) *CommonGenerator {
	commonGenerator := &CommonGenerator{
		Flags: manifestToFlag(manifest),
	}
	for _, option := range options {
		option(commonGenerator)
	}
	return commonGenerator
}

func manifestToFlag(manifest *manifest.Manifest) []Flag {
	flags := []Flag{}
	for key, flag := range manifest.Flags {
		flagData := flag.(map[string]any)
		flagType := flagData["flagType"].(string)
		description := flagData["description"].(string)
		defaultValue := flagData["defaultValue"]
		
		flags = append(flags, Flag{
			Key:          key,
			Type:         flagType,
			Description:  description,
			DefaultValue: defaultValue,
		})
	}

	// Ensure consistency of order of flag generation.
	sort.Slice(flags, func(i, j int) bool {
		return flags[i].Key < flags[j].Key
	})

	return flags
}