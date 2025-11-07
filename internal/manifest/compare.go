package manifest

import (
	"fmt"
	"maps"
	"path/filepath"
	"reflect"
	"strings"
)

type Change struct {
	Type     string `json:"type"`
	Path     string `json:"path"`
	OldValue any    `json:"oldValue,omitempty"`
	NewValue any    `json:"newValue,omitempty"`
}

// CompareOptions holds options for comparing manifests
type CompareOptions struct {
	IgnorePatterns []string
}

// Compare compares two manifests and returns differences, optionally ignoring specified fields
func Compare(oldManifest, newManifest *Manifest, opts CompareOptions) ([]Change, error) {
	var changes []Change
	oldFlags := oldManifest.Flags
	newFlags := newManifest.Flags

	// Check for changes and additions
	for key, newFlag := range newFlags {
		if oldFlag, exists := oldFlags[key]; exists {
			// Compare flags semantically with ignore patterns
			if flagHasChanges(oldFlag, newFlag, key, opts.IgnorePatterns) {
				changes = append(changes, Change{
					Type:     "change",
					Path:     fmt.Sprintf("flags.%s", key),
					OldValue: oldFlag,
					NewValue: newFlag,
				})
			}
		} else {
			changes = append(changes, Change{
				Type:     "add",
				Path:     fmt.Sprintf("flags.%s", key),
				NewValue: newFlag,
			})
		}
	}

	// Check for removals
	for key, oldFlag := range oldFlags {
		if _, exists := newFlags[key]; !exists {
			changes = append(changes, Change{
				Type:     "remove",
				Path:     fmt.Sprintf("flags.%s", key),
				OldValue: oldFlag,
			})
		}
	}

	return changes, nil
}

// getKnownFlagProperties returns the set of known schema properties for flags
// by extracting JSON field names from the BaseFlag struct
func getKnownFlagProperties() map[string]bool {
	props := make(map[string]bool)

	// Extract fields from BaseFlag struct
	t := reflect.TypeOf(BaseFlag{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			// Extract field name from json tag (e.g., "flagType,omitempty" -> "flagType")
			fieldName := strings.Split(jsonTag, ",")[0]
			if fieldName != "" && fieldName != "-" {
				props[fieldName] = true
			}
		}
	}

	// All flag types have defaultValue, but it's defined in the specific flag type structs
	// (BooleanFlag, StringFlag, etc.), not in BaseFlag, so we add it manually here.
	// This is necessary because each flag type has a different type for defaultValue
	// (bool, string, int, float64, any), so it can't be in the base struct.
	props["defaultValue"] = true

	return props
}

// Cache the known properties (initialized once)
var knownFlagProperties = getKnownFlagProperties()

// flagHasChanges checks if two flags have semantic differences, ignoring specified patterns
func flagHasChanges(oldFlag, newFlag any, flagKey string, ignorePatterns []string) bool {
	// Flatten both flags to path->value maps
	oldPaths := flattenToMap(oldFlag, fmt.Sprintf("flags.%s", flagKey))
	newPaths := flattenToMap(newFlag, fmt.Sprintf("flags.%s", flagKey))

	// Get all unique paths
	allPaths := make(map[string]bool)
	for path := range oldPaths {
		allPaths[path] = true
	}
	for path := range newPaths {
		allPaths[path] = true
	}

	// Check each path for differences (skipping ignored paths and unknown properties)
	for path := range allPaths {
		// Skip paths that should be ignored
		if shouldIgnorePath(path, ignorePatterns) {
			continue
		}

		// Skip unknown properties (only compare known schema properties)
		if !isKnownProperty(path, flagKey) {
			continue
		}

		oldVal, oldExists := oldPaths[path]
		newVal, newExists := newPaths[path]

		// If existence differs or values differ, we have a change
		if oldExists != newExists || !reflect.DeepEqual(oldVal, newVal) {
			return true
		}
	}

	return false
}

// isKnownProperty checks if a path represents a known schema property
func isKnownProperty(path, flagKey string) bool {
	// Extract the property name from the path
	// Path format: flags.<flagKey>.<property> or flags.<flagKey>.<nested>.<property>
	prefix := fmt.Sprintf("flags.%s.", flagKey)
	if !strings.HasPrefix(path, prefix) {
		return false
	}

	// Get the first property after the flag key
	remainder := strings.TrimPrefix(path, prefix)
	parts := strings.Split(remainder, ".")
	if len(parts) == 0 {
		return false
	}

	propertyName := parts[0]

	// Check if it's a known property
	return knownFlagProperties[propertyName]
}

// IsKnownPropertyForTest checks if a path is a known property (exported for testing/rendering)
func IsKnownPropertyForTest(path, flagKey string) bool {
	return isKnownProperty(path, flagKey)
}

