package cmd

import (
	"fmt"

	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators/golang"
	"github.com/open-feature/cli/internal/generators/react"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func GetGenerateReactCmd() *cobra.Command {
	reactCmd := &cobra.Command{
		Use:   "react",
		Short: "Generate typesafe React Hooks.",
		Long:  `Generate typesafe React Hooks compatible with the OpenFeature React SDK.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := config.GetString(config.ManifestFlag)
			flagset, err := flagset.Load(manifestPath)
			if err != nil {
				return err
			}
	
			generator := react.NewGenerator(flagset)
			err = generator.Generate()
			if err != nil {
				return err
			}

			pterm.Success.Println("Generated React Hooks.")
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
		// PreRun executes before flags are validated
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("Generating Golang flag accessors...")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			params := golang.Params{
				GoPackage: config.GetString("package-name"),
			}
			flagset, err := flagset.Load("sample/sample_manifest.json")
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

	goCmd.Flags().String("package-name", "", "Name of the Go package to be generated.")
	goCmd.MarkFlagRequired("package-name")

	return goCmd
}

func GetGenerateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate typesafe OpenFeature accessors.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO print overview of help message
			return nil
		},
	}

	generateCmd.AddCommand(GetGenerateReactCmd())
	generateCmd.AddCommand(GetGenerateGoCmd())

	return generateCmd
}
