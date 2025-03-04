package react

import (
	// "github.com/open-feature/cli/internal/generate"
	// "github.com/open-feature/cli/internal/generate/plugins/react"

	"github.com/open-feature/cli/internal/generators/react"
	"github.com/open-feature/cli/internal/manifest"

	"github.com/spf13/cobra"
)

// Cmd for "generate" command, handling code generation for flag accessors
var Cmd = &cobra.Command{
	Use:   "react",
	Short: "Generate typesafe React Hooks.",
	Long:  `Generate typesafe React Hooks compatible with the OpenFeature React SDK.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		manifest, err := manifest.Load("sample/sample_manifest.json")
		if err != nil {
			return err
		}

		params := react.Params{}
		generator := react.NewGenerator(manifest)
		return generator.Generate(params)
	},
}

func init() {
}
