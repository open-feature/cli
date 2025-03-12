package cmd

import (
	"fmt"
	"os"

	"github.com/open-feature/cli/internal/config"
	"github.com/pterm/pterm"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	Version = "dev"
	Commit  string
	Date    string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string, commit string, date string) {
	Version = version
	Commit = commit
	Date = date
	if err := GetRootCmd().Execute(); err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}

func GetRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "openfeature",
		Short: "CLI for OpenFeature.",
		Long:  `CLI for OpenFeature related functionalities.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			printBanner()
			pterm.Println()
			pterm.Println("To see all the options, try 'openfeature --help'")
			pterm.Println()

			return nil
		},
		// Custom error handling for invalid commands
		SilenceErrors: true,
		SilenceUsage:  true,
		// Handle unknown commands
		DisableSuggestions:         false,
		SuggestionsMinimumDistance: 2,
		DisableAutoGenTag:          true,
	}

	// Add global flags using the config package
	config.AddRootFlags(rootCmd)

	// Add subcommands
	rootCmd.AddCommand(GetVersionCmd())
	rootCmd.AddCommand(GetInitCmd())
	rootCmd.AddCommand(GetGenerateCmd())
	rootCmd.AddCommand(GetPullCmd())

	// Add a custom error handler after the command is created
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		pterm.Error.Printf("Invalid flag: %s\n", err)
		pterm.Println("Run 'openfeature --help' for usage information")
		return err
	})

	return rootCmd
}

func initializeConfig(cmd *cobra.Command, bindPrefix string) error {
	v := viper.New()

	// Set the base name of the config file, without the file extension.
	v.SetConfigName(".openfeature")

	// Set as many paths as you like where viper should look for the
	// config file. We are only looking in the current working directory.
	v.AddConfigPath(".")

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// Bind the current command's flags to viper
	bindFlags(cmd, v, bindPrefix)

	return nil
}

// Bind each cobra flag to its associated viper configuration
func bindFlags(cmd *cobra.Command, v *viper.Viper, bindPrefix string) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Determine the naming convention of the flags when represented in the config file
		configName := f.Name
		if bindPrefix != "" {
			configName = bindPrefix + "." + f.Name
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				pterm.Error.Println(err)
			}
		}
	})
}
