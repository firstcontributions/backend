package usersstore

import (
	"time"

	"github.com/firstcontributions/backend/pkg/authorizer"
	"github.com/firstcontributions/backend/pkg/cursor"
)

type UserSortBy uint8

const (
	UserSortByDefault UserSortBy = iota
	UserSortByTimeCreated
)

type User struct {
	Avatar               string                  `bson:"avatar,omitempty"`
	Bio                  string                  `bson:"bio,omitempty"`
	CursorCheckpoints    *CursorCheckpoints      `bson:"cursor_checkpoints,omitempty"`
	GitContributionStats *GitContributionStats   `bson:"git_contribution_stats,omitempty"`
	Handle               string                  `bson:"handle,omitempty"`
	Id                   string                  `bson:"_id"`
	Name                 string                  `bson:"name,omitempty"`
	Reputation           *Reputation             `bson:"reputation,omitempty"`
	Tags                 *Tags                   `bson:"tags,omitempty"`
	TimeCreated          time.Time               `bson:"time_created,omitempty"`
	TimeUpdated          time.Time               `bson:"time_updated,omitempty"`
	Token                *Token                  `bson:"token,omitempty"`
	Permissions          []authorizer.Permission `bson:"permissions,omitempty"`
	Ownership            *authorizer.Scope       `bson:"ownership,omitempty"`
}

func NewUser() *User {
	return &User{}
}
func (user *User) Get(field string) interface{} {
	switch field {
	case "avatar":
		return user.Avatar
	case "bio":
		return user.Bio
	case "cursor_checkpoints":
		return user.CursorCheckpoints
	case "git_contribution_stats":
		return user.GitContributionStats
	case "handle":
		return user.Handle
	case "_id":
		return user.Id
	case "name":
		return user.Name
	case "reputation":
		return user.Reputation
	case "tags":
		return user.Tags
	case "time_created":
		return user.TimeCreated
	case "time_updated":
		return user.TimeUpdated
	case "token":
		return user.Token
	default:
		return nil
	}
}

type UserUpdate struct {
	Avatar               *string               `bson:"avatar,omitempty"`
	Bio                  *string               `bson:"bio,omitempty"`
	CursorCheckpoints    *CursorCheckpoints    `bson:"cursor_checkpoints,omitempty"`
	GitContributionStats *GitContributionStats `bson:"git_contribution_stats,omitempty"`
	Name                 *string               `bson:"name,omitempty"`
	Reputation           *Reputation           `bson:"reputation,omitempty"`
	Tags                 *Tags                 `bson:"tags,omitempty"`
	TimeUpdated          *time.Time            `bson:"time_updated,omitempty"`
	Token                *Token                `bson:"token,omitempty"`
}

type UserFilters struct {
	Ids    []string
	Search *string
	Handle *string
}

func (s UserSortBy) String() string {
	switch s {
	case UserSortByTimeCreated:
		return "time_created"
	default:
		return "time_created"
	}
}

func GetUserSortByFromString(s string) UserSortBy {
	switch s {
	case "time_created":
		return UserSortByTimeCreated
	default:
		return UserSortByDefault
	}
}

func (s UserSortBy) CursorType() cursor.ValueType {
	switch s {
	case UserSortByTimeCreated:
		return cursor.ValueTypeTime
	default:
		return cursor.ValueTypeTime
	}
}
