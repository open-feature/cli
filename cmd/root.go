package cmd

import (
	"os"

	"github.com/open-feature/cli/internal/config"
	"github.com/pterm/pterm"
	"golang.org/x/term"

	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Commit  string
	Date    string

	noInput bool
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string, commit string, date string) {
	Version = version
	Commit = commit
	Date = date
	if err := GetRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func GetRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "openfeature",
		Short: "CLI for OpenFeature.",
		Long:  `CLI for OpenFeature related functionalities.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// TODO support custom path for config
			err := config.Load("", cmd)
			if err != nil {
				pterm.Error.Printf("%s", err)
			}

			// Check if the terminal is interactive
			// Ref: https://clig.dev/#interactivity
			if (!term.IsTerminal(int(os.Stdin.Fd())) || config.GetBool(config.NoInputFlag)) {
				pterm.Debug.Println("Disabling interactive prompts")
				noInput = true
			}

			return err
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
		DisableSuggestions: false,
		SuggestionsMinimumDistance: 2,
		DisableAutoGenTag: true,
	}

	// Add global flags
	// rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to the configuration (defaults to .openfeature.yaml)")
	rootCmd.PersistentFlags().StringP("manifest", "m", "flags.json", "Path to the flag manifest")
	rootCmd.PersistentFlags().Bool(config.NoInputFlag, false, "Disable interactive prompts")
	// viper.BindPFlag(("manifest"), rootCmd.PersistentFlags().Lookup("manifest"))


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
