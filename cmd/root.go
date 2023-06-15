package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var ProjectID string
var rootCmd = &cobra.Command{
	Use:   "bqsweeper",
	Short: "bqsweeper is a tool for managing and sweeping BigQuery tables",
	Long:  "bqsweeper is a tool for managing and sweeping BigQuery tables",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&ProjectID, "project", "p", "", "GCP project ID")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
