package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/mayura-andrew/goodrive/internal/auth"
	"github.com/mayura-andrew/goodrive/internal/utils"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Google Drive",
	Long: `Authenticate with Google Drive using OAuth 2.0.

This command will open a browser window to sign in with your Google account
and grant gooDrive access to your Google Drive files.

Run this command if:
- You're using gooDrive for the first time
- Your authentication token has expired
- You want to switch to a different Google account`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize OAuth from credentials file
		if err := auth.InitOAuthFromFile(); err != nil {
			utils.LogError(fmt.Errorf("failed to initialize OAuth: %v", err))
			return
		}

		fmt.Println("\n╔════════════════════════════════════════════════════════════════╗")
		fmt.Println("║          🔐 Google Drive Authentication                        ║")
		fmt.Println("╚════════════════════════════════════════════════════════════════╝")

		// Start authentication flow
		tok := auth.GetTokenFromWeb(auth.OAuthConfig)

		// Save the token
		if err := auth.SaveToken(tok); err != nil {
			utils.LogError(fmt.Errorf("failed to save token: %v", err))
			return
		}

		fmt.Println("\n╔════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                    ✅ Authentication Successful                 ║")
		fmt.Println("╚════════════════════════════════════════════════════════════════╝")
		fmt.Println("\nYour authentication token has been saved.")
		fmt.Println("You can now use gooDrive commands like:")
		fmt.Println("  • gooDrive list    - List files in your Drive")
		fmt.Println("  • gooDrive upload  - Upload files")
		fmt.Println("  • gooDrive download - Download files")
		fmt.Println("  • gooDrive search  - Search for files")
		fmt.Println("")

		// Verify authentication by listing a few files
		fmt.Println("🔄 Verifying access...")

		client, err := auth.GetHTTPClient()
		if err != nil {
			log.Printf("Warning: Could not verify access: %v", err)
			return
		}

		// Test API call
		_, err = client.Get("https://www.googleapis.com/drive/v3/about?fields=user")
		if err != nil {
			log.Printf("Warning: Could not verify Google Drive access: %v", err)
			return
		}

		fmt.Println("✅ Successfully verified access to your Google Drive!\n")
	},
}

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh authentication token",
	Long: `Manually refresh your authentication token.

This is usually not necessary as tokens are automatically refreshed when needed,
but you can use this command to force a refresh.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize OAuth from credentials file
		if err := auth.InitOAuthFromFile(); err != nil {
			utils.LogError(fmt.Errorf("failed to initialize OAuth: %v", err))
			return
		}

		fmt.Println("\n🔄 Attempting to refresh authentication token...")

		tok, err := auth.GetToken(context.Background())
		if err != nil {
			utils.LogError(fmt.Errorf("failed to refresh token: %v. Please run 'gooDrive auth' to re-authenticate", err))
			return
		}

		fmt.Println("✅ Token refreshed successfully!")
		fmt.Printf("Token expires at: %v\n", tok.Expiry)
	},
}

func init() {
	// Add refresh as a subcommand of auth
	authCmd.AddCommand(refreshCmd)
}
