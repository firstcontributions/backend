package usersstore

import (
	"github.com/firstcontributions/backend/pkg/authorizer"
	"github.com/firstcontributions/backend/pkg/cursor"
)

type TagsSortBy uint8

const (
	TagsSortByDefault TagsSortBy = iota
)

type Tags struct {
	Languages   []*string         `bson:"languages,omitempty"`
	RecentRepos []*string         `bson:"recent_repos,omitempty"`
	Topics      []*string         `bson:"topics,omitempty"`
	Ownership   *authorizer.Scope `bson:"ownership,omitempty"`
}

func NewTags() *Tags {
	return &Tags{}
}
func (tags *Tags) Get(field string) interface{} {
	switch field {
	case "languages":
		return tags.Languages
	case "recent_repos":
		return tags.RecentRepos
	case "topics":
		return tags.Topics
	default:
		return nil
	}
}

type TagsFilters struct {
	Ids []string
}

func (s TagsSortBy) String() string {
	switch s {
	default:
		return "time_created"
	}
}

func GetTagsSortByFromString(s string) TagsSortBy {
	switch s {
	default:
		return TagsSortByDefault
	}
}

func (s TagsSortBy) CursorType() cursor.ValueType {
	switch s {
	default:
		return cursor.ValueTypeTime
	}
}
