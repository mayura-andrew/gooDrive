package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/mayura-andrew/goodrive/internal/drive"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload [file]",
	Short: "Upload a file to Google Drive",
	Long:  `Upload a specified file to Google Drive and track it. You can provide the file path as an argument.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := drive.InitDriveClient()
		if err != nil {
			log.Fatal(err)
		}

		filePath := args[0]
		if err := client.AddRemote(filePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error uploading file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File uploaded successfully!")
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
