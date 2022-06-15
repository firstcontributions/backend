package storiesstore

import "time"

type Story struct {
	UserID          string    `bson:"user_id"`
	AbstractContent string    `bson:"abstract_content,omitempty"`
	ContentJson     string    `bson:"content_json,omitempty"`
	CreatedBy       string    `bson:"created_by,omitempty"`
	Id              string    `bson:"_id"`
	Thumbanil       string    `bson:"thumbanil,omitempty"`
	TimeCreated     time.Time `bson:"time_created,omitempty"`
	TimeUpdated     time.Time `bson:"time_updated,omitempty"`
	UrlSuffix       string    `bson:"url_suffix,omitempty"`
}

func NewStory() *Story {
	return &Story{}
}

type StoryUpdate struct {
	TimeUpdated *time.Time `bson:"time_updated,omitempty"`
}