// flattenToMap recursively flattens a nested structure into a map of path->value
func flattenToMap(obj any, prefix string) map[string]any {
	result := make(map[string]any)

	switch v := obj.(type) {
	case map[string]any:
		for key, value := range v {
			path := fmt.Sprintf("%s.%s", prefix, key)

			// If the value is a simple type, add it directly
			if isSimpleType(value) {
				result[path] = value
			} else {
				// Otherwise, recursively flatten
				nested := flattenToMap(value, path)
				maps.Copy(result, nested)
			}
		}
	case []any:
		// For arrays, include them as-is at their path
		result[prefix] = v
	default:
		// Simple value
		result[prefix] = v
	}

	return result
}

// isSimpleType checks if a value is a simple type (not a map or struct)
func isSimpleType(v any) bool {
	if v == nil {
		return true
	}

	switch v.(type) {
	case map[string]any:
		return false
	case []any:
		// Arrays are considered simple for comparison purposes
		return true
	default:
		return true
	}
}

// shouldIgnorePath checks if a path should be ignored based on the ignore patterns
func shouldIgnorePath(path string, ignorePatterns []string) bool {
	for _, pattern := range ignorePatterns {
		if matchesPattern(path, pattern) {
			return true
		}
	}
	return false
}

// ShouldIgnorePathForTest checks if a path should be ignored (exported for testing/rendering)
func ShouldIgnorePathForTest(path string, ignorePatterns []string) bool {
	return shouldIgnorePath(path, ignorePatterns)
}

// matchesPattern checks if a path matches a given pattern
// Supports:
// - Full paths with globs: flags.*.description
// - Shorthand: description (matches **.description anywhere)
// - Partial wildcards: metadata.* (matches any path ending with metadata.<something>)
func matchesPattern(path, pattern string) bool {
	// Check for exact match first
	if path == pattern {
		return true
	}

	// If pattern contains a dot, treat it as a path pattern
	if strings.Contains(pattern, ".") {
		// Try direct filepath.Match first
		matched, err := filepath.Match(pattern, path)
		if err == nil && matched {
			return true
		}

		// Try matching with wildcard support
		if matchesPathSegments(path, pattern) {
			return true
		}

		// For patterns like "metadata.*", check if path ends with that pattern
		// This allows "metadata.*" to match "flags.flagName.metadata.field"
		if strings.Contains(pattern, "*") {
			// Extract the suffix after the last non-wildcard part
			// e.g., "metadata.*" should match any path containing "metadata.<anything>"
			if matchesPartialPath(path, pattern) {
				return true
			}
		}
	} else {
		// Shorthand pattern - matches field anywhere in path
		// E.g., "description" matches "flags.myFlag.description"
		pathParts := strings.Split(path, ".")
		for _, part := range pathParts {
			if part == pattern {
				return true
			}
		}
	}

	return false
}

// matchesPartialPath checks if a path matches a pattern with wildcards anywhere in the path
// E.g., "metadata.*" matches "flags.flagName.metadata.author"
func matchesPartialPath(path, pattern string) bool {
	pathParts := strings.Split(path, ".")
	patternParts := strings.Split(pattern, ".")

	// Try to find a matching subsequence in the path
	for i := 0; i <= len(pathParts)-len(patternParts); i++ {
		if matchesSubsequence(pathParts[i:], patternParts) {
			return true
		}
	}

	return false
}

// matchesSubsequence checks if a path subsequence matches a pattern
func matchesSubsequence(pathParts, patternParts []string) bool {
	if len(pathParts) < len(patternParts) {
		return false
	}

	for i, patternPart := range patternParts {
		if patternPart == "*" {
			// Wildcard matches any single part
			continue
		}
		if pathParts[i] != patternPart {
			return false
		}
	}

	return true
}

// matchesPathSegments handles patterns with * wildcards matching any single segment
func matchesPathSegments(path, pattern string) bool {
	pathParts := strings.Split(path, ".")
	patternParts := strings.Split(pattern, ".")

	// If pattern has more parts than path, can't match
	if len(patternParts) > len(pathParts) {
		return false
	}

	// Match each pattern part against corresponding path part
	pathIdx := 0
	for patternIdx, patternPart := range patternParts {
		if patternPart == "*" {
			// Wildcard matches any single segment
			pathIdx++
		} else {
			// Exact match required
			if pathIdx >= len(pathParts) || pathParts[pathIdx] != patternPart {
				return false
			}
			pathIdx++
		}

		// If this is the last pattern part and it's *, and we've matched all path parts, success
		if patternIdx == len(patternParts)-1 && patternPart == "*" && pathIdx == len(pathParts) {
			return true
		}
	}

	// All pattern parts matched and we consumed all path parts
	return pathIdx == len(pathParts)
}
