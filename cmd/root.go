package cmd

import (
	"fmt"
	"os"

	"github.com/open-feature/cli/cmd/generate"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version = "dev"
	Commit  string
	Date    string

	ManifestPath = "flags.json"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "openfeature",
	Short: "CLI for OpenFeature.",
	Long:  `CLI for OpenFeature related functionalities.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		ManifestPath = viper.GetString("manifest-path")
		return nil
	},
	DisableAutoGenTag: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string, commit string, date string) {
	Version = version
	Commit = commit
	Date = date
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&ManifestPath, "manifest-path", "m", ManifestPath, "Specify the path and name for the manifest file")
	
	viper.BindPFlag("manifest-path", rootCmd.PersistentFlags().Lookup("manifest-path"))
	
	viper.SetConfigName(".openfeature")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	rootCmd.AddCommand(generate.Root)
	rootCmd.AddCommand(versionCmd)
}
