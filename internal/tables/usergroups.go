package tables

import (
	"strings"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/slack-go/slack"
)

// UserGroupsRow follows the structure of the WebAPI. For field descriptions see the official
// documentation: https://api.slack.com/types/usergroup
type UserGroupsRow struct {
	Org         string         `bigquery:"org"`
	ID          string         `bigquery:"id"`
	TeamID      string         `bigquery:"team_id"`
	IsUserGroup bool           `bigquery:"is_usergroup"`
	Name        string         `bigquery:"name"`
	Description string         `bigquery:"description"`
	Handle      string         `bigquery:"handle"`
	IsExternal  bool           `bigquery:"is_external"`
	DateUpdate  string         `bigquery:"date_update"`
	DateDelete  string         `bigquery:"date_delete"`
	AutoType    string         `bigquery:"auto_type"`
	CreatedBy   string         `bigquery:"created_by"`
	UpdatedBy   string         `bigquery:"updated_by"`
	DeletedBy   string         `bigquery:"deleted_by"`
	Prefs       UserGroupPrefs `bigquery:"prefs"`
	UserCount   int            `bigquery:"user_count"`
	Users       []string       `bigquery:"users"`
}

var _ = &UserGroupsRow{}

type UserGroupPrefs struct {
	Channels []string `bigquery:"channels"`
	Groups   []string `bigquery:"groups"`
}

func (u *UserGroupsRow) TableID(date civil.Date) string {
	return "usergroups_" + strings.ReplaceAll(date.String(), "-", "")
}

func (u *UserGroupsRow) ValueSaver(jobID uuid.UUID) bigquery.ValueSaver {
	return &bigquery.StructSaver{
		Schema:   u.Schema(),
		InsertID: u.InsertID(jobID),
		Struct:   u,
	}
}

func (u *UserGroupsRow) Schema() bigquery.Schema {
	schema, _ := bigquery.InferSchema(u)
	return schema
}

func (u *UserGroupsRow) TableMetadata() *bigquery.TableMetadata {
	return &bigquery.TableMetadata{
		Description: "usergroups follows the structure of the WebAPI. For field descriptions see the official " +
			"documentation: https://api.slack.com/types/usergroup",
		Schema: u.Schema(),
	}
}

func (u *UserGroupsRow) InsertID(jobID uuid.UUID) string {
	return strings.Join([]string{
		jobID.String(),
		u.ID,
	}, "-")
}

func (u *UserGroupsRow) UnmarshalSlackUserGroup(su *slack.UserGroup) {
	if su == nil {
		*u = UserGroupsRow{}
		return
	}
	u.ID = su.ID
	u.TeamID = su.TeamID
	u.IsUserGroup = su.IsUserGroup
	u.Name = su.Name
	u.Description = su.Description
	u.Handle = su.Handle
	u.IsExternal = su.IsExternal
	u.DateUpdate = su.DateUpdate.String()
	u.DateDelete = su.DateDelete.String()
	u.DeletedBy = su.DeletedBy
	u.Prefs.UnmarshalUserGroupPrefs(&su.Prefs)
	u.UserCount = su.UserCount
	u.Users = su.Users
}

func (u *UserGroupPrefs) UnmarshalUserGroupPrefs(up *slack.UserGroupPrefs) {
	u.Groups = up.Groups
	u.Channels = up.Channels
}
