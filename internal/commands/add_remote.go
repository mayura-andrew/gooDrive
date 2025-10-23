package commands

import (
	"github.com/mayura-andrew/gooDrive/internal/drive"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addRemoteCmd)
}

var addRemoteCmd = &cobra.Command{
	Use:   "add-remote [local-file-path]",
	Short: "Upload local file to Google Drive and track it",
	Long: `Upload a local file to Google Drive and begin tracking it for sync operations.
	
The file will be uploaded to your Drive root folder and tracked for future
pull/push operations.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := drive.NewClient()
		if err != nil {
			return err
		}

		return client.AddRemote(args[0])
	},
}
