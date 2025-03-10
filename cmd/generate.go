package cmd

import (
	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators"
	"github.com/open-feature/cli/internal/generators/golang"
	"github.com/open-feature/cli/internal/generators/react"
	"github.com/spf13/cobra"
)

func GetGenerateReactCmd() *cobra.Command {
	reactCmd := &cobra.Command{
		Use:   "react",
		Short: "Generate typesafe React Hooks.",
		Long:  `Generate typesafe React Hooks compatible with the OpenFeature React SDK.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "generate.react")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := config.GetManifestPath(cmd)
			outputPath := config.GetOutputPath(cmd)

			params := generators.Params[react.Params]{
				OutputPath: outputPath,
				Custom: react.Params{},
			}
			flagset, err := flagset.Load(manifestPath)
			if err != nil {
				return err
			}

			generator := react.NewGenerator(flagset)
			err = generator.Generate(&params)
			if err != nil {
				return err
			}
			return nil
		},
	}

	return reactCmd
}

func GetGenerateGoCmd() *cobra.Command {
	goCmd := &cobra.Command{
		Use:   "go",
		Short: "Generate typesafe accessors for OpenFeature.",
		Long:  `Generate typesafe accessors compatible with the OpenFeature Go SDK.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "generate.go")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
				// Use the helper functions to get flag values
			goPackageName := config.GetGoPackageName(cmd)
			manifestPath := config.GetManifestPath(cmd)
			outputPath := config.GetOutputPath(cmd)

			params := generators.Params[golang.Params]{
				OutputPath: outputPath,
				Custom: golang.Params{
					GoPackage: goPackageName,
				},
			}

			flagset, err := flagset.Load(manifestPath)
			if err != nil {
				return err
			}

			generator := golang.NewGenerator(flagset)
			err = generator.Generate(&params)
			if err != nil {
				return err
			}
			return nil
		},
	}

	// Add Go-specific flags
	config.AddGoGenerateFlags(goCmd)

	return goCmd
}

func GetGenerateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate typesafe OpenFeature accessors.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// If this command has a parent with PersistentPreRunE, call it
			if cmd.Parent() != nil && cmd.Parent().PersistentPreRunE != nil {
				err := cmd.Parent().PersistentPreRunE(cmd.Parent(), args)
				if err != nil {
					return err
				}
			}

			return initializeConfig(cmd, "generate")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO print overview of help message
			return nil
		},
	}

	// Add generate flags using the config package
	config.AddGenerateFlags(generateCmd)

	// Add generate subcommands
	generateCmd.AddCommand(GetGenerateReactCmd())
	generateCmd.AddCommand(GetGenerateGoCmd())

	return generateCmd
}
