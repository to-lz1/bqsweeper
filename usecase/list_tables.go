package usecase

import (
	"context"
	"strings"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func ListTables(projectID, datasetID, tableIDPrefix string) (tableIDs []string, e error) {
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
		if strings.HasPrefix(t.TableID, tableIDPrefix) {
			tableIDs = append(tableIDs, t.TableID)
		}
	}
	return tableIDs, nil
}
