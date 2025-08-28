package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/manifest"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func promptWithValidation[T any](
	input *pterm.InteractiveTextInputPrinter,
	prompt string,
	parser func(string) (T, error),
	typeName string,
) (T, error) {
	for {
		inputString, err := input.Show(prompt)
		if err != nil {
			var zero T
			return zero, fmt.Errorf("failed to prompt for %s: %w", typeName, err)
		}

		value, err := parser(inputString)
		if err == nil {
			return value, nil
		}

		pterm.Error.Printf("Input a valid %s\n", typeName)
	}
}

func promptForDefaultValue(flag *flagset.Flag) (any, error) {
	switch flag.Type {
	case flagset.BoolType:
		options := []string{"false", "true"}
		prompt := fmt.Sprintf("Enter default value for flag '%s' (%s)", flag.Key, flag.Type)
		boolStr, err := pterm.DefaultInteractiveSelect.WithOptions(options).WithFilter(false).Show(prompt)
		if err != nil {
			return nil, fmt.Errorf("failed to prompt for bool value: %w", err)
		}
		boolValue, err := strconv.ParseBool(boolStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bool value: %w", err)
		}
		return boolValue, nil
	case flagset.IntType:
		prompt := fmt.Sprintf("Enter default value for flag '%s' (%s)", flag.Key, flag.Type)
		input := pterm.DefaultInteractiveTextInput.WithDefaultText("0")
		return promptWithValidation(input, prompt, strconv.Atoi, "integer")
	case flagset.FloatType:
		prompt := fmt.Sprintf("Enter default value for flag '%s' (%s)", flag.Key, flag.Type)
		input := pterm.DefaultInteractiveTextInput.WithDefaultText("0.0")
		parser := func(s string) (float64, error) {
			return strconv.ParseFloat(s, 64)
		}
		return promptWithValidation(input, prompt, parser, "float")
	case flagset.StringType:
		prompt := fmt.Sprintf("Enter default value for flag '%s' (%s)", flag.Key, flag.Type)
		defaultValue, err := pterm.DefaultInteractiveTextInput.WithDefaultText("").Show(prompt)
		if err != nil {
			return nil, fmt.Errorf("failed to prompt for string value: %w", err)
		}
		return defaultValue, nil
	case flagset.ObjectType:
		return nil, fmt.Errorf("object flags require a default value to be specified in the source - cannot safely prompt for object structure")
	default:
		return nil, fmt.Errorf("unsupported flag type: %s", flag.Type)
	}
}

func GetPullCmd() *cobra.Command {
	pullCmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull a flag manifest from a remote source",
		Long: `Pull a flag manifest from a remote source.

This command fetches feature flag configurations from a specified remote source and saves them locally as a manifest file.

Supported URL schemes:
- http:// - HTTP remote sources
- https:// - HTTPS remote sources  
- file:// - Local file paths

How it works:
1. Connects to the specified flag source URL
2. Downloads the flag configuration data
3. Validates and processes each flag definition
4. Prompts for missing default values (unless --no-prompt is used)
5. Writes the complete manifest to the local file system

Why pull from a remote source:
- Centralized flag management: Keep all flag definitions in a central repository or service
- Team collaboration: Share flag configurations across team members and environments
- Version control: Track changes to flag configurations over time
- Environment consistency: Ensure the same flag definitions are used across different environments
- Configuration as code: Treat flag definitions as versioned artifacts that can be reviewed and deployed`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "pull")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			flagSourceUrl := config.GetFlagSourceUrl(cmd)
			manifestPath := config.GetManifestPath(cmd)
			authToken := config.GetAuthToken(cmd)
			noPrompt := config.GetNoPrompt(cmd)

			if flagSourceUrl == "" {
				return fmt.Errorf("flagSourceUrl not set in config")
			}

			// fetch the flags from the remote source
			// Check if the URL is a local file path
			var flags *flagset.Flagset
			if strings.HasPrefix(flagSourceUrl, "file://") {
				localPath := strings.TrimPrefix(flagSourceUrl, "file://")
				loadedFlags, err := manifest.LoadFromLocal(localPath)
				if err != nil {
					return fmt.Errorf("error loading flags from local file: %w", err)
				}
				flags = loadedFlags
			} else if strings.HasPrefix(flagSourceUrl, "http://") || strings.HasPrefix(flagSourceUrl, "https://") {
				loadedFlags, err := manifest.LoadFromRemote(flagSourceUrl, authToken)
				if err != nil {
					return fmt.Errorf("error fetching flags from remote source: %w", err)
				}
				flags = loadedFlags
			} else {
				return fmt.Errorf("unsupported URL scheme: %s. Supported schemes are file://, http://, and https://", flagSourceUrl)
			}

			// Check each flag for null defaultValue
			for index, flag := range flags.Flags {
				if flag.DefaultValue == nil {
					if noPrompt {
						return fmt.Errorf("flag '%s' is missing a default value and --no-prompt was specified", flag.Key)
					}
					defaultValue, err := promptForDefaultValue(&flag)
					if err != nil {
						return fmt.Errorf("failed to get default value for flag '%s': %w", flag.Key, err)
					}
					flags.Flags[index].DefaultValue = defaultValue
				}
			}

			pterm.Success.Printfln("Successfully fetched flags from %s", flagSourceUrl)
			if err := manifest.Write(manifestPath, *flags); err != nil {
				return fmt.Errorf("error writing manifest: %w", err)
			}

			return nil
		},
	}

	config.AddPullFlags(pullCmd)

	return pullCmd
}
