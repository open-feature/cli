package config

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// Flag name constants to avoid duplication
const (
	DebugFlagName         = "debug"
	ManifestFlagName      = "manifest"
	OutputFlagName        = "output"
	NoInputFlagName       = "no-input"
	GoPackageFlagName     = "package-name"
	CSharpNamespaceName   = "namespace"
	OverrideFlagName      = "override"
	JavaPackageFlagName   = "package-name"
	ProviderFlagName      = "provider"
	FlagSourceUrlFlagName = "flag-source-url" // Deprecated: use ProviderFlagName instead
	AuthTokenFlagName     = "auth-token"
	NoPromptFlagName      = "no-prompt"
	DryRunFlagName        = "dry-run"
	TypeFlagName          = "type"
	DefaultValueFlagName  = "default-value"
	DescriptionFlagName   = "description"
)

// Default values for flags
const (
	DefaultManifestPath    = "flags.json"
	DefaultOutputPath      = ""
	DefaultGoPackageName   = "openfeature"
	DefaultCSharpNamespace = "OpenFeature"
	DefaultJavaPackageName = "com.example.openfeature"
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

// AddJavaGenerateFlags adds the Java generator specific flags to the given command
func AddJavaGenerateFlags(cmd *cobra.Command) {
	cmd.Flags().String(JavaPackageFlagName, DefaultJavaPackageName, "Name of the generated Java package")
}

// AddInitFlags adds the init command specific flags
func AddInitFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(OverrideFlagName, false, "Override an existing configuration")
	cmd.Flags().String(ProviderFlagName, "", "The URL of the flag provider")
	cmd.Flags().String(FlagSourceUrlFlagName, "", "The URL of the flag source (deprecated: use --provider-url instead)")
	_ = cmd.Flags().MarkDeprecated(FlagSourceUrlFlagName, "use --provider-url instead")
}

// AddPullFlags adds the pull command specific flags
func AddPullFlags(cmd *cobra.Command) {
	cmd.Flags().String(ProviderFlagName, "", "The URL of the flag provider")
	cmd.Flags().String(FlagSourceUrlFlagName, "", "The URL of the flag source (deprecated: use --provider-url instead)")
	_ = cmd.Flags().MarkDeprecated(FlagSourceUrlFlagName, "use --provider-url instead")
	cmd.Flags().String(AuthTokenFlagName, "", "The auth token for the flag provider")
	cmd.Flags().Bool(NoPromptFlagName, false, "Disable interactive prompts for missing default values")
}

// AddPushFlags adds the push command specific flags
func AddPushFlags(cmd *cobra.Command) {
	cmd.Flags().String(ProviderFlagName, "", "The URL of the flag provider")
	cmd.Flags().String(FlagSourceUrlFlagName, "", "The URL of the flag destination (deprecated: use --provider-url instead)")
	_ = cmd.Flags().MarkDeprecated(FlagSourceUrlFlagName, "use --provider-url instead")
	cmd.Flags().String(AuthTokenFlagName, "", "The auth token for the flag provider")
	cmd.Flags().Bool(DryRunFlagName, false, "Preview changes without pushing")
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

// GetJavaPackageName gets the Java package name from the given command
func GetJavaPackageName(cmd *cobra.Command) string {
	javaPackageName, _ := cmd.Flags().GetString(JavaPackageFlagName)
	return javaPackageName
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

// getConfigValueWithFallback is a helper function that attempts to get a value from Viper config
// if the provided value is empty. This reduces duplication for flag source/destination URLs.
// It checks both the new and legacy config keys, with the new key taking precedence.
func getConfigValueWithFallback(value string, newConfigKey string, legacyConfigKey string) string {
	if value != "" {
		return value
	}
	// Check the new config key first
	if configValue := viper.GetString(newConfigKey); configValue != "" {
		return configValue
	}
	// Fall back to legacy config key for backward compatibility
	return viper.GetString(legacyConfigKey)
}

// GetFlagSourceUrl gets the flag source URL from the given command
// It checks the new --provider-url flag first, then falls back to the deprecated --flag-source-url flag
// for backward compatibility. Finally, it checks the config file for both keys.
func GetFlagSourceUrl(cmd *cobra.Command) string {
	// Check new flag first
	provider, _ := cmd.Flags().GetString(ProviderFlagName)
	if provider != "" {
		return provider
	}

	// Fall back to deprecated flag for backward compatibility
	flagSourceUrl, _ := cmd.Flags().GetString(FlagSourceUrlFlagName)

	// Use the fallback helper which checks both config keys
	return getConfigValueWithFallback(flagSourceUrl, "provider", "flagSourceUrl")
}

// GetAuthToken gets the auth token from the given command
func GetAuthToken(cmd *cobra.Command) string {
	authToken, _ := cmd.Flags().GetString(AuthTokenFlagName)
	return authToken
}

// GetNoPrompt gets the no-prompt flag from the given command
func GetNoPrompt(cmd *cobra.Command) bool {
	noPrompt, _ := cmd.Flags().GetBool(NoPromptFlagName)
	return noPrompt
}

// GetDryRun gets the dry-run flag from the given command
func GetDryRun(cmd *cobra.Command) bool {
	dryRun, _ := cmd.Flags().GetBool(DryRunFlagName)
	return dryRun
}

// AddManifestAddFlags adds the manifest add command specific flags
func AddManifestAddFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(TypeFlagName, "t", "boolean", "Type of the flag (boolean, string, integer, float, object)")
	cmd.Flags().StringP(DefaultValueFlagName, "d", "", "Default value for the flag (required)")
	cmd.Flags().String(DescriptionFlagName, "", "Description of the flag")
}

// AddManifestListFlags adds the manifest list command specific flags
func AddManifestListFlags(cmd *cobra.Command) {
	// Currently no specific flags for list command, but function exists for consistency
}

// ShouldDisableInteractivePrompts returns true if interactive prompts should be disabled
// This happens when:
// - The --no-input flag is set, OR
// - stdin is not a terminal (e.g., in tests, CI, or when input is piped)
func ShouldDisableInteractivePrompts(cmd *cobra.Command) bool {
	noInput := GetNoInput(cmd)
	if noInput {
		return true
	}
	// Automatically disable prompting if stdin is not a terminal
	return !term.IsTerminal(int(os.Stdin.Fd()))
}
