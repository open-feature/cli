package cmd

import (
	"strings"

	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/generators"
	"github.com/open-feature/cli/internal/generators/golang"
	"github.com/open-feature/cli/internal/generators/nodejs"
	"github.com/open-feature/cli/internal/generators/python"
	"github.com/open-feature/cli/internal/generators/react"
	"github.com/open-feature/cli/internal/logger"
	"github.com/spf13/cobra"
)

// addStabilityInfo adds stability information to the command's help template before "Usage:"
func addStabilityInfo(cmd *cobra.Command) {
	// Only modify commands that have a stability annotation
	if stability, ok := cmd.Annotations["stability"]; ok {
		originalTemplate := cmd.UsageTemplate()

		// Find the "Usage:" section and insert stability info before it
		if strings.Contains(originalTemplate, "Usage:") {
			customTemplate := strings.Replace(
				originalTemplate,
				"Usage:",
				"Stability: "+stability+"\n\nUsage:",
				1, // Replace only the first occurrence
			)
			cmd.SetUsageTemplate(customTemplate)
		} else {
			// Fallback if "Usage:" not found - prepend to the template
			customTemplate := "Stability: " + stability + "\n\n" + originalTemplate
			cmd.SetUsageTemplate(customTemplate)
		}
	}
}

func GetGenerateNodeJSCmd() *cobra.Command {
	nodeJSCmd := &cobra.Command{
		Use:   "nodejs",
		Short: "Generate typesafe Node.js client.",
		Long:  `Generate typesafe Node.js client compatible with the OpenFeature JavaScript Server SDK.`,
		Annotations: map[string]string{
			"stability": string(generators.Alpha),
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "generate.nodejs")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := config.GetManifestPath(cmd)
			outputPath := config.GetOutputPath(cmd)

			logger.Default.GenerationStarted("Node.js")

			params := generators.Params[nodejs.Params]{
				OutputPath: outputPath,
				Custom:     nodejs.Params{},
			}
			flagset, err := flagset.Load(manifestPath)
			if err != nil {
				return err
			}

			generator := nodejs.NewGenerator(flagset)
			logger.Default.Debug("Executing Node.js generator")
			err = generator.Generate(&params)
			if err != nil {
				return err
			}
	
			logger.Default.GenerationComplete("Node.js")
			
			return nil
		},
	}

	addStabilityInfo(nodeJSCmd)

	return nodeJSCmd
}

func GetGenerateReactCmd() *cobra.Command {
	reactCmd := &cobra.Command{
		Use:   "react",
		Short: "Generate typesafe React Hooks.",
		Long:  `Generate typesafe React Hooks compatible with the OpenFeature React SDK.`,
		Annotations: map[string]string{
			"stability": string(generators.Alpha),
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "generate.react")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := config.GetManifestPath(cmd)
			outputPath := config.GetOutputPath(cmd)
			
			logger.Default.GenerationStarted("React")

			params := generators.Params[react.Params]{
				OutputPath: outputPath,
				Custom:     react.Params{},
			}
			flagset, err := flagset.Load(manifestPath)
			if err != nil {
				return err
			}

			generator := react.NewGenerator(flagset)
			logger.Default.Debug("Executing React generator")
			err = generator.Generate(&params)
			if err != nil {
				return err
			}
			
			logger.Default.GenerationComplete("React")
			
			return nil
		},
	}

	addStabilityInfo(reactCmd)

	return reactCmd
}

func GetGenerateGoCmd() *cobra.Command {
	goCmd := &cobra.Command{
		Use:   "go",
		Short: "Generate typesafe accessors for OpenFeature.",
		Long:  `Generate typesafe accessors compatible with the OpenFeature Go SDK.`,
		Annotations: map[string]string{
			"stability": string(generators.Alpha),
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "generate.go")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			goPackageName := config.GetGoPackageName(cmd)
			manifestPath := config.GetManifestPath(cmd)
			outputPath := config.GetOutputPath(cmd)
			
			logger.Default.GenerationStarted("Go")

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
			logger.Default.Debug("Executing Go generator")
			err = generator.Generate(&params)
			if err != nil {
				return err
			}
			
			logger.Default.GenerationComplete("Go")
			
			return nil
		},
	}

	// Add Go-specific flags
	config.AddGoGenerateFlags(goCmd)

	addStabilityInfo(goCmd)

	return goCmd
}

func getGeneratePythonCmd() *cobra.Command {
	pythonCmd := &cobra.Command{
		Use:   "python",
		Short: "Generate typesafe Python client.",
		Long:  `Generate typesafe Python client compatible with the OpenFeature Python SDK.`,
		Annotations: map[string]string{
			"stability": string(generators.Alpha),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestPath := config.GetManifestPath(cmd)
			outputPath := config.GetOutputPath(cmd)

			logger.Default.GenerationStarted("Python")

			params := generators.Params[python.Params]{
				OutputPath: outputPath,
				Custom:     python.Params{},
			}
			flagset, err := flagset.Load(manifestPath)
			if err != nil {
				return err
			}

			generator := python.NewGenerator(flagset)
			logger.Default.Debug("Executing Python generator")
			err = generator.Generate(&params)
			if err != nil {
				return err
			}

			logger.Default.GenerationComplete("Python")

			return nil
		},
	}

	addStabilityInfo(pythonCmd)

	return pythonCmd
}

func init() {
	// Register generators with the manager
	generators.DefaultManager.Register(GetGenerateReactCmd)
	generators.DefaultManager.Register(GetGenerateGoCmd)
	generators.DefaultManager.Register(GetGenerateNodeJSCmd)
	generators.DefaultManager.Register(getGeneratePythonCmd)
}

func GetGenerateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate typesafe OpenFeature accessors.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "generate")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println("Available generators:")
			return generators.DefaultManager.PrintGeneratorsTable()
		},
	}

	// Add generate flags using the config package
	config.AddGenerateFlags(generateCmd)

	// Add all registered generator commands
	for _, subCmd := range generators.DefaultManager.GetCommands() {
		generateCmd.AddCommand(subCmd)
	}
	
	addStabilityInfo(generateCmd)

	return generateCmd
}
