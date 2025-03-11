package cmd

import (
	"fmt"

	"github.com/open-feature/cli/internal/filesystem"
	"github.com/spf13/cobra"
)

func GetPullCmd() *cobra.Command {
	pullCmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull a flag manifest from a remote source",
		Long:  "Pull a flag manifest from a remote source.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "pull")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			flagSourceUrl, err := cmd.Flags().GetString("flagSourceUrl")
			if err != nil {
				flagSourceUrl, err = filesystem.GetFromYaml("flagSourceUrl")
				if err != nil {
					return fmt.Errorf("error getting flagSourceUrl from config: %w", err)
				}
			}

			fmt.Println(flagSourceUrl)

			// // fetch the flags from the remote source
			// resp, err := http.Get(flagSourceUrl)
			// if err != nil {
			// 	return fmt.Errorf("error fetching flags: %w", err)
			// }
			// defer resp.Body.Close()

			// flags, err := io.ReadAll(resp.Body)
			// if err != nil {
			// 	return fmt.Errorf("error reading response body: %w", err)
			// }

			return nil
		},
	}

	pullCmd.Flags().String("flagSourceUrl", "", "The URL of the flag source")

	return pullCmd
}
