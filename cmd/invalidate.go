package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/to-lz1/bqsweeper/usecase"
)

func init() {
	rootCmd.AddCommand(invalidateCmd)
}

var dateFormat = "20060102"

var invalidateCmd = &cobra.Command{
	Use:   "invalidate [datasetID] [tableIDRegex] [expiration(yyyyMMdd)]",
	Short: "set an expiration date for specified BigQuery table(s)",
	Long:  "set an expiration date for specified BigQuery table(s)",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 {
			return errors.New("requires 3 arguments")
		}
		_, err := regexp.Compile(args[1])
		if err != nil {
			return fmt.Errorf("Invalid tableIDRegex: %s", args[1])
		}
		expiration, err := time.Parse(dateFormat, args[2])
		if err != nil {
			return fmt.Errorf("Invalid expiration date(yyyyMMdd format is expected): %s", args[2])
		}
		if expiration.Before(time.Now()) {
			return errors.New("Expiration date must be in the future")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		datasetID := args[0]
		tableIDRegex := regexp.MustCompile(args[1])
		expiration, _ := time.Parse(dateFormat, args[2])

		tableIDs, err := usecase.ListTables(ProjectID, datasetID, tableIDRegex)
		if err != nil {
			return err
		}
		if len(tableIDs) == 0 {
			return fmt.Errorf("No tables found in %s.%s with regex %s", ProjectID, datasetID, tableIDRegex)
		}

		fmt.Printf("Found %d table(s) in %s.%s.\n", len(tableIDs), ProjectID, datasetID)
		fmt.Printf("Are you sure you want to invalidate them all on %s? [y/n]:", expiration.Format(dateFormat))
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if err != nil || strings.TrimSpace(text) != "y" {
			fmt.Println("Aborted.")
			return nil
		}

		for _, tableID := range tableIDs {
			err = usecase.UpdateTableExpiration(ProjectID, datasetID, tableID, expiration)
			if err != nil {
				return err
			}
		}
		return nil
	},
}
