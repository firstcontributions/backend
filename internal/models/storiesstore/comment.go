package storiesstore

import "time"

type CommentSortBy uint8

const (
	CommentSortByDefault = iota
	CommentSortByTimeCreated
)

type Comment struct {
	StoryID         string    `bson:"story_id"`
	AbstractContent string    `bson:"abstract_content,omitempty"`
	ContentJson     string    `bson:"content_json,omitempty"`
	CreatedBy       string    `bson:"created_by,omitempty"`
	Id              string    `bson:"_id"`
	TimeCreated     time.Time `bson:"time_created,omitempty"`
	TimeUpdated     time.Time `bson:"time_updated,omitempty"`
}

func NewComment() *Comment {
	return &Comment{}
}

type CommentUpdate struct {
	TimeUpdated *time.Time `bson:"time_updated,omitempty"`
}

type CommentFilters struct {
	Ids       []string
	CreatedBy *string
	Story     *Story
}

func (s CommentSortBy) String() string {
	switch s {
	case CommentSortByTimeCreated:
		return "time_created"
	default:
		return "time_created"
	}
}

func GetCommentSortByFromString(s string) CommentSortBy {
	switch s {
	case "time_created":
		return CommentSortByTimeCreated
	default:
		return CommentSortByDefault
	}
}
