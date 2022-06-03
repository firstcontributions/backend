package usersstore

import "time"

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
