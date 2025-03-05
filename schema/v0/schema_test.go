package schema

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xeipuuv/gojsonschema"
)

func TestPositiveFlagManifest(t *testing.T) {
	if err := walkPath(true, "./testdata/positive"); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestNegativeFlagManifest(t *testing.T) {
	if err := walkPath(false, "./testdata/negative"); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func walkPath(shouldPass bool, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ps := strings.Split(path, ".")
		if ps[len(ps)-1] != "json" {
			return nil
		}

		file, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var v any
		if err := json.Unmarshal([]byte(file), &v); err != nil {
			log.Fatal(err)
		}

		schemaLoader := gojsonschema.NewStringLoader(SchemaFile)
		manifestLoader := gojsonschema.NewGoLoader(v)
		result, err := gojsonschema.Validate(schemaLoader, manifestLoader)
		if (err != nil) {
			return fmt.Errorf("Error validating json schema: %v", err)
		}

		if (len(result.Errors()) >= 1 && shouldPass == true) {
			var errorMessage strings.Builder

			errorMessage.WriteString("file " + path + " should be valid, but had the following issues:\n")
			for _, error := range result.Errors() {
				errorMessage.WriteString(" - " + error.String() + "\n")
			}
			return fmt.Errorf("%s", errorMessage.String())
		}

		if (len(result.Errors()) == 0 && shouldPass == false) {
			return fmt.Errorf("file %s should be invalid, but no issues were detected", path)
		}

		return nil
	})
}
