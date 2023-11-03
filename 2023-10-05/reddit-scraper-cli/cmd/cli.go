package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/turnage/graw/reddit"
)

var companyName string

var rootCmd = &cobra.Command{
	Use:   "reddit-cli",
	Short: "CLI tool to process Reddit streams for a company name",
	Run:   execute,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&companyName, "company", "c", "", "Company name to search for (required)")
	rootCmd.MarkFlagRequired("company")
}
func execute(_ *cobra.Command, _ []string) {
	bot, err := reddit.NewScript("your user agent", 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to create Reddit script: %v", err)
	}

	if err := streamRedditWithBot(bot); err != nil {
		log.Fatal(err)
	}
}
