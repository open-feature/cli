package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/open-feature/cli/internal/logger"
	"github.com/spf13/cobra"
)

func GetVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of the OpenFeature CLI",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if Version == "dev" {
				logger.Default.Debug("Development version detected, attempting to get build info")
				details, ok := debug.ReadBuildInfo()
				if ok && details.Main.Version != "" && details.Main.Version != "(devel)" {
					Version = details.Main.Version
					for _, i := range details.Settings {
						if i.Key == "vcs.time" {
							Date = i.Value
							logger.Default.Debug(fmt.Sprintf("Found build date: %s", Date))
						}
						if i.Key == "vcs.revision" {
							Commit = i.Value
							logger.Default.Debug(fmt.Sprintf("Found commit: %s", Commit))
						}
					}
				}
			}

			versionInfo := fmt.Sprintf("OpenFeature CLI: %s (%s), built at: %s", Version, Commit, Date)
			logger.Default.Info(versionInfo)
		},
	}

	return versionCmd
}
