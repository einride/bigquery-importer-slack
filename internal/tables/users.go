package tables

import (
	"strings"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/slack-go/slack"
)

// UsersRow follows the structure of the WebAPI. For field descriptions see the official
// documentation: https://api.slack.com/types/user
type UsersRow struct {
	Org               string      `bigquery:"org"`
	ID                string      `bigquery:"id"`
	TeamID            string      `bigquery:"team_id"`
	Deleted           bool        `bigquery:"deleted"`
	RealName          string      `bigquery:"real_name"`
	TZ                string      `bigquery:"tz"`
	TZLabel           string      `bigquery:"tz_label"`
	TZOffset          int         `bigquery:"tz_offset"`
	Profile           UserProfile `bigquery:"profile"`
	IsBot             bool        `bigquery:"is_bot"`
	IsAdmin           bool        `bigquery:"is_admin"`
	IsOwner           bool        `bigquery:"is_owner"`
	IsPrimaryOwner    bool        `bigquery:"is_primary_owner"`
	IsRestricted      bool        `bigquery:"is_restricted"`
	IsUltraRestricted bool        `bigquery:"is_ultra_restricted"`
	IsStranger        bool        `bigquery:"is_stranger"`
	IsAppUser         bool        `bigquery:"is_app_user"`
	IsInvitedUser     bool        `bigquery:"is_invited_user"`
	Has2FA            bool        `bigquery:"has_2fa"`
	HasFiles          bool        `bigquery:"has_files"`
	Presence          string      `bigquery:"presence"`
	Locale            string      `bigquery:"locale"`
}

var _ Row = &UsersRow{}

type UserProfile struct {
	FirstName             string `bigquery:"first_name"`
	LastName              string `bigquery:"last_name"`
	RealName              string `bigquery:"real_name"`
	RealNameNormalized    string `bigquery:"real_name_normalized"`
	DisplayName           string `bigquery:"display_name"`
	DisplayNameNormalized string `bigquery:"display_name_normalized"`
	Email                 string `bigquery:"email"`
	Skype                 string `bigquery:"skype"`
	Phone                 string `bigquery:"phone"`
	Title                 string `bigquery:"title"`
	BotID                 string `bigquery:"bot_id"`
	APIAppID              string `bigquery:"api_app_id"`
	StatusText            string `bigquery:"status_text"`
	StatusEmoji           string `bigquery:"status_emoji"`
	StatusExpiration      int    `bigquery:"status_expiration"`
	Team                  string `bigquery:"team"`
}

func (u *UsersRow) TableID(date civil.Date) string {
	return "users_" + strings.ReplaceAll(date.String(), "-", "")
}

func (u *UsersRow) ValueSaver(jobID uuid.UUID) bigquery.ValueSaver {
	return &bigquery.StructSaver{
		Schema:   u.Schema(),
		InsertID: u.InsertID(jobID),
		Struct:   u,
	}
}

func (u *UsersRow) Schema() bigquery.Schema {
	schema, _ := bigquery.InferSchema(u)
	return schema
}

func (u *UsersRow) TableMetadata() *bigquery.TableMetadata {
	return &bigquery.TableMetadata{
		Description: "users follows the structure of the WebAPI. For field descriptions see the official " +
			"documentation: https://api.slack.com/types/user",
		Schema: u.Schema(),
	}
}

func (u *UsersRow) InsertID(jobID uuid.UUID) string {
	return strings.Join([]string{
		jobID.String(),
		u.ID,
	}, "-")
}

func (u *UsersRow) UnmarshalSlackUser(su *slack.User) {
	if su == nil {
		*u = UsersRow{}
		return
	}
	u.ID = su.ID
	u.TeamID = su.TeamID
	u.Deleted = su.Deleted
	u.TZ = su.TZ
	u.TZLabel = su.TZLabel
	u.TZOffset = su.TZOffset
	u.Profile.UnmarshalSlackUserProfile(&su.Profile)
	u.IsBot = su.IsBot
	u.IsAdmin = su.IsAdmin
	u.IsOwner = su.IsOwner
	u.IsPrimaryOwner = su.IsPrimaryOwner
	u.IsRestricted = su.IsRestricted
	u.IsUltraRestricted = su.IsUltraRestricted
	u.IsStranger = su.IsStranger
	u.IsAppUser = su.IsAppUser
	u.IsInvitedUser = su.IsInvitedUser
	u.Has2FA = su.Has2FA
	u.HasFiles = su.HasFiles
	u.Presence = su.Presence
	u.Locale = su.Locale
}

func (u *UserProfile) UnmarshalSlackUserProfile(up *slack.UserProfile) {
	u.FirstName = up.FirstName
	u.LastName = up.LastName
	u.RealName = up.RealName
	u.RealNameNormalized = up.RealNameNormalized
	u.DisplayName = up.DisplayName
	u.DisplayNameNormalized = up.DisplayNameNormalized
	u.Email = up.Email
	u.Skype = up.Skype
	u.Phone = up.Phone
	u.Title = up.Title
	u.BotID = up.BotID
	u.APIAppID = up.ApiAppID
	u.StatusText = up.StatusText
	u.StatusEmoji = up.StatusEmoji
	u.StatusExpiration = up.StatusExpiration
	u.Team = up.Team
}
