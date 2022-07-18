package usersstore

import "time"

type UserSortBy uint8

const (
	UserSortByDefault = iota
	UserSortByTimeCreated
)

type User struct {
	Avatar               string                `bson:"avatar,omitempty"`
	Bio                  string                `bson:"bio,omitempty"`
	CursorCheckpoints    *CursorCheckpoints    `bson:"cursor_checkpoints,omitempty"`
	GitContributionStats *GitContributionStats `bson:"git_contribution_stats,omitempty"`
	Handle               string                `bson:"handle,omitempty"`
	Id                   string                `bson:"_id"`
	Name                 string                `bson:"name,omitempty"`
	Reputation           *Reputation           `bson:"reputation,omitempty"`
	Tags                 *Tags                 `bson:"tags,omitempty"`
	TimeCreated          time.Time             `bson:"time_created,omitempty"`
	TimeUpdated          time.Time             `bson:"time_updated,omitempty"`
	Token                *Token                `bson:"token,omitempty"`
}

func NewUser() *User {
	return &User{}
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
