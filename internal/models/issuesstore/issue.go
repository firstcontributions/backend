package issuesstore

import (
	"time"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/firstcontributions/backend/pkg/authorizer"
	"github.com/firstcontributions/backend/pkg/cursor"
)

type IssueSortBy uint8

const (
	IssueSortByDefault IssueSortBy = iota
	IssueSortByRepositoryUpdatedAt
)

type Issue struct {
	UserID              string            `bson:"user_id"`
	Body                string            `bson:"body,omitempty"`
	CommentCount        int64             `bson:"comment_count,omitempty"`
	Id                  string            `bson:"_id"`
	IssueType           string            `bson:"issue_type,omitempty"`
	Labels              []*string         `bson:"labels,omitempty"`
	Repository          string            `bson:"repository,omitempty"`
	RepositoryAvatar    string            `bson:"repository_avatar,omitempty"`
	RepositoryUpdatedAt time.Time         `bson:"repository_updated_at,omitempty"`
	Title               string            `bson:"title,omitempty"`
	Url                 string            `bson:"url,omitempty"`
	Ownership           *authorizer.Scope `bson:"ownership,omitempty"`
}

func NewIssue() *Issue {
	return &Issue{}
}
func (issue *Issue) Get(field string) interface{} {
	switch field {
	case "user_id":
		return issue.UserID
	case "body":
		return issue.Body
	case "comment_count":
		return issue.CommentCount
	case "_id":
		return issue.Id
	case "issue_type":
		return issue.IssueType
	case "labels":
		return issue.Labels
	case "repository":
		return issue.Repository
	case "repository_avatar":
		return issue.RepositoryAvatar
	case "repository_updated_at":
		return issue.RepositoryUpdatedAt
	case "title":
		return issue.Title
	case "url":
		return issue.Url
	default:
		return nil
	}
}

type IssueFilters struct {
	Ids       []string
	IssueType *string
	User      *usersstore.User
}

func (s IssueSortBy) String() string {
	switch s {
	case IssueSortByRepositoryUpdatedAt:
		return "repository_updated_at"
	default:
		return "time_created"
	}
}

func GetIssueSortByFromString(s string) IssueSortBy {
	switch s {
	case "repository_updated_at":
		return IssueSortByRepositoryUpdatedAt
	default:
		return IssueSortByDefault
	}
}

func (s IssueSortBy) CursorType() cursor.ValueType {
	switch s {
	case IssueSortByRepositoryUpdatedAt:
		return cursor.ValueTypeTime
	default:
		return cursor.ValueTypeTime
	}
}
