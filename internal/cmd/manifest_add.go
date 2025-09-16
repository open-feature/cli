package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/flagset"
	"github.com/open-feature/cli/internal/logger"
	"github.com/open-feature/cli/internal/manifest"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func GetManifestAddCmd() *cobra.Command {
	manifestAddCmd := &cobra.Command{
		Use:   "add [flag-name]",
		Short: "Add a new flag to the manifest",
		Long: `Add a new flag to the manifest file with the specified configuration.

Examples:
  # Add a boolean flag (default type)
  openfeature manifest add new-feature --default-value false

  # Add a string flag with description
  openfeature manifest add welcome-message --type string --default-value "Hello!" --description "Welcome message for users"

  # Add an integer flag
  openfeature manifest add max-retries --type integer --default-value 3

  # Add a float flag
  openfeature manifest add discount-rate --type float --default-value 0.15

  # Add an object flag
  openfeature manifest add config --type object --default-value '{"key":"value"}'`,
		Args: cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "manifest.add")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			flagName := args[0]
			manifestPath := config.GetManifestPath(cmd)

			// Get flag configuration from command flags
			flagType, _ := cmd.Flags().GetString("type")
			defaultValueStr, _ := cmd.Flags().GetString("default-value")
			description, _ := cmd.Flags().GetString("description")

			// Validate that default-value is provided
			if !cmd.Flags().Changed("default-value") {
				return errors.New("--default-value is required")
			}

			// Parse flag type
			parsedType, err := parseFlagTypeString(flagType)
			if err != nil {
				return fmt.Errorf("invalid flag type: %w", err)
			}

			// Parse and validate default value
			defaultValue, err := parseDefaultValue(defaultValueStr, parsedType)
			if err != nil {
				return fmt.Errorf("invalid default value for type %s: %w", flagType, err)
			}

			// Load existing manifest
			fs, err := manifest.LoadFlagSet(manifestPath)
			if err != nil {
				// If manifest doesn't exist, create a new one
				if strings.Contains(err.Error(), "error reading contents") {
					fs = &flagset.Flagset{
						Flags: []flagset.Flag{},
					}
				} else {
					return fmt.Errorf("failed to load manifest: %w", err)
				}
			}

			// Check if flag already exists
			for _, flag := range fs.Flags {
				if flag.Key == flagName {
					return fmt.Errorf("flag '%s' already exists in the manifest", flagName)
				}
			}

			// Add new flag
			newFlag := flagset.Flag{
				Key:          flagName,
				Type:         parsedType,
				Description:  description,
				DefaultValue: defaultValue,
			}
			fs.Flags = append(fs.Flags, newFlag)

			// Write updated manifest
			if err := manifest.Write(manifestPath, *fs); err != nil {
				return fmt.Errorf("failed to write manifest: %w", err)
			}

			// Success message
			pterm.Success.Printfln("Flag '%s' added successfully to %s", flagName, manifestPath)
			logger.Default.Debug(fmt.Sprintf("Added flag: name=%s, type=%s, defaultValue=%v, description=%s",
				flagName, flagType, defaultValue, description))

			return nil
		},
	}

	// Add command-specific flags
	config.AddManifestAddFlags(manifestAddCmd)
	addStabilityInfo(manifestAddCmd)

	return manifestAddCmd
}

// parseFlagTypeString converts a string flag type to FlagType enum
func parseFlagTypeString(typeStr string) (flagset.FlagType, error) {
	switch strings.ToLower(typeStr) {
	case "boolean", "bool":
		return flagset.BoolType, nil
	case "string":
		return flagset.StringType, nil
	case "integer", "int":
		return flagset.IntType, nil
	case "float", "number":
		return flagset.FloatType, nil
	case "object", "json":
		return flagset.ObjectType, nil
	default:
		return flagset.UnknownFlagType, fmt.Errorf("unknown flag type: %s", typeStr)
	}
}

// parseDefaultValue parses and validates the default value based on flag type
func parseDefaultValue(value string, flagType flagset.FlagType) (interface{}, error) {
	switch flagType {
	case flagset.BoolType:
		switch strings.ToLower(value) {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return nil, fmt.Errorf("boolean value must be 'true' or 'false', got '%s'", value)
		}
	case flagset.StringType:
		return value, nil
	case flagset.IntType:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid integer value: %s", value)
		}
		return intVal, nil
	case flagset.FloatType:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float value: %s", value)
		}
		return floatVal, nil
	case flagset.ObjectType:
		var jsonObj interface{}
		if err := json.Unmarshal([]byte(value), &jsonObj); err != nil {
			return nil, fmt.Errorf("invalid JSON object: %s", err.Error())
		}
		return jsonObj, nil
	default:
		return nil, fmt.Errorf("unsupported flag type: %v", flagType)
	}
}