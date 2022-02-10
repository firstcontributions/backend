package usersstore

import "time"

type CursorCheckpoints struct {
	PullRequests string    `bson:"pull_requests,omitempty"`
	TimeCreated  time.Time `bson:"time_created,omitempty"`
	TimeUpdated  time.Time `bson:"time_updated,omitempty"`
}

func NewCursorCheckpoints() *CursorCheckpoints {
	return &CursorCheckpoints{}
}

type CursorCheckpointsUpdate struct {
	TimeUpdated *time.Time `bson:"time_updated,omitempty"`
}
