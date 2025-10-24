package cmd

import (
	"log"
	"strings"

	"github.com/mayura-andrew/goodrive/internal/drive"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download [file-id or sharing-link]",
	Short: "Download a file from Google Drive",
	Long:  `Download a file or folder from Google Drive using its file ID or sharing link. The downloaded file will be saved in the current directory.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := drive.InitDriveClient()
		if err != nil {
			log.Fatal(err)
		}

		fileID := extractFileID(args[0])
		if err := client.CloneFile(fileID); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}

func extractFileID(input string) string {
	// Extract file ID from sharing link or return as-is
	if strings.Contains(input, "drive.google.com") {
		parts := strings.Split(input, "/")
		for i, part := range parts {
			if part == "d" && i+1 < len(parts) {
				return parts[i+1]
			}
		}
	}
	return input
}
