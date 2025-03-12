package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

func GetVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of the OpenFeature CLI",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if Version == "dev" {
				details, ok := debug.ReadBuildInfo()
				if ok && details.Main.Version != "" && details.Main.Version != "(devel)" {
					Version = details.Main.Version
					for _, i := range details.Settings {
						if i.Key == "vcs.time" {
							Date = i.Value
						}
						if i.Key == "vcs.revision" {
							Commit = i.Value
						}
					}
				}
			}

			fmt.Printf("OpenFeature CLI: %s (%s), built at: %s\n", Version, Commit, Date)
		},
	}

	return versionCmd
}