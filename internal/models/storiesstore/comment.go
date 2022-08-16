package storiesstore

import (
	"time"

	"github.com/firstcontributions/backend/pkg/authorizer"
	"github.com/firstcontributions/backend/pkg/cursor"
)

type CommentSortBy uint8

const (
	CommentSortByDefault CommentSortBy = iota
	CommentSortByTimeCreated
)

type Comment struct {
	StoryID         string            `bson:"story_id"`
	AbstractContent string            `bson:"abstract_content,omitempty"`
	ContentJson     string            `bson:"content_json,omitempty"`
	CreatedBy       string            `bson:"created_by,omitempty"`
	Id              string            `bson:"_id"`
	TimeCreated     time.Time         `bson:"time_created,omitempty"`
	TimeUpdated     time.Time         `bson:"time_updated,omitempty"`
	Ownership       *authorizer.Scope `bson:"ownership,omitempty"`
}

func NewComment() *Comment {
	return &Comment{}
}
func (comment *Comment) Get(field string) interface{} {
	switch field {
	case "story_id":
		return comment.StoryID
	case "abstract_content":
		return comment.AbstractContent
	case "content_json":
		return comment.ContentJson
	case "created_by":
		return comment.CreatedBy
	case "_id":
		return comment.Id
	case "time_created":
		return comment.TimeCreated
	case "time_updated":
		return comment.TimeUpdated
	default:
		return nil
	}
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

func (s CommentSortBy) CursorType() cursor.ValueType {
	switch s {
	case CommentSortByTimeCreated:
		return cursor.ValueTypeTime
	default:
		return cursor.ValueTypeTime
	}
}
