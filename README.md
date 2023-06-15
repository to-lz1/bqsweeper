# bqsweeper

[![Actions Status](https://github.com/to-lz1/bqsweeper/workflows/golangci-lint/badge.svg)](https://github.com/to-lz1/bqsweeper/workflows)

tool for managing and sweeping BigQuery tables.

## Install

```
go install github.com/to-lz1/bqsweeper@latest
```

## Requirements

Users must authenticate access using the command line via [gcloud CLI](https://cloud.google.com/sdk/gcloud), and have the following permissions:

* bigquery.tables.list
* bigquery.tables.update
* bigquery.tables.delete

See Also: https://cloud.google.com/bigquery/docs/access-control


## Usage

Set expiration date for tables. Table ID can be regular expression.

```
bqsweeper --project [BIGQUERY_PROJECT_ID] invalidate [DATASET_ID] [TABLE_ID_REGEXP] [EXPIRATION_DATE(yyyyMMdd)]
```
