package app

import "github.com/einride/bigquery-importer-slack/internal/api/bigqueryapi"

type Config struct {
	Logger struct {
		ServiceName string `required:"true"`
		Level       string `required:"true"`
		Development bool   `required:"true"`
	}

	BigQueryClient struct {
		ProjectID string `required:"true"`
	}

	SlackClient struct {
		APIKeySecret string `required:"true"`
	}

	Job bigqueryapi.JobConfig
}
