package config

import (
	"github.com/spf13/cobra"
)

// Flag name constants to avoid duplication
const (
	DebugFlagName       = "debug"
	ManifestFlagName    = "manifest"
	OutputFlagName      = "output"
	NoInputFlagName     = "no-input"
	GoPackageFlagName   = "package-name"
	CSharpNamespaceName = "namespace"
	OverrideFlagName    = "override"
)

// Default values for flags
const (
	DefaultManifestPath    = "flags.json"
	DefaultOutputPath      = ""
	DefaultGoPackageName   = "openfeature"
	DefaultCSharpNamespace = "OpenFeature"
)

// AddRootFlags adds the common flags to the given command
func AddRootFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(ManifestFlagName, "m", DefaultManifestPath, "Path to the flag manifest")
	cmd.PersistentFlags().Bool(NoInputFlagName, false, "Disable interactive prompts")
	cmd.PersistentFlags().Bool(DebugFlagName, false, "Enable debug logging")
}

// AddGenerateFlags adds the common generate flags to the given command
func AddGenerateFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(OutputFlagName, "o", DefaultOutputPath, "Path to where the generated files should be saved")
}

// AddGoGenerateFlags adds the go generator specific flags to the given command
func AddGoGenerateFlags(cmd *cobra.Command) {
	cmd.Flags().String(GoPackageFlagName, DefaultGoPackageName, "Name of the generated Go package")
}

// AddCSharpGenerateFlags adds the C# generator specific flags to the given command
func AddCSharpGenerateFlags(cmd *cobra.Command) {
	cmd.Flags().String(CSharpNamespaceName, DefaultCSharpNamespace, "Namespace for the generated C# code")
}

// AddInitFlags adds the init command specific flags
func AddInitFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(OverrideFlagName, false, "Override an existing configuration")
}

// GetManifestPath gets the manifest path from the given command
func GetManifestPath(cmd *cobra.Command) string {
	manifestPath, _ := cmd.Flags().GetString(ManifestFlagName)
	return manifestPath
}

// GetOutputPath gets the output path from the given command
func GetOutputPath(cmd *cobra.Command) string {
	outputPath, _ := cmd.Flags().GetString(OutputFlagName)
	return outputPath
}

// GetGoPackageName gets the Go package name from the given command
func GetGoPackageName(cmd *cobra.Command) string {
	goPackageName, _ := cmd.Flags().GetString(GoPackageFlagName)
	return goPackageName
}

// GetCSharpNamespace gets the C# namespace from the given command
func GetCSharpNamespace(cmd *cobra.Command) string {
	namespace, _ := cmd.Flags().GetString(CSharpNamespaceName)
	return namespace
}

// GetNoInput gets the no-input flag from the given command
func GetNoInput(cmd *cobra.Command) bool {
	noInput, _ := cmd.Flags().GetBool(NoInputFlagName)
	return noInput
}

// GetOverride gets the override flag from the given command
func GetOverride(cmd *cobra.Command) bool {
	override, _ := cmd.Flags().GetBool(OverrideFlagName)
	return override
}
