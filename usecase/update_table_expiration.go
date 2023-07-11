package usecase

import (
	"context"
	"fmt"
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
	fmt.Printf("Expiration date for %s.%s set to %s\n", datasetID, tableID, expiration)
	return nil
}
