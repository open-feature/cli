package cmd

import (
	"os"

	"github.com/open-feature/cli/internal/config"
	"github.com/pterm/pterm"

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
		pterm.Error.Println(err)
		os.Exit(1)
	}
}

func GetRootCmd() *cobra.Command {
	// Execute all parent's persistent hooks
	cobra.EnableTraverseRunHooks =true

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
		SilenceErrors: true,
		SilenceUsage:  true,
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

	// Add a custom error handler after the command is created
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		pterm.Error.Printf("Invalid flag: %s", err)
		pterm.Println("Run 'openfeature --help' for usage information")
		return err
	})

	return rootCmd
}
