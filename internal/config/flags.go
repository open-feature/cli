package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	FlagSourceUrlFlagName = "flag-source-url"
	AuthTokenFlagName     = "auth-token"
	NoPromptFlagName      = "no-prompt"
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
	cmd.Flags().String(FlagSourceUrlFlagName, "", "The URL of the flag source")
}

// AddPullFlags adds the pull command specific flags
func AddPullFlags(cmd *cobra.Command) {
	cmd.Flags().String(FlagSourceUrlFlagName, "", "The URL of the flag source")
	cmd.Flags().String(AuthTokenFlagName, "", "The auth token for the flag source")
	cmd.Flags().Bool(NoPromptFlagName, false, "Disable interactive prompts for missing default values")
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

// GetFlagSourceUrl gets the flag source URL from the given command
func GetFlagSourceUrl(cmd *cobra.Command) string {
	flagSourceUrl, _ := cmd.Flags().GetString(FlagSourceUrlFlagName)
	if flagSourceUrl == "" {
		viper.SetConfigName(".openfeature")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			return ""
		}
		if !viper.IsSet("flagSourceUrl") {
			return ""
		}
		flagSourceUrl = viper.GetString("flagSourceUrl")
	}
	return flagSourceUrl
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
