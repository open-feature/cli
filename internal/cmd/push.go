package cmd

import (
	"fmt"
	"net/url"

	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/manifest"
	"github.com/spf13/cobra"
)

// GetPushCmd returns the command for pushing flags to a remote source
func GetPushCmd() *cobra.Command {
	pushCmd := &cobra.Command{
		Use:   "push",
		Short: "Push flag configurations to a remote source",
		Long: `The push command uploads local flag configurations to a remote flag management service.

This command reads your local flag manifest and pushes it to a specified remote destination.
It supports HTTP and HTTPS protocols for remote endpoints.

Note: The file:// scheme is not supported for push operations.
For local file operations, use standard shell commands like cp or mv.`,
		Example: `  # Push flags to a remote HTTPS endpoint
  openfeature push --flag-destination-url https://api.example.com/flags --auth-token secret-token

  # Push flags to an HTTP endpoint (development)
  openfeature push --flag-destination-url http://localhost:8080/flags

  # Push using PUT method instead of POST
  openfeature push --flag-destination-url https://api.example.com/flags/my-app --method PUT

  # Dry run to preview what would be sent
  openfeature push --flag-destination-url https://api.example.com/flags --dry-run`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "push")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get configuration values
			flagDestinationUrl := config.GetFlagDestinationUrl(cmd)
			manifestPath := config.GetManifestPath(cmd)
			authToken := config.GetAuthToken(cmd)
			httpMethod := config.GetPushMethod(cmd)
			dryRun := config.GetDryRun(cmd)

			// Validate destination URL is provided
			if flagDestinationUrl == "" {
				return fmt.Errorf("flag destination URL is required. Please provide --flag-destination-url")
			}

			// Parse and validate URL
			parsedURL, err := url.Parse(flagDestinationUrl)
			if err != nil {
				return fmt.Errorf("invalid destination URL: %w", err)
			}

			// Load the local manifest
			flags, err := manifest.LoadFlagSet(manifestPath)
			if err != nil {
				return fmt.Errorf("error loading manifest from %s: %w", manifestPath, err)
			}

			// Validation of required fields is handled by manifest.LoadFlagSet

			// If dry run, show what would be pushed
			if dryRun {
				fmt.Printf("DRY RUN: Would push %d flags to %s\n", len(flags.Flags), flagDestinationUrl)
				fmt.Printf("Method: %s\n", httpMethod)
				if authToken != "" {
					fmt.Println("Authentication: Bearer token provided")
				}
				fmt.Println("\nFlags to push:")
				for _, flag := range flags.Flags {
					fmt.Printf("  - %s (%s): %v\n", flag.Key, flag.Type.String(), flag.DefaultValue)
				}
				return nil
			}

			// Handle URL schemes
			switch parsedURL.Scheme {
			case "file":
				return fmt.Errorf("file:// scheme is not supported for push. Use standard shell commands (cp, mv) for local file operations")
			case "http", "https":
				err = manifest.SaveToRemote(flagDestinationUrl, flags, authToken, httpMethod)
				if err != nil {
					return fmt.Errorf("error pushing flags to remote destination: %w", err)
				}
				fmt.Printf("Successfully pushed %d flags to %s\n", len(flags.Flags), flagDestinationUrl)
			default:
				return fmt.Errorf("unsupported URL scheme: %s. Supported schemes are http:// and https://", parsedURL.Scheme)
			}

			return nil
		},
	}

	// Add push-specific flags
	config.AddPushFlags(pushCmd)

	// Add common flags (like --manifest)
	config.AddRootFlags(pushCmd)

	return pushCmd
}
