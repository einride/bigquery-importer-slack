package tables

import (
	"strings"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/slack-go/slack"
)

type ChannelsRow struct {
	GroupConversation GroupConversation `bigquery:"group_conversation"`
	IsChannel         bool              `bigquery:"is_channel"`
	IsGeneral         bool              `bigquery:"is_general"`
	IsMember          bool              `bigquery:"is_member"`
	Locale            string            `bigquery:"locale"`
}

var _ Row = &ChannelsRow{}

type GroupConversation struct {
	Conversation Conversation `bigquery:"conversation"`
	Name         string       `bigquery:"name"`
	Creator      string       `bigquery:"creator"`
	IsArchived   bool         `bigquery:"is_archived"`
	Members      []string     `bigquery:"members"`
	Topic        Topic        `bigquery:"topic"`
	Purpose      Purpose      `bigquery:"purpose"`
}

type Conversation struct {
	ID                 string         `bigquery:"id"`
	Created            civil.DateTime `bigquery:"created"`
	IsOpen             bool           `bigquery:"is_open"`
	LastRead           string         `bigquery:"last_read"`
	UnreadCount        int            `bigquery:"unread_count"`
	UnreadCountDisplay int            `bigquery:"unread_count_display"`
	IsGroup            bool           `bigquery:"is_group"`
	IsShared           bool           `bigquery:"is_shared"`
	IsIM               bool           `bigquery:"is_im"`
	IsExtShared        bool           `bigquery:"is_ext_shared"`
	IsOrgShared        bool           `bigquery:"is_org_shared"`
	IsPendingExtShared bool           `bigquery:"is_pending_ext_shared"`
	IsPrivate          bool           `bigquery:"is_private"`
	IsMpIM             bool           `bigquery:"is_mpim"`
	Unlinked           int            `bigquery:"unlinked"`
	NameNormalized     string         `bigquery:"name_normalized"`
	NumMembers         int            `bigquery:"num_members"`
	Priority           float64        `bigquery:"priority"`
	User               string         `bigquery:"user"`
}

type Topic struct {
	Value   string         `bigquery:"value"`
	Creator string         `bigquery:"creator"`
	LastSet civil.DateTime `bigquery:"last_set"`
}
type Purpose struct {
	Value   string         `bigquery:"value"`
	Creator string         `bigquery:"creator"`
	LastSet civil.DateTime `bigquery:"last_set"`
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
		Schema: c.Schema(),
	}
}

func (c *ChannelsRow) InsertID(jobID uuid.UUID) string {
	return strings.Join([]string{
		jobID.String(),
		c.GroupConversation.Conversation.ID,
	}, "-")
}

func (c *ChannelsRow) UnmarshalSlackChannel(sc *slack.Channel) {
	if sc == nil {
		*c = ChannelsRow{}
		return
	}
	c.GroupConversation.UnmarshalGroupsConversation(&sc.GroupConversation)
	c.IsChannel = sc.IsChannel
	c.IsGeneral = sc.IsGeneral
	c.IsMember = sc.IsMember
	c.Locale = sc.Locale
}

func (g *GroupConversation) UnmarshalGroupsConversation(sg *slack.GroupConversation) {
	g.Conversation.UnmarshalConversation(&sg.Conversation)
	g.Name = sg.Name
	g.Creator = sg.Creator
	g.IsArchived = sg.IsArchived
	g.Members = sg.Members
	g.Topic.UnmarshalTopic(&sg.Topic)
	g.Purpose.UnmarshallPurpose(&sg.Purpose)
}

func (c *Conversation) UnmarshalConversation(sc *slack.Conversation) {
	c.ID = sc.ID
	c.Created = civil.DateTimeOf(sc.Created.Time())
	c.IsOpen = sc.IsOpen
	c.LastRead = sc.LastRead
	c.UnreadCount = sc.UnreadCount
	c.UnreadCountDisplay = sc.UnreadCountDisplay
	c.IsGroup = sc.IsGroup
	c.IsShared = sc.IsShared
	c.IsIM = sc.IsIM
	c.IsExtShared = sc.IsExtShared
	c.IsOrgShared = sc.IsOrgShared
	c.IsPendingExtShared = sc.IsPendingExtShared
	c.IsPrivate = sc.IsPrivate
	c.IsMpIM = sc.IsMpIM
	c.Unlinked = sc.Unlinked
	c.NameNormalized = sc.NameNormalized
	c.NumMembers = sc.NumMembers
	c.Priority = sc.Priority
	c.User = sc.User
}

func (t *Topic) UnmarshalTopic(st *slack.Topic) {
	t.Value = st.Value
	t.Creator = st.Creator
	t.LastSet = civil.DateTimeOf(st.LastSet.Time())
}

func (p *Purpose) UnmarshallPurpose(sp *slack.Purpose) {
	p.Value = sp.Value
	p.Creator = sp.Creator
	p.LastSet = civil.DateTimeOf(sp.LastSet.Time())
}
