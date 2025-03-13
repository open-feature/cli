package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// initializeConfig reads in config file and ENV variables if set.
// It applies configuration values to command flags based on hierarchical priority.
func initializeConfig(cmd *cobra.Command, bindPrefix string) error {
	v := viper.New()

	// Set the config file name and path
	v.SetConfigName(".openfeature")
	v.AddConfigPath(".")

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// Track which flags were set directly via command line
	cmdLineFlags := make(map[string]bool)
	cmd.Flags().Visit(func(f *pflag.Flag) {
		cmdLineFlags[f.Name] = true
	})

	// Apply the configuration values
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Skip if flag was set on command line
		if cmdLineFlags[f.Name] {
			return
		}

		// Build configuration paths from most specific to least specific
		configPaths := []string{}
		
		// Check the most specific path (e.g., generate.go.package-name)
		if bindPrefix != "" {
			configPaths = append(configPaths, bindPrefix + "." + f.Name)
			
			// Check parent paths (e.g., generate.package-name)
			parts := strings.Split(bindPrefix, ".")
			for i := len(parts) - 1; i > 0; i-- {
				parentPath := strings.Join(parts[:i], ".") + "." + f.Name
				configPaths = append(configPaths, parentPath)
			}
		}
		
		// Check the base path (e.g., package-name)
		configPaths = append(configPaths, f.Name)
		
		// Try each path in order until we find a match
		for _, path := range configPaths {
			if v.IsSet(path) {
				val := v.Get(path)
				_ = f.Value.Set(fmt.Sprintf("%v", val))
				break
			}
		}
	})

	return nil
}
