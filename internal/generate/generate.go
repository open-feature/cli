// Package generate contains the top level functions used for generating flag accessors.
package generate

import (
	"bytes"
	"codegen/internal/flagkeys"
	"codegen/internal/generate/manifestutils"
	"codegen/internal/generate/types"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"text/template"

	"github.com/spf13/viper"
)

// GenerateFile receives data for the Go template engine and outputs the contents to the file.
// Intended to be invoked by each language generator with appropriate data.
func GenerateFile(funcs template.FuncMap, contents string, data types.TmplDataInterface) error {
	contentsTmpl, err := template.New("contents").Funcs(funcs).Parse(contents)
	if err != nil {
		return fmt.Errorf("error initializing template: %v", err)
	}

	var buf bytes.Buffer
	if err := contentsTmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}
	outputPath := data.BaseTmplDataInfo().OutputPath
	if err := os.MkdirAll(filepath.Dir(outputPath), 0770); err != nil {
		return err
	}
	f, err := os.Create(path.Join(outputPath))
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

// Takes as input a generator and outputs file with the appropriate flag accessors.
// The flag data is taken from the provided flag manifest.
func CreateFlagAccessors(gen types.Generator) error {
	bt, err := manifestutils.LoadData(viper.GetString(flagkeys.FlagManifestPath), gen.SupportedFlagTypes())
	if err != nil {
		return fmt.Errorf("error loading flag manifest: %v", err)
	}
	input := types.Input{
		BaseData: bt,
	}
	return gen.Generate(input)
}
