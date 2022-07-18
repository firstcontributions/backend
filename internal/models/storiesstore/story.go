package storiesstore

import (
	"time"

	"github.com/firstcontributions/backend/internal/models/usersstore"
)

type StorySortBy uint8

const (
	StorySortByDefault = iota
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
