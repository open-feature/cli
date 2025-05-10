package manifest

// OutputFormat represents the available output formats for the compare command
type OutputFormat string

const (
	// OutputFormatTree represents the tree output format (default)
	OutputFormatTree OutputFormat = "tree"
	// OutputFormatFlat represents the flat output format
	OutputFormatFlat OutputFormat = "flat"
	// OutputFormatJSON represents the JSON output format
	OutputFormatJSON OutputFormat = "json"
	// OutputFormatYAML represents the YAML output format
	OutputFormatYAML OutputFormat = "yaml"
)

// IsValidOutputFormat checks if the given format is a valid output format
func IsValidOutputFormat(format string) bool {
	switch OutputFormat(format) {
	case OutputFormatTree, OutputFormatFlat, OutputFormatJSON, OutputFormatYAML:
		return true
	default:
		return false
	}
}

// GetValidOutputFormats returns a list of all valid output formats
func GetValidOutputFormats() []string {
	return []string{
		string(OutputFormatTree),
		string(OutputFormatFlat),
		string(OutputFormatJSON),
		string(OutputFormatYAML),
	}
}