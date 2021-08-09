package tables

import (
	"strings"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/slack-go/slack"
)

// ChannelsRow follows the structure of the WebAPI. For field descriptions see the official
// documentation: https://api.slack.com/types/channel
type ChannelsRow struct {
	Org       string  `bigquery:"org"`
	Id        string  `bigquery:"id"`
	Name      string  `bigquery:"name"`
	Creator   string  `bigquery:"creator"`
	Topic     Topic   `bigquery:"topic"`
	Purpose   Purpose `bigquery:"purpose"`
	IsChannel bool    `bigquery:"is_channel"`
	IsGeneral bool    `bigquery:"is_general"`
	Locale    string  `bigquery:"locale"`
	Created   string  `bigquery:"created"`
}

var _ Row = &ChannelsRow{}

type Topic struct {
	Value   string `bigquery:"value"`
	Creator string `bigquery:"creator"`
	LastSet string `bigquery:"last_set"`
}
type Purpose struct {
	Value   string `bigquery:"value"`
	Creator string `bigquery:"creator"`
	LastSet string `bigquery:"last_set"`
}

func (c *ChannelsRow) TableID(date civil.Date) string {
	return "channels_" + strings.ReplaceAll(date.String(), "-", "")
}

func (c *ChannelsRow) ValueSaver(jobID uuid.UUID) bigquery.ValueSaver {
	return &bigquery.StructSaver{
		Schema:   c.Schema(),
		InsertID: c.InsertID(jobID),
		Struct:   c,
	}
}

func (c *ChannelsRow) Schema() bigquery.Schema {
	schema, _ := bigquery.InferSchema(c)
	return schema
}

func (c *ChannelsRow) TableMetadata() *bigquery.TableMetadata {
	return &bigquery.TableMetadata{
		Description: "channels follows the structure of the WebAPI. For field descriptions see the official " +
			"documentation: https://api.slack.com/types/channel",
		Schema: c.Schema(),
	}
}

func (c *ChannelsRow) InsertID(jobID uuid.UUID) string {
	return strings.Join([]string{
		jobID.String(),
		c.Id,
	}, "-")
}

func (c *ChannelsRow) UnmarshallSlackChannel(sc *slack.Channel) {
	if sc == nil {
		*c = ChannelsRow{}
		return
	}
	c.Id = sc.GroupConversation.ID
	c.Name = sc.GroupConversation.Name
	c.Creator = sc.GroupConversation.Creator
	c.Topic.UnmarshallTopic(&sc.Topic)
	c.Purpose.UnmarshallPurpose(&sc.Purpose)
	c.IsChannel = sc.IsChannel
	c.IsGeneral = sc.IsGeneral
	c.Locale = sc.Locale
	c.Created = sc.Created.String()
}

func (t *Topic) UnmarshallTopic(st *slack.Topic) {
	t.Value = st.Value
	t.Creator = st.Creator
	t.LastSet = st.LastSet.String()
}

func (p *Purpose) UnmarshallPurpose(sp *slack.Purpose) {
	p.Value = sp.Value
	p.Creator = sp.Creator
	p.LastSet = sp.LastSet.String()
}
