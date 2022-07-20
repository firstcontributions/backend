package storiesstore

import (
	"time"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/firstcontributions/backend/pkg/cursor"
)

type StorySortBy uint8

const (
	StorySortByDefault StorySortBy = iota
	StorySortByTimeCreated
)

type Story struct {
	UserID          string    `bson:"user_id"`
	AbstractContent string    `bson:"abstract_content,omitempty"`
	ContentJson     string    `bson:"content_json,omitempty"`
	CreatedBy       string    `bson:"created_by,omitempty"`
	Id              string    `bson:"_id"`
	Thumbnail       string    `bson:"thumbnail,omitempty"`
	TimeCreated     time.Time `bson:"time_created,omitempty"`
	TimeUpdated     time.Time `bson:"time_updated,omitempty"`
	Title           string    `bson:"title,omitempty"`
	UrlSuffix       string    `bson:"url_suffix,omitempty"`
}

func NewStory() *Story {
	return &Story{}
}
func (story *Story) Get(field string) interface{} {
	switch field {
	case "user_id":
		return story.UserID
	case "abstract_content":
		return story.AbstractContent
	case "content_json":
		return story.ContentJson
	case "created_by":
		return story.CreatedBy
	case "_id":
		return story.Id
	case "thumbnail":
		return story.Thumbnail
	case "time_created":
		return story.TimeCreated
	case "time_updated":
		return story.TimeUpdated
	case "title":
		return story.Title
	case "url_suffix":
		return story.UrlSuffix
	default:
		return nil
	}
}

type StoryUpdate struct {
	TimeUpdated *time.Time `bson:"time_updated,omitempty"`
}

type StoryFilters struct {
	Ids       []string
	CreatedBy *string
	User      *usersstore.User
}

func (s StorySortBy) String() string {
	switch s {
	case StorySortByTimeCreated:
		return "time_created"
	default:
		return "time_created"
	}
}

func GetStorySortByFromString(s string) StorySortBy {
	switch s {
	case "time_created":
		return StorySortByTimeCreated
	default:
		return StorySortByDefault
	}
}

func (s StorySortBy) CursorType() cursor.ValueType {
	switch s {
	case StorySortByTimeCreated:
		return cursor.ValueTypeTime
	default:
		return cursor.ValueTypeTime
	}
}
