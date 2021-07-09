//+build wireinject

package app

import (
	"context"

	"github.com/einride/bigquery-importer-slack/internal/api/bigqueryapi"
	"github.com/einride/bigquery-importer-slack/internal/api/slackapi"
	"github.com/google/wire"
	"go.uber.org/zap"
)

func InitApp(ctx context.Context, logger *zap.Logger, config *Config) (*App, func(), error) {
	panic(
		wire.Build(
			wire.Struct(new(App), "*"),
			InitBigQueryClient,
			InitSlackClient,
			InitSecretManagerClient,
			wire.Struct(new(slackapi.SlackClient), "*"),
			wire.Struct(new(bigqueryapi.JobClient), "*"), wire.FieldsOf(&config, "Job"),
		),
	)
}
