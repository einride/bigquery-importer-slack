package bigqueryapi

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/bigquery"
	"github.com/einride/bigquery-importer-slack/internal/tables"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
	"google.golang.org/api/googleapi"
)

type JobClient struct {
	Config         JobConfig
	BigQueryClient *bigquery.Client
	Logger         *zap.Logger
}

func allTableRows() []tables.Row {
	return []tables.Row{
		&tables.UsersRow{},
		&tables.UserGroupsRow{},
		&tables.ChannelsRow{},
		&tables.ChannelMembersRow{},
		&tables.FilesRow{},
	}
}

// EnsureTables creates new tables.
// If a table already exists an error will be returned.
func (c *JobClient) EnsureTables(ctx context.Context) error {
	c.Logger.Info("ensuring tables")
	for _, tableRow := range allTableRows() {
		if err := c.createTable(ctx, tableRow); err != nil {
			return err
		}
	}
	return nil
}

// PutUsers adds an array of slack.User to the corresponding BigQuery table.
func (c *JobClient) PutUsers(ctx context.Context, users []slack.User) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("put users: %w", err)
		}
	}()
	if len(users) == 0 {
		return nil
	}
	valueSavers := make([]bigquery.ValueSaver, 0, len(users))
	for _, user := range users {
		row := tables.UsersRow{
			Org: c.Config.Org,
		}
		row.UnmarshalSlackUser(&user)
		valueSavers = append(valueSavers, row.ValueSaver(c.Config.ID))
	}
	c.Logger.Debug("inserting users", zap.Int("count", len(valueSavers)))
	return c.inserter(&tables.UsersRow{}).Put(ctx, valueSavers)
}

// PutUserGroups adds an array of slack.UserGroup to the corresponding BigQuery table.
func (c *JobClient) PutUserGroups(ctx context.Context, usergroups []slack.UserGroup) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("put usergroups: %w", err)
		}
	}()
	if len(usergroups) == 0 {
		return nil
	}
	valueSavers := make([]bigquery.ValueSaver, 0, len(usergroups))
	for _, usergroup := range usergroups {
		row := tables.UserGroupsRow{
			Org: c.Config.Org,
		}
		row.UnmarshalSlackUserGroup(&usergroup)
		valueSavers = append(valueSavers, row.ValueSaver(c.Config.ID))
	}
	c.Logger.Debug("inserting usergroups", zap.Int("count", len(valueSavers)))
	return c.inserter(&tables.UserGroupsRow{}).Put(ctx, valueSavers)
}

// PutChannels adds an array of slack.Channel to the corresponding BigQuery table.
func (c *JobClient) PutChannels(ctx context.Context, channels []slack.Channel) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("put channels: %w", err)
		}
	}()
	if len(channels) == 0 {
		return nil
	}
	valueSavers := make([]bigquery.ValueSaver, 0, len(channels))
	for _, channel := range channels {
		row := tables.ChannelsRow{
			Org: c.Config.Org,
		}
		row.UnmarshallSlackChannel(&channel)
		valueSavers = append(valueSavers, row.ValueSaver(c.Config.ID))
	}
	c.Logger.Debug("inserting channels", zap.Int("count", len(valueSavers)))
	return c.inserter(&tables.ChannelsRow{}).Put(ctx, valueSavers)
}

// PutChannelMembers adds an array of channel members to the corresponding BigQuery table.
func (c *JobClient) PutChannelMembers(ctx context.Context, channel *slack.Channel, members []string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("put channelmembers: %w", err)
		}
	}()
	if len(members) == 0 {
		return nil
	}
	valueSavers := make([]bigquery.ValueSaver, 0, len(members))
	for _, member := range members {
		row := tables.ChannelMembersRow{
			ChannelID:   channel.ID,
			ChannelName: channel.Name,
			Member:      member,
		}
		valueSavers = append(valueSavers, row.ValueSaver(c.Config.ID))
	}
	c.Logger.Debug("inserting channelmembers", zap.Int("count", len(valueSavers)))
	return c.inserter(&tables.ChannelMembersRow{}).Put(ctx, valueSavers)
}

// PutFiles adds an array of slack.File to the corresponding BigQuery table.
func (c *JobClient) PutFiles(ctx context.Context, files []slack.File) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("put files: %w", err)
		}
	}()
	if len(files) == 0 {
		return nil
	}
	valueSavers := make([]bigquery.ValueSaver, 0, len(files))
	for _, file := range files {
		row := &tables.FilesRow{}
		row.UnmarshalFile(&file)
		valueSavers = append(valueSavers, row.ValueSaver(c.Config.ID))
	}
	c.Logger.Debug("inserting files", zap.Int("count", len(valueSavers)))
	return c.inserter(&tables.FilesRow{}).Put(ctx, valueSavers)
}

func (c *JobClient) inserter(row tables.Row) *bigquery.Inserter {
	tableID := row.TableID(c.Config.Date)
	if c.Config.AppendIDSuffix {
		tableID = tableID + "_" + c.Config.ID.String()
	}
	return c.BigQueryClient.Dataset(c.Config.Dataset).Table(tableID).Inserter()
}

func (c *JobClient) createTable(ctx context.Context, row tables.Row) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("recreate table %s: %w", row.TableID(c.Config.Date), err)
		}
	}()
	tableID := row.TableID(c.Config.Date)
	if c.Config.AppendIDSuffix {
		tableID = tableID + "_" + c.Config.ID.String()
	}
	table := c.BigQueryClient.Dataset(c.Config.Dataset).Table(tableID)
	if _, err = table.Metadata(ctx); err == nil {
		return fmt.Errorf("table already exists: %s", table.FullyQualifiedName())
	}
	if errAPI, ok := err.(*googleapi.Error); err != nil && (!ok || errAPI.Code != http.StatusNotFound) {
		return err
	}
	c.Logger.Info("creating table", zap.Any("fullyQualifiedName", table.FullyQualifiedName()))
	return table.Create(ctx, row.TableMetadata())
}
