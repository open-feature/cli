package cmd

import (
	"fmt"
	"os"

	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/logger"

	"github.com/spf13/cobra"
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
		logger.Default.Error(err.Error())
		os.Exit(1)
	}
}

func GetRootCmd() *cobra.Command {
	// Execute all parent's persistent hooks
	cobra.EnableTraverseRunHooks = true

	rootCmd := &cobra.Command{
		Use:   "openfeature",
		Short: "CLI for OpenFeature.",
		Long:  `CLI for OpenFeature related functionalities.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			debug, _ := cmd.Flags().GetBool("debug")
			logger.Default.SetDebug(debug)
			logger.Default.Debug("Debug logging enabled")
			return initializeConfig(cmd, "")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			printBanner()
			logger.Default.Println("")
			logger.Default.Println("To see all the options, try 'openfeature --help'")
			return nil
		},
		SilenceErrors:              true,
		SilenceUsage:               true,
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
	rootCmd.AddCommand(GetCompareCmd())

	// Add a custom error handler after the command is created
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		logger.Default.Error(fmt.Sprintf("Invalid flag: %s", err))
		logger.Default.Info("Run 'openfeature --help' for usage information")
		return err
	})

	return rootCmd
}
