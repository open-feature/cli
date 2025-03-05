package golang

import (
	"github.com/open-feature/cli/internal/flagkeys"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators/golang"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Cmd for 'generate' command, handling code generation for flag accessors
var Cmd = &cobra.Command{
	Use:   "go",
	Short: "Generate Golang flag accessors for OpenFeature.",
	Long:  `Generate Golang flag accessors for OpenFeature.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := golang.Params{
			GoPackage: viper.GetString(flagkeys.GoPackageName),
		}
		flagset, err := flagset.Load("sample/sample_manifest.json")
		if err != nil {
			return err
		}

		generator := golang.NewGenerator(flagset)
		return generator.Generate(&params)
	},
}

func init() {
	Cmd.Flags().String(flagkeys.GoPackageName, "", "Name of the Go package to be generated.")
	Cmd.MarkFlagRequired(flagkeys.GoPackageName)
	viper.BindPFlag(flagkeys.GoPackageName, Cmd.Flags().Lookup(flagkeys.GoPackageName))

}
