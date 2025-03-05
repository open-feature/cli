package generate

import (
	"github.com/open-feature/cli/cmd/generate/golang"
	"github.com/open-feature/cli/cmd/generate/react"

	"github.com/spf13/cobra"
)

// Root for `generate“ sub-commands, handling code generation for flag accessors.
var Root = &cobra.Command{
	Use:   "generate",
	Short: "Code generation for flag accessors for OpenFeature.",
	Long:  `Code generation for flag accessors for OpenFeature.`,
}

func init() {
	// Add subcommands.
	Root.AddCommand(golang.Cmd)
	Root.AddCommand(react.Cmd)

	// Add flags.
	// Root.PersistentFlags().String(flagkeys.FlagManifestPath, "", "Path to the flag manifest.")
	// Root.MarkPersistentFlagRequired(flagkeys.FlagManifestPath)
	// viper.BindPFlag(flagkeys.FlagManifestPath, Root.PersistentFlags().Lookup(flagkeys.FlagManifestPath))
	// Root.PersistentFlags().String(flagkeys.OutputPath, "", "Output path for the codegen")
	// viper.BindPFlag(flagkeys.OutputPath, Root.PersistentFlags().Lookup(flagkeys.OutputPath))
	// Root.MarkPersistentFlagRequired(flagkeys.OutputPath)
}
