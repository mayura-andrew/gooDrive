package cmd

import (
	"fmt"
	"log"

	"github.com/mayura-andrew/goodrive/internal/drive"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List files in Google Drive",
	Long:  `Retrieve and display a list of files stored in your Google Drive.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := drive.InitDriveClient()
		if err != nil {
			log.Fatal(err)
		}

		nameFilter, _ := cmd.Flags().GetString("name")
		typeFilter, _ := cmd.Flags().GetString("type")

		if err := client.ViewFiles(nameFilter, typeFilter); err != nil {
			fmt.Println("Error retrieving files:", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("name", "n", "", "Filter by file name")
	listCmd.Flags().StringP("type", "t", "", "Filter by file type")
}
