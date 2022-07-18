package storiesstore

import "time"

type ReactionSortBy uint8

const (
	ReactionSortByDefault = iota
	ReactionSortByTimeCreated
)

type Reaction struct {
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

type ReactionFilters struct {
	Ids       []string
	CreatedBy *string
	Story     *Story
}

func (s ReactionSortBy) String() string {
	switch s {
	case ReactionSortByTimeCreated:
		return "time_created"
	default:
		return "time_created"
	}
}

func GetReactionSortByFromString(s string) ReactionSortBy {
	switch s {
	case "time_created":
		return ReactionSortByTimeCreated
	default:
		return ReactionSortByDefault
	}
}
