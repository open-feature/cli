package manifest

import (
	"fmt"
	"strings"

	"github.com/open-feature/cli/schema/v0"
	"github.com/xeipuuv/gojsonschema"
)

type ValidationError struct {
	Type    string `json:"type"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

func Validate(data []byte) ([]ValidationError, error) {
	schemaLoader := gojsonschema.NewStringLoader(schema.SchemaFile)
	manifestLoader := gojsonschema.NewBytesLoader(data)

	result, err := gojsonschema.Validate(schemaLoader, manifestLoader)
	if err != nil {
		return nil, fmt.Errorf("failed to validate manifest: %w", err)
	}

	var issues []ValidationError
	for _, err := range result.Errors() {
		if strings.HasPrefix(err.Field(), "flags") && err.Type() == "number_one_of" {
			issues = append(issues, ValidationError{
				Type:    err.Type(),
				Path:    err.Field(),
				Message: "flagType must be 'boolean', 'string', 'integer', 'float', or 'object'",
			})
		} else {
			issues = append(issues, ValidationError{
				Type:    err.Type(),
				Path:    err.Field(),
				Message: err.Description(),
			})
		}
	}

	return issues, nil
}
