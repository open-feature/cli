package manifest

import (
	"bytes"
	"encoding/json"
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

	// Check for duplicate flag keys
	duplicates := findDuplicateFlagKeys(data)
	for _, key := range duplicates {
		issues = append(issues, ValidationError{
			Type:    "duplicate_key",
			Path:    fmt.Sprintf("flags.%s", key),
			Message: fmt.Sprintf("flag '%s' is defined multiple times in the manifest", key),
		})
	}

	return issues, nil
}

// findDuplicateFlagKeys parses the raw JSON to detect duplicate keys within the "flags" object.
// Standard JSON unmarshaling silently accepts duplicates (taking the last value), so we use
// a token-based approach to detect them.
func findDuplicateFlagKeys(data []byte) []string {
	decoder := json.NewDecoder(bytes.NewReader(data))

	// Navigate to the root object
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		return nil
	}

	// Look for the "flags" key at the top level
	for decoder.More() {
		keyToken, err := decoder.Token()
		if err != nil {
			return nil
		}

		key, ok := keyToken.(string)
		if !ok {
			continue
		}

		if key == "flags" {
			return findDuplicatesInObject(decoder)
		}

		// Skip the value for non-"flags" keys
		skipValue(decoder)
	}

	return nil
}

// findDuplicatesInObject reads an object from the decoder and returns any duplicate keys.
func findDuplicatesInObject(decoder *json.Decoder) []string {
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		return nil
	}

	seen := make(map[string]bool)
	var duplicates []string

	for decoder.More() {
		keyToken, err := decoder.Token()
		if err != nil {
			break
		}

		key, ok := keyToken.(string)
		if !ok {
			continue
		}

		if seen[key] {
			duplicates = append(duplicates, key)
		} else {
			seen[key] = true
		}

		// Skip the value
		skipValue(decoder)
	}

	// Consume the closing brace
	_, err = decoder.Token()
	if err != nil {
		return duplicates
	}

	// Sort for consistent output
	sort.Strings(duplicates)

	return duplicates
}

// skipValue advances the decoder past one complete JSON value (object, array, or primitive).
func skipValue(decoder *json.Decoder) {
	token, err := decoder.Token()
	if err != nil {
		return
	}

	switch t := token.(type) {
	case json.Delim:
		switch t {
		case '{':
			// Skip object contents
			for decoder.More() {
				if _, err := decoder.Token(); err != nil { // key
					return
				}
				skipValue(decoder)
			}
			if _, err := decoder.Token(); err != nil { // closing }
				return
			}
		case '[':
			// Skip array contents
			for decoder.More() {
				skipValue(decoder)
			}
			if _, err := decoder.Token(); err != nil { // closing ]
				return
			}
		}
	}
	// Primitives (string, number, bool, null) are already consumed by the Token() call
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
