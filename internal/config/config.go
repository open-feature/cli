package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	ManifestFlag = "manifest"
	NoInputFlag  = "no-input"
	OverrideFlag = "override"

	// vCfg is the global configuration object
	vCfg = viper.New()
)

func Load(path string, cmd *cobra.Command) error {
	var err error
	if len(path) != 0 {
		err = useUserSuppliedConfigFile(path)
	} else {
		err = useDefaultConfigFile()
	}
	if err != nil {
		return err
	}

	// Bind global flags to viper (so flag values override config)
	cmd.InheritedFlags().VisitAll(func(f *pflag.Flag) {
		if err := vCfg.BindPFlag(f.Name, f); err != nil {
			fmt.Fprintf(os.Stderr, "Error binding flag '%s': %v\n", f.Name, err)
		}
	})

	return err
}

// func Get(key string) any {
// 	return vCfg.Get(key)
// }

func GetString(key string) string {
	return vCfg.GetString(key)
}

func GetBool(key string) bool {
	return vCfg.GetBool(key)
}

func Update(key string, value any) {
	vCfg.Set(key, value)
}

func Save() error {
	return vCfg.WriteConfig()
}

func useDefaultConfigFile() error {
	vCfg.SetConfigName(".openfeature")
	vCfg.SetConfigType("yaml")
	vCfg.AddConfigPath(".")
	err := vCfg.ReadInConfig()
	if err == nil {
		return nil
	}
	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		return err
	}
	// config file isn't required
	return nil
}

func useUserSuppliedConfigFile(configFilePath string) error {
	vCfg.SetConfigFile(os.ExpandEnv(configFilePath))
	return vCfg.ReadInConfig()
}
