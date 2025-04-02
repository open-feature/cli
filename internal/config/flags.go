package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Flag name constants to avoid duplication
const (
	ManifestFlagName  = "manifest"
	OutputFlagName    = "output"
	NoInputFlagName   = "no-input"
	GoPackageFlagName = "package-name"
	OverrideFlagName  = "override"
	FlagSourceUrlFlagName = "flag-source-url"
	AuthTokenFlagName = "auth-token"
)

// Default values for flags
const (
	DefaultManifestPath  = "flags.json"
	DefaultOutputPath    = ""
	DefaultGoPackageName = "openfeature"
)

// AddRootFlags adds the common flags to the given command
func AddRootFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(ManifestFlagName, "m", DefaultManifestPath, "Path to the flag manifest")
	cmd.PersistentFlags().Bool(NoInputFlagName, false, "Disable interactive prompts")
}

// AddGenerateFlags adds the common generate flags to the given command
func AddGenerateFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(OutputFlagName, "o", DefaultOutputPath, "Path to where the generated files should be saved")
}

// AddGoGenerateFlags adds the go generator specific flags to the given command
func AddGoGenerateFlags(cmd *cobra.Command) {
	cmd.Flags().String(GoPackageFlagName, DefaultGoPackageName, "Name of the generated Go package")
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
