package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/civil"
	"github.com/einride/bigquery-importer-slack/internal/app"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

func main() {
	var config app.Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	logger, cleanupLogger, err := app.InitLogger(&config)
	if err != nil {
		log.Panic(err)
	}
	defer cleanupLogger()
	logger.Info("initializing", zap.Any("config", &config))
	if config.Job.Date == (civil.Date{}) {
		config.Job.Date = civil.DateOf(time.Now())
		logger.Info("setting job date", zap.Stringer("date", config.Job.Date))
	}
	if config.Job.ID == (uuid.UUID{}) {
		config.Job.ID = uuid.New()
		logger.Info("setting job ID", zap.Stringer("id", config.Job.ID))
	}
	logger = logger.With(zap.Any("job", config.Job))
	app, cleanupApp, err := app.InitApp(ctx, logger, &config)
	if err != nil {
		logger.Panic("failed to initialize", zap.Error(err))
	}
	defer cleanupApp()
	if err := app.Run(ctx); err != nil {
		logger.Error("failed to run", zap.Error(err))
	}
}
