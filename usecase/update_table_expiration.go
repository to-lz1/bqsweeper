package usecase

import (
	"context"
	"log/slog"
	"time"

	"cloud.google.com/go/bigquery"
)

func UpdateTableExpiration(ctx context.Context, projectID, datasetID, tableID string, expiration time.Time) error {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()
	tableRef := client.Dataset(datasetID).Table(tableID)
	metadata, err := tableRef.Metadata(ctx)
	if err != nil {
		return err
	}
	metadataToUpdate := bigquery.TableMetadataToUpdate{
		ExpirationTime: expiration,
	}
	if _, err := tableRef.Update(ctx, metadataToUpdate, metadata.ETag); err != nil {
		return err
	}
	slog.Info("Expiration date set", "dataset", datasetID, "table", tableID, "expiration", expiration)
	return nil
}
