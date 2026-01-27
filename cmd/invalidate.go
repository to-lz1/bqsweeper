package cmd

import (
	"bufio"
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
	Use:   "invalidate [DATASET_ID] [TABLE_ID_REGEXP] [EXPIRATION_DATE(yyyyMMdd)]",
	Short: "set an expiration date for specified BigQuery table(s)",
	Long:  "set an expiration date for specified BigQuery table(s)",
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 3 {
			return fmt.Errorf("%w: requires 3 arguments", ErrInvalidArgs)
		}
		_, err := regexp.Compile(args[1])
		if err != nil {
			return fmt.Errorf("%w: invalid table id regex: %q: %v", ErrInvalidArgs, args[1], err)
		}
		expiration, err := time.Parse(dateFormat, args[2])
		if err != nil {
			return fmt.Errorf("%w: yyyyMMdd format is expected for expiration date. your input: %q: %v", ErrInvalidArgs, args[2], err)
		}
		if expiration.Before(time.Now()) {
			return fmt.Errorf("%w: expiration date must be in the future. your input: %q", ErrInvalidArgs, args[2])
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
			return fmt.Errorf("%w: there is no table in dataset %s.%s with regex %q", ErrNoTables, ProjectID, datasetID, tableIDRegex)
		}

		cmd.Printf("Found %d table(s) in %s.%s.\n", len(tableIDs), ProjectID, datasetID)
		cmd.Printf("Are you sure you want to invalidate them all on %s? [y/n]:", expiration.Format(dateFormat))
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if err != nil || strings.TrimSpace(text) != "y" {
			cmd.Println("Aborted.")
			return nil
		}

		err = usecase.UpdateTableExpirations(ProjectID, datasetID, tableIDs, expiration)
		if err != nil {
			return err
		}
		return nil
	},
}
