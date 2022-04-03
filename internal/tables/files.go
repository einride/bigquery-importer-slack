package tables

import (
	"strings"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/slack-go/slack"
)

// FilesRow follows the structure of the WebAPI. For field descriptions see the official
// documentation: https://api.slack.com/types/file
type FilesRow struct {
	ID                 string     `bigquery:"id"`
	Created            civil.Time `bigquery:"created"`
	Name               string     `bigquery:"name"`
	Title              string     `bigquery:"title"`
	Mimetype           string     `bigquery:"mimetype"`
	ImageExifRotation  int        `bigquery:"image_exif_rotation"`
	Filetype           string     `bigquery:"filetype"`
	PrettyType         string     `bigquery:"pretty_type"`
	User               string     `bigquery:"user"`
	Mode               string     `bigquery:"mode"`
	Editable           bool       `bigquery:"editable"`
	IsExternal         bool       `bigquery:"is_external"`
	ExternalType       string     `bigquery:"external_type"`
	Size               int        `bigquery:"size"`
	URL                string     `bigquery:"url"`
	URLDownload        string     `bigquery:"url_download"`
	URLPrivate         string     `bigquery:"url_private"`
	URLPrivateDownload string     `bigquery:"url_private_download"`
	OriginalH          int        `bigquery:"original_h"`
	OriginalW          int        `bigquery:"original_w"`
	Thumb64            string     `bigquery:"thumb_64"`
	Permalink          string     `bigquery:"permalink"`
	PermalinkPublic    string     `bigquery:"permalink_public"`
	EditLink           string     `bigquery:"edit_link"`
	Preview            string     `bigquery:"preview"`
	PreviewHighlight   string     `bigquery:"preview_highlight"`
	Lines              int        `bigquery:"lines"`
	LinesMore          int        `bigquery:"lines_more"`
	IsPublic           bool       `bigquery:"is_public"`
	PublicURLShared    bool       `bigquery:"public_url_shared"`
	Channels           []string   `bigquery:"channels"`
	Groups             []string   `bigquery:"groups"`
	IMs                []string   `bigquery:"ims"`
	InitialComment     Comment    `bigquery:"initial_comment"`
	CommentsCount      int        `bigquery:"comments_count"`
	NumStars           int        `bigquery:"num_stars"`
	IsStarred          bool       `bigquery:"is_starred"`
	Shares             Share      `bigquery:"shares"`
}

var _ Row = &FilesRow{}

type Comment struct {
	ID      string     `bigquery:"id"`
	Created civil.Time `bigquery:"created"`
	User    string     `bigquery:"user"`
	Comment string     `bigquery:"comment"`
}

type Share struct {
	Public  []ShareFileInfo `bigquery:"public"`
	Private []ShareFileInfo `bigquery:"private"`
}

type ShareFileInfo struct {
	ID              string   `bigquery:"id"`
	ReplyUsers      []string `bigquery:"reply_users"`
	ReplyUsersCount int      `bigquery:"reply_users_count"`
	ReplyCount      int      `bigquery:"reply_count"`
	Timestamp       string   `bigquery:"ts"`
	ThreadTimestamp string   `bigquery:"thread_ts"`
	LatestReply     string   `bigquery:"latest_reply"`
	ChannelName     string   `bigquery:"channel_name"`
	TeamID          string   `bigquery:"team_id"`
}

func (f *FilesRow) TableID(date civil.Date) string {
	return "files_" + strings.ReplaceAll(date.String(), "-", "")
}

func (f *FilesRow) ValueSaver(jobID uuid.UUID) bigquery.ValueSaver {
	return &bigquery.StructSaver{
		Schema:   f.Schema(),
		InsertID: f.InsertID(jobID),
		Struct:   f,
	}
}

func (f *FilesRow) Schema() bigquery.Schema {
	schema, _ := bigquery.InferSchema(f)
	return schema
}

func (f *FilesRow) TableMetadata() *bigquery.TableMetadata {
	return &bigquery.TableMetadata{
		Description: "files follows the structure of the WebAPI. For field descriptions see the official " +
			"documentation: https://api.slack.com/types/file",
		Schema: f.Schema(),
	}
}

func (f *FilesRow) InsertID(jobID uuid.UUID) string {
	return strings.Join([]string{
		jobID.String(),
		f.ID,
	}, "-")
}

func (f *FilesRow) UnmarshalFile(sf *slack.File) {
	if sf == nil {
		*f = FilesRow{}
		return
	}
	f.ID = sf.ID
	f.Created = civil.TimeOf(sf.Created.Time())
	f.Name = sf.Name
	f.Title = sf.Title
	f.Mimetype = sf.Mimetype
	f.ImageExifRotation = sf.ImageExifRotation
	f.Filetype = sf.Filetype
	f.PrettyType = sf.PrettyType
	f.User = sf.User
	f.Mode = sf.Mode
	f.Editable = sf.Editable
	f.IsExternal = sf.IsExternal
	f.ExternalType = sf.ExternalType
	f.Size = sf.Size
	f.URL = sf.URL
	f.URLDownload = sf.URLDownload
	f.URLPrivate = sf.URLPrivate
	f.URLPrivateDownload = sf.URLPrivateDownload
	f.OriginalH = sf.OriginalH
	f.OriginalW = sf.OriginalW
	f.Thumb64 = sf.Thumb64
	f.Permalink = sf.Permalink
	f.PermalinkPublic = sf.PermalinkPublic
	f.EditLink = sf.EditLink
	f.Preview = sf.Preview
	f.PreviewHighlight = sf.PreviewHighlight
	f.Lines = sf.Lines
	f.LinesMore = sf.LinesMore
	f.IsPublic = sf.IsPublic
	f.PublicURLShared = sf.PublicURLShared
	f.Channels = sf.Channels
	f.Groups = sf.Groups
	f.IMs = sf.IMs
	f.InitialComment.UnmarshalComment(&sf.InitialComment)
	f.CommentsCount = sf.CommentsCount
	f.NumStars = sf.NumStars
	f.IsStarred = sf.IsStarred
	f.Shares.UnmarshalShare(&sf.Shares)
}

func (c *Comment) UnmarshalComment(sc *slack.Comment) {
	c.ID = sc.ID
	c.Created = civil.TimeOf(sc.Created.Time())
	c.User = sc.User
	c.Comment = sc.Comment
}

func (s *Share) UnmarshalShare(ss *slack.Share) {
	public := make([]ShareFileInfo, 0, len(ss.Public))
	private := make([]ShareFileInfo, 0, len(ss.Private))
	for key, value := range ss.Public {
		public = UnmarshalFileInfoArray(key, value)
	}
	for key, value := range ss.Private {
		private = UnmarshalFileInfoArray(key, value)
	}
	s.Public = public
	s.Private = private
}

func UnmarshalFileInfoArray(id string, files []slack.ShareFileInfo) []ShareFileInfo {
	result := make([]ShareFileInfo, 0, len(files))
	for _, file := range files {
		file := file
		info := ShareFileInfo{ID: id}
		info.UnmarshalShareFileInfo(&file)
		result = append(result, info)
	}
	return result
}

func (s *ShareFileInfo) UnmarshalShareFileInfo(sfi *slack.ShareFileInfo) {
	s.ReplyUsers = sfi.ReplyUsers
	s.ReplyUsersCount = sfi.ReplyUsersCount
	s.ReplyCount = sfi.ReplyCount
	s.Timestamp = sfi.Ts
	s.ThreadTimestamp = sfi.ThreadTs
	s.LatestReply = sfi.LatestReply
	s.ChannelName = sfi.ChannelName
	s.TeamID = sfi.TeamID
}
