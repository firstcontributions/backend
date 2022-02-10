package usersstore

import "time"

type Badge struct {
	UserID      string    `bson:"user_id"`
	DisplayName string    `bson:"display_name,omitempty"`
	Id          string    `bson:"_id"`
	TimeCreated time.Time `bson:"time_created,omitempty"`
	TimeUpdated time.Time `bson:"time_updated,omitempty"`
}

func NewBadge() *Badge {
	return &Badge{}
}
