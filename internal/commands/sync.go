package commands

import (
	"github.com/mayura-andrew/gooDrive/internal/drive"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(pushCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of tracked files",
	Long:  `Display the synchronization status of all tracked files in the current directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := drive.NewClient()
		if err != nil {
			return err
		}
		return client.CheckStatus()
	},
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull latest changes from Drive",
	Long:  `Download the latest versions of tracked files from Google Drive.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := drive.NewClient()
		if err != nil {
			return err
		}
		return client.PullChanges()
	},
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push local changes to Drive",
	Long:  `Upload modified tracked files to Google Drive.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := drive.NewClient()
		if err != nil {
			return err
		}
		return client.PushChanges()
	},
}
