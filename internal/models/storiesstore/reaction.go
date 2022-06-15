package storiesstore

import "time"

type Reaction struct {
	CommentID   string    `bson:"comment_id"`
	StoryID     string    `bson:"story_id"`
	CreatedBy   string    `bson:"created_by,omitempty"`
	Id          string    `bson:"_id"`
	TimeCreated time.Time `bson:"time_created,omitempty"`
	TimeUpdated time.Time `bson:"time_updated,omitempty"`
}

func NewReaction() *Reaction {
	return &Reaction{}
}

type ReactionUpdate struct {
	TimeUpdated *time.Time `bson:"time_updated,omitempty"`
}
