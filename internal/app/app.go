package app

import (
	"context"
	"fmt"

	"github.com/einride/bigquery-importer-slack/internal/api/bigqueryapi"
	"github.com/einride/bigquery-importer-slack/internal/api/slackapi"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type App struct {
	BigQueryJobClient *bigqueryapi.JobClient
	SlackClient       *slackapi.SlackClient
	Logger            *zap.Logger
}

// Run export all the fetched data into its corresponding table
func (a *App) Run(ctx context.Context) error {
	a.Logger.Info("running")
	defer a.Logger.Info("stopped")
	if err := a.BigQueryJobClient.EnsureTables(ctx); err != nil {
		return err
	}
	if err := a.exportUsers(ctx); err != nil {
		return err
	}
	if err := a.exportUserGroups(ctx); err != nil {
		return err
	}
	if err := a.exportChannels(ctx); err != nil {
		return err
	}
	if err := a.exportFiles(ctx); err != nil {
		return err
	}
	return nil
}

func (a *App) exportUsers(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("export users: %w", err)
		}
	}()
	a.Logger.Info("exporting users")
	return a.SlackClient.ListUsers(ctx, a.BigQueryJobClient.PutUsers)
}

func (a *App) exportUserGroups(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("export usersgroups: %w", err)
		}
	}()
	a.Logger.Info("exporting usersgroups")
	return a.SlackClient.ListUserGroups(ctx, a.BigQueryJobClient.PutUserGroups)
}

func (a *App) exportChannels(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("export channels: %w", err)
		}
	}()
	a.Logger.Info("exporting channels")
	return a.SlackClient.ListChannels(ctx, func(ctx context.Context, channels []slack.Channel) error {
		if err := a.BigQueryJobClient.PutChannels(ctx, channels); err != nil {
			return nil
		}
		for _, channel := range channels {
			if err := a.exportChannelMembers(ctx, &channel); err != nil {
				return err
			}
		}
		return nil
	})
}

func (a *App) exportChannelMembers(ctx context.Context, channel *slack.Channel) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("export channelmembers: %w", err)
		}
	}()
	a.Logger.Info("exporting channelmembers")
	return a.SlackClient.ListChannelMembers(ctx, channel, a.BigQueryJobClient.PutChannelMembers)
}

func (a *App) exportFiles(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("exporting files: %w", err)
		}
	}()
	a.Logger.Info("exporting files")
	return a.SlackClient.ListFiles(ctx, a.BigQueryJobClient.PutFiles)
}
