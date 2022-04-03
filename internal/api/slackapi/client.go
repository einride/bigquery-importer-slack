package slackapi

import (
	"context"
	"fmt"

	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type SlackClient struct {
	Client *slack.Client
	Logger *zap.Logger
}

// ListUsers returns all the users in a workspace.
//
// Required Scopes: users:read, users:read.email (for email).
func (c *SlackClient) ListUsers(
	ctx context.Context,
	put func(context.Context, []slack.User) error,
) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("list users: %v", err)
		}
	}()
	users, err := c.Client.GetUsers()
	if err != nil {
		return err
	}
	return put(ctx, users)
}

// ListUserGroups returns all slack.UserGroup's in a workspace.
//
// Required Scopes: usergroups:read.
func (c *SlackClient) ListUserGroups(
	ctx context.Context,
	put func(context.Context, []slack.UserGroup) error,
) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("list usersgroups: %v", err)
		}
	}()
	groups, err := c.Client.GetUserGroups(slack.GetUserGroupsOptionIncludeUsers(true))
	if err != nil {
		return err
	}
	return put(ctx, groups)
}

// ListChannels returns all public and private channels in a workspace.
// Only private channels that the slack bot have been added to will be returned.
//
// Required scopes: channels:read, groups:read.
func (c *SlackClient) ListChannels(
	ctx context.Context,
	put func(context.Context, []slack.Channel) error,
) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("list channels: %v", err)
		}
	}()
	var cursor string
	for {
		channels, nextCursor, err := c.Client.GetConversations(&slack.GetConversationsParameters{
			Cursor:          cursor,
			ExcludeArchived: true,
			Types:           []string{"public_channel", "private_channel"},
		})
		if err != nil {
			return err
		}
		if err := put(ctx, channels); err != nil {
			return fmt.Errorf("list channels: %v", err)
		}
		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}
	return nil
}

// ListChannelMembers returns the ids of the members in a channel.
func (c *SlackClient) ListChannelMembers(
	ctx context.Context,
	channel *slack.Channel,
	put func(context.Context, *slack.Channel, []string) error,
) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("list channel members: %w", err)
		}
	}()
	var cursor string
	for {
		users, nextCursor, err := c.Client.GetUsersInConversation(&slack.GetUsersInConversationParameters{
			ChannelID: channel.ID,
			Cursor:    cursor,
		})
		if err != nil {
			return fmt.Errorf("list channel members: %w", err)
		}
		if err = put(ctx, channel, users); err != nil {
			return fmt.Errorf("list channel members: %w", err)
		}
		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}
	return nil
}

// ListFiles returns an array of slack.File.
// The bot only has access to files in channels that it has been added to.
//
// Required Scopes: files:read.
func (c *SlackClient) ListFiles(
	ctx context.Context,
	put func(context.Context, []slack.File) error,
) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("list files: %w", err)
		}
	}()
	params := slack.ListFilesParameters{Cursor: ""}
	for {
		files, newParams, err := c.Client.ListFiles(params)
		if err != nil {
			return err
		}
		if err = put(ctx, files); err != nil {
			return err
		}
		if newParams.Cursor == "" {
			break
		}
		params = *newParams
	}
	return nil
}
