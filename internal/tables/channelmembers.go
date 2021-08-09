package tables

import (
	"strings"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
)

// ChannelMembersRow is a connection between a channel and a member user.
type ChannelMembersRow struct {
	ChannelID   string `bigquery:"channel_id"`
	ChannelName string `bigquery:"channel_name"`
	Member      string `bigquery:"member"`
}

var _ Row = &ChannelMembersRow{}

func (c *ChannelMembersRow) TableID(date civil.Date) string {
	return "channel_members_" + strings.ReplaceAll(date.String(), "-", "")
}

func (c *ChannelMembersRow) ValueSaver(jobID uuid.UUID) bigquery.ValueSaver {
	return &bigquery.StructSaver{
		Schema:   c.Schema(),
		InsertID: c.InsertID(jobID),
		Struct:   c,
	}
}

func (c *ChannelMembersRow) Schema() bigquery.Schema {
	schema, _ := bigquery.InferSchema(c)
	return schema
}

func (c *ChannelMembersRow) TableMetadata() *bigquery.TableMetadata {
	return &bigquery.TableMetadata{
		Description: "channel_members is a connection between a channel and a member user.",
		Schema:      c.Schema(),
	}
}

func (c *ChannelMembersRow) InsertID(jobID uuid.UUID) string {
	return strings.Join([]string{
		jobID.String(),
		c.ChannelID,
		c.Member,
	}, "-")
}
