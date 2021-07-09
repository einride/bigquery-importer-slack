package bigqueryapi

import (
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
)

type JobConfig struct {
	Dataset        string `required:"true"`
	Org            string `required:"true"`
	Date           civil.Date
	ID             uuid.UUID
	AppendIDSuffix bool
}
