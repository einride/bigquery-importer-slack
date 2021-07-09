package tables

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
)

type Row interface {
	TableID(civil.Date) string
	TableMetadata() *bigquery.TableMetadata
	ValueSaver(uuid.UUID) bigquery.ValueSaver
}
