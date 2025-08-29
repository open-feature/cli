package manifest

import (
	"fmt"
	"sort"
	"strings"

	schema "github.com/open-feature/cli/schema/v0"
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

func FormatValidationError(issues []ValidationError) string {
	var sb strings.Builder
	sb.WriteString("flag manifest validation failed:\n\n")

	// Group messages by flag path
	grouped := make(map[string]struct {
		flagType string
		messages []string
	})

	for _, issue := range issues {
		entry := grouped[issue.Path]
		entry.flagType = issue.Type
		entry.messages = append(entry.messages, issue.Message)
		grouped[issue.Path] = entry
	}

	// Sort paths for consistent output
	paths := make([]string, 0, len(grouped))
	for path := range grouped {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	// Format each row
	for _, path := range paths {
		entry := grouped[path]
		flagType := entry.flagType
		if flagType == "" {
			flagType = "missing"
		}
		sb.WriteString(fmt.Sprintf(
			"- flagType: %s\n  flagPath: %s\n  errors:\n    ~ %s\n  \tSuggestions:\n      \t- flagType: boolean\n      \t- defaultValue: true\n\n",
			flagType,
			path,
			strings.Join(entry.messages, "\n    ~ "),
		))
	}
	return sb.String()
}
