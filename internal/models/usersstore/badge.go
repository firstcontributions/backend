package usersstore

import "time"

type BadgeSortBy uint8

const (
	BadgeSortByDefault = iota
	BadgeSortByTimeCreated
)

type Badge struct {
	UserID                        string    `bson:"user_id"`
	CurrentLevel                  int64     `bson:"current_level,omitempty"`
	DisplayName                   string    `bson:"display_name,omitempty"`
	Id                            string    `bson:"_id"`
	Points                        int64     `bson:"points,omitempty"`
	ProgressPercentageToNextLevel int64     `bson:"progress_percentage_to_next_level,omitempty"`
	TimeCreated                   time.Time `bson:"time_created,omitempty"`
	TimeUpdated                   time.Time `bson:"time_updated,omitempty"`
}

func NewBadge() *Badge {
	return &Badge{}
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
	case BadgeSortByTimeCreated:
		return "time_created"
	default:
		return "time_created"
	}
}

func GetBadgeSortByFromString(s string) BadgeSortBy {
	switch s {
	case "time_created":
		return BadgeSortByTimeCreated
	default:
		return BadgeSortByDefault
	}
}
