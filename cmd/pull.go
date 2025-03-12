package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/filesystem"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/manifest"
	"github.com/open-feature/cli/internal/requests"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func promptForDefaultValue(flag *flagset.Flag) (any) {
	var prompt string
	switch flag.Type {
	case flagset.BoolType:
		var options []string = []string{"false", "true"}
		prompt = fmt.Sprintf("Enter default value for flag '%s' (%s)", flag.Key, flag.Type)
		boolStr, _ := pterm.DefaultInteractiveSelect.WithOptions(options).WithFilter(false).Show(prompt)
		boolValue, _ := strconv.ParseBool(boolStr)
		return boolValue
	case flagset.IntType:
		var err error = errors.New("Input a valid integer")
		prompt = fmt.Sprintf("Enter default value for flag '%s' (%s)", flag.Key, flag.Type)
		var defaultValue int
		for err != nil {
			defaultValueString, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("0").Show(prompt)
			defaultValue, err = strconv.Atoi(defaultValueString)
		}
		return defaultValue
	case flagset.FloatType:
		var err error = errors.New("Input a valid float")
		prompt = fmt.Sprintf("Enter default value for flag '%s' (%s)", flag.Key, flag.Type)
		var defaultValue float64
		for err != nil {
			defaultValueString, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("0.0").Show(prompt)
			defaultValue, err = strconv.ParseFloat(defaultValueString, 64)
			if err != nil {
				pterm.Error.Println("Input a valid float")
			}
		}
		return defaultValue
	case flagset.StringType:
		prompt = fmt.Sprintf("Enter default value for flag '%s' (%s)", flag.Key, flag.Type)
		defaultValue, _ := pterm.DefaultInteractiveTextInput.WithDefaultText("").Show(prompt)
		return defaultValue
	// TODO: Add proper support for object type
	case flagset.ObjectType:
		return map[string]any{}
	default:
		return nil
	}
}

func GetPullCmd() *cobra.Command {
	pullCmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull a flag manifest from a remote source",
		Long:  "Pull a flag manifest from a remote source.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "pull")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			flagSourceUrl := config.GetFlagSourceUrl(cmd)
			manifestPath := config.GetManifestPath(cmd)
			authToken := config.GetAuthToken(cmd)

			var err error
			if flagSourceUrl == "" {
				flagSourceUrl, err = filesystem.GetFromYaml("flagSourceUrl")
				if err != nil {
					return fmt.Errorf("error getting flagSourceUrl from config: %w", err)
				}
			}

			// fetch the flags from the remote source
			flags, err := requests.FetchFlags(flagSourceUrl, authToken)
			if err != nil {
				return fmt.Errorf("error fetching flags: %w", err)
			}

			// Check each flag for null defaultValue
			for index, flag := range flags.Flags {
				if flag.DefaultValue == nil {
					defaultValue := promptForDefaultValue(&flag)
					flags.Flags[index].DefaultValue = defaultValue
				}
			}

			pterm.Success.Printf("Successfully fetched flags from %s", flagSourceUrl)
			err = manifest.Write(manifestPath, flags)
			if err != nil {
				return fmt.Errorf("error writing manifest: %w", err)
			}

			return nil
		},
	}

	config.AddPullFlags(pullCmd)

	return pullCmd
}
