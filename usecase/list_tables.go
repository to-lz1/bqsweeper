package usecase

import (
	"context"
	"regexp"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func ListTables(projectID, datasetID string, tableIDRegex *regexp.Regexp) (tableIDs []string, e error) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ts := client.Dataset(datasetID).Tables(ctx)
	for {
		t, err := ts.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		if tableIDRegex.MatchString(t.TableID) {
			tableIDs = append(tableIDs, t.TableID)
		}
	}
	return tableIDs, nil
}
