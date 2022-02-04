package usersstore

import "time"

type User struct {
	Handle      string    `bson:"handle,omitempty"`
	Id          string    `bson:"_id"`
	Name        string    `bson:"name,omitempty"`
	TimeCreated time.Time `bson:"time_created,omitempty"`
	TimeUpdated time.Time `bson:"time_updated,omitempty"`
	Token       *Token    `bson:"token,omitempty"`
}

func NewUser() *User {
	return &User{}
}

type UserUpdate struct {
	Handle      *string    `bson:"handle"`
	Name        *string    `bson:"name"`
	TimeUpdated *time.Time `bson:"time_updated"`
}
