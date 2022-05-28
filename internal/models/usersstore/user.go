package usersstore

import "time"

type User struct {
	CursorCheckpoints *CursorCheckpoints `bson:"cursor_checkpoints,omitempty"`
	Handle            string             `bson:"handle,omitempty"`
	Id                string             `bson:"_id"`
	Name              string             `bson:"name,omitempty"`
	Tags              *Tags              `bson:"tags,omitempty"`
	TimeCreated       time.Time          `bson:"time_created,omitempty"`
	TimeUpdated       time.Time          `bson:"time_updated,omitempty"`
	Token             *Token             `bson:"token,omitempty"`
}

func NewUser() *User {
	return &User{}
}

type UserUpdate struct {
	CursorCheckpoints *CursorCheckpoints `bson:"cursor_checkpoints,omitempty"`
	Name              *string            `bson:"name,omitempty"`
	Tags              *Tags              `bson:"tags,omitempty"`
	TimeUpdated       *time.Time         `bson:"time_updated,omitempty"`
	Token             *Token             `bson:"token,omitempty"`
}
