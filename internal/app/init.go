package app

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/blendle/zapdriver"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func InitSlackClient(
	ctx context.Context,
	config *Config,
	secretmanager *secretmanager.Client,
	logger *zap.Logger,
) (*slack.Client, error) {
	logger.Info("init Slack client", zap.Any("cfg", config.SlackClient))
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: config.SlackClient.APIKeySecret,
	}
	APIKey, err := secretmanager.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return nil, err
	}
	return slack.New(string(APIKey.Payload.Data)), nil
}

func InitSecretManagerClient(
	ctx context.Context,
	logger *zap.Logger,
) (*secretmanager.Client, func(), error) {
	logger.Info("init Secret Manager client")
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("init Secret Manager client: %w", err)
	}
	cleanup := func() {
		logger.Info("closing Secret Manager client")
		if err := client.Close(); err != nil {
			logger.Warn("close Secret Manager client", zap.Error(err))
		}
	}
	return client, cleanup, nil
}

func InitBigQueryClient(
	ctx context.Context,
	config *Config,
	logger *zap.Logger,
) (_ *bigquery.Client, _ func(), err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("init BigQuery client: %w", err)
		}
	}()
	logger.Info("init BigQuery client", zap.Any("config", config.BigQueryClient))
	client, err := bigquery.NewClient(ctx, config.BigQueryClient.ProjectID)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		logger.Info("closing BigQuery client")
		if err := client.Close(); err != nil {
			logger.Warn("close BigQuery client", zap.Error(err))
		}
	}
	return client, cleanup, nil
}

func InitLogger(
	config *Config,
) (_ *zap.Logger, _ func(), err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("init logger: %w", err)
		}
	}()
	var zapConfig zap.Config
	var zapOptions []zap.Option
	if config.Logger.Development {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	} else {
		zapConfig = zap.NewProductionConfig()
		zapConfig.EncoderConfig = zapdriver.NewProductionEncoderConfig()
		zapOptions = append(
			zapOptions,
			zapdriver.WrapCore(
				zapdriver.ServiceName(config.Logger.ServiceName),
				zapdriver.ReportAllErrors(true),
			),
		)
	}
	if err := zapConfig.Level.UnmarshalText([]byte(config.Logger.Level)); err != nil {
		return nil, nil, err
	}
	logger, err := zapConfig.Build(zapOptions...)
	if err != nil {
		return nil, nil, err
	}
	logger = logger.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	logger.Info("logger initialized")
	cleanup := func() {
		logger.Info("closing logger, goodbye")
		_ = logger.Sync()
	}
	return logger, cleanup, nil
}
