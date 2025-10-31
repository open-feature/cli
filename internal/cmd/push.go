package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/open-feature/cli/internal/api/push"
	"github.com/open-feature/cli/internal/config"
	"github.com/open-feature/cli/internal/manifest"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// GetPushCmd returns the command for pushing flags to a remote source
func GetPushCmd() *cobra.Command {
	pushCmd := &cobra.Command{
		Use:   "push",
		Short: "Push flag configurations to a remote source",
		Long: `The push command uploads local flag configurations to a remote flag management service.

This command reads your local flag manifest and intelligently pushes it to a specified
remote destination. It performs a smart push by:

1. Fetching existing flags from the remote
2. Comparing local flags with remote flags
3. Creating new flags that don't exist remotely
4. Updating existing flags that have changed

This approach ensures idempotent operations and prevents conflicts.

The pushed data follows the Manifest Management API OpenAPI specification defined at:
api/v1/push.yaml

The API uses individual flag endpoints:
- POST /api/v1/manifest/flags - Creates new flags
- PUT /api/v1/manifest/flags/{key} - Updates existing flags
- GET /api/v1/manifest - Fetches existing flags for comparison

Remote services implementing this API should accept the flag data in the format
specified by the OpenFeature flag manifest schema.

Note: The file:// scheme is not supported for push operations.
For local file operations, use standard shell commands like cp or mv.`,
		Example: `  # Push flags to a remote HTTPS endpoint (smart push: creates and updates as needed)
  openfeature push --flag-source-url https://api.example.com --auth-token secret-token

  # Push flags to an HTTP endpoint (development)
  openfeature push --flag-source-url http://localhost:8080

  # Dry run to preview what would be sent
  openfeature push --flag-source-url https://api.example.com --dry-run`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, "push")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get configuration values
			flagSourceUrl := config.GetFlagSourceUrl(cmd)
			manifestPath := config.GetManifestPath(cmd)
			authToken := config.GetAuthToken(cmd)
			dryRun := config.GetDryRun(cmd)

			// Validate destination URL is provided
			if flagSourceUrl == "" {
				return fmt.Errorf("flag source URL is required. Please provide --flag-source-url")
			}

			// Parse and validate URL
			parsedURL, err := url.Parse(flagSourceUrl)
			if err != nil {
				return fmt.Errorf("invalid source URL: %w", err)
			}

			// Load the local manifest
			flags, err := manifest.LoadFlagSet(manifestPath)
			if err != nil {
				return fmt.Errorf("error loading manifest from %s: %w", manifestPath, err)
			}

			// Validation of required fields is handled by manifest.LoadFlagSet

			// If dry run, show what would be pushed
			if dryRun {
				fmt.Printf("DRY RUN: Would push %d flags to %s\n", len(flags.Flags), flagSourceUrl)
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
				result, err := manifest.SaveToRemote(flagSourceUrl, flags, authToken)
				if err != nil {
					return fmt.Errorf("error pushing flags to remote destination: %w", err)
				}

				// Display the results
				displayPushResults(result, flagSourceUrl)
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

// displayPushResults renders the push operation results with color-coded output
func displayPushResults(result *push.PushResult, destination string) {
	totalChanges := len(result.Created) + len(result.Updated)

	if totalChanges == 0 {
		pterm.Success.Println("No changes needed - all flags are already up to date.")
		return
	}

	pterm.Success.Printf("Successfully pushed %d flag(s) to %s\n\n", totalChanges, destination)

	// Display created flags
	if len(result.Created) > 0 {
		pterm.FgGreen.Printf("◆ Created (%d):\n", len(result.Created))
		for _, flag := range result.Created {
			pterm.FgGreen.Printf("  + %s", flag.Key)
			if flag.Description != "" {
				fmt.Printf(" - %s", flag.Description)
			}
			fmt.Println()

			// Show flag details
			flagJSON, _ := json.MarshalIndent(map[string]interface{}{
				"type":         flag.Type.String(),
				"defaultValue": flag.DefaultValue,
			}, "    ", "  ")
			fmt.Printf("    %s\n", flagJSON)
		}
		fmt.Println()
	}

	// Display updated flags
	if len(result.Updated) > 0 {
		pterm.FgYellow.Printf("◆ Updated (%d):\n", len(result.Updated))
		for _, flag := range result.Updated {
			pterm.FgYellow.Printf("  ~ %s", flag.Key)
			if flag.Description != "" {
				fmt.Printf(" - %s", flag.Description)
			}
			fmt.Println()

			// Show flag details
			flagJSON, _ := json.MarshalIndent(map[string]interface{}{
				"type":         flag.Type.String(),
				"defaultValue": flag.DefaultValue,
			}, "    ", "  ")
			fmt.Printf("    %s\n", flagJSON)
		}
		fmt.Println()
	}
}
