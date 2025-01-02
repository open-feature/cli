package golang

import (
	"github.com/open-feature/cli/internal/flagkeys"
	"github.com/open-feature/cli/internal/generate"
	"github.com/open-feature/cli/internal/generate/plugins/golang"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Cmd for `generate go` command, handling code generation for flag accessors
// for Golang.
var Cmd = &cobra.Command{
	Use:   "go",
	Short: "Generate Golang flag accessors for OpenFeature.",
	Long:  `Generate Golang flag accessors for OpenFeature.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := golang.Params{
			GoPackage: viper.GetString(flagkeys.GoPackageName),
		}
		gen := golang.NewGenerator(params)
		err := generate.CreateFlagAccessors(gen)
		return err
	},
}

func init() {
	Cmd.Flags().String(flagkeys.GoPackageName, "", "Name of the Go package to be generated.")
	Cmd.MarkFlagRequired(flagkeys.GoPackageName)
	viper.BindPFlag(flagkeys.GoPackageName, Cmd.Flags().Lookup(flagkeys.GoPackageName))

}
