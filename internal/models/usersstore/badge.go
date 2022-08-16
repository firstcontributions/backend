package usersstore

import (
	"time"

	"github.com/firstcontributions/backend/pkg/authorizer"
	"github.com/firstcontributions/backend/pkg/cursor"
)

type BadgeSortBy uint8

const (
	BadgeSortByDefault BadgeSortBy = iota
	BadgeSortByPoints
	BadgeSortByTimeCreated
)

type Badge struct {
	UserID                        string            `bson:"user_id"`
	CurrentLevel                  int64             `bson:"current_level,omitempty"`
	DisplayName                   string            `bson:"display_name,omitempty"`
	Id                            string            `bson:"_id"`
	Points                        int64             `bson:"points,omitempty"`
	ProgressPercentageToNextLevel int64             `bson:"progress_percentage_to_next_level,omitempty"`
	TimeCreated                   time.Time         `bson:"time_created,omitempty"`
	TimeUpdated                   time.Time         `bson:"time_updated,omitempty"`
	Ownership                     *authorizer.Scope `bson:"ownership,omitempty"`
}

func NewBadge() *Badge {
	return &Badge{}
}
func (badge *Badge) Get(field string) interface{} {
	switch field {
	case "user_id":
		return badge.UserID
	case "current_level":
		return badge.CurrentLevel
	case "display_name":
		return badge.DisplayName
	case "_id":
		return badge.Id
	case "points":
		return badge.Points
	case "progress_percentage_to_next_level":
		return badge.ProgressPercentageToNextLevel
	case "time_created":
		return badge.TimeCreated
	case "time_updated":
		return badge.TimeUpdated
	default:
		return nil
	}
}

type BadgeUpdate struct {
	CurrentLevel                  *int64     `bson:"current_level,omitempty"`
	Points                        *int64     `bson:"points,omitempty"`
	ProgressPercentageToNextLevel *int64     `bson:"progress_percentage_to_next_level,omitempty"`
	TimeUpdated                   *time.Time `bson:"time_updated,omitempty"`
}

type BadgeFilters struct {
	Ids []string

	User *User
}

func (s BadgeSortBy) String() string {
	switch s {
	case BadgeSortByPoints:
		return "points"
	case BadgeSortByTimeCreated:
		return "time_created"
	default:
		return "time_created"
	}
}

func GetBadgeSortByFromString(s string) BadgeSortBy {
	switch s {
	case "points":
		return BadgeSortByPoints
	case "time_created":
		return BadgeSortByTimeCreated
	default:
		return BadgeSortByDefault
	}
}

func (s BadgeSortBy) CursorType() cursor.ValueType {
	switch s {
	case BadgeSortByPoints:
		return cursor.ValueTypeInt
	case BadgeSortByTimeCreated:
		return cursor.ValueTypeTime
	default:
		return cursor.ValueTypeTime
	}
}
