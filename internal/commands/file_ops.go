package commands

import (
	"github.com/mayura-andrew/gooDrive/internal/drive"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(catCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm [file-id]",
	Short: "Remove file from Drive",
	Long:  `Delete a file from Google Drive. You will be prompted for confirmation.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := drive.NewClient()
		if err != nil {
			return err
		}
		return client.RemoveFile(args[0])
	},
}

var catCmd = &cobra.Command{
	Use:   "cat [file-id]",
	Short: "View file contents without downloading",
	Long:  `Display the contents of a text-based file from Google Drive without downloading it.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := drive.NewClient()
		if err != nil {
			return err
		}
		return client.ViewFileContents(args[0])
	},
}
