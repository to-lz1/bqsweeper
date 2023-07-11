package usecase

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

func UpdateTableExpirations(projectID, datasetID string, tableIDs []string, expiration time.Time) error {
	ctx := context.Background()

	g, ctx := errgroup.WithContext(ctx)
	tasks := make(chan string)

	// see also: https://cloud.google.com/bigquery/quotas#api_request_quotas
	const workersNum = 16
	for i := 0; i < workersNum; i++ {
		g.Go(func() error {
			for tableID := range tasks {
				UpdateTableExpiration(ctx, projectID, datasetID, tableID, expiration)
			}
			return nil
		})
	}

	go func() {
		defer close(tasks)
		for _, tableID := range tableIDs {
			select {
			case tasks <- tableID:
			case <-ctx.Done():
				return
			}
		}
	}()

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}
