package cmd

import (
	"fmt"
	"log"

	"github.com/mayura-andrew/goodrive/internal/drive"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for files in Google Drive",
	Long:  `Search for files in Google Drive using the specified query. The query can include file names, types, and other metadata.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := drive.InitDriveClient()
		if err != nil {
			log.Fatal(err)
		}

		query := args[0]
		results, err := client.SearchFiles(query)
		if err != nil {
			fmt.Printf("Error searching files: %v\n", err)
			return
		}

		if len(results) == 0 {
			fmt.Println("No files found.")
			return
		}

		fmt.Println("Search Results:")
		for _, file := range results {
			fmt.Printf("- %s (ID: %s)\n", file.Name, file.ID)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
