package usecase

import (
	"context"
	"errors"
	"log/slog"
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

	pagesize := 10000
	ts := client.Dataset(datasetID).Tables(ctx)
	ts.PageInfo().MaxSize = pagesize

	i := 0

	for {
		t, err := ts.Next()
		if i++; i%pagesize == 0 {
			slog.Info("Scanning tables...", "count", i)
		}
		if errors.Is(err, iterator.Done) {
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
