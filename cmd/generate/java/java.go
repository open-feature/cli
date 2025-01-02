package java

import (
	"github.com/open-feature/cli/internal/flagkeys"
	"github.com/open-feature/cli/internal/generate"
	"github.com/open-feature/cli/internal/generate/plugins/java"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Cmd for `generate java“ command, handling code generation for flag accessors
// for Java.
var Cmd = &cobra.Command{
	Use:   "java",
	Short: "Generate Java flag accessors for OpenFeature.",
	Long:  `Generate Java flag accessors for OpenFeature.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		params := java.Params{
			JavaPackage: viper.GetString(flagkeys.JavaPackageName),
		}
		gen := java.NewGenerator(params)
		err := generate.CreateFlagAccessors(gen)
		return err
	},
}

func init() {
	Cmd.Flags().String(flagkeys.JavaPackageName, "", "Name of the Java package to be generated.")
	Cmd.MarkFlagRequired(flagkeys.JavaPackageName)
	viper.BindPFlag(flagkeys.JavaPackageName, Cmd.Flags().Lookup(flagkeys.JavaPackageName))

}
