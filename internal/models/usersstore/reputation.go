package usersstore

import (
	"github.com/firstcontributions/backend/pkg/authorizer"
	"github.com/firstcontributions/backend/pkg/cursor"
)

type ReputationSortBy uint8

const (
	ReputationSortByDefault ReputationSortBy = iota
)

type Reputation struct {
	ContributionsToPopularRepos   int64             `bson:"contributions_to_popular_repos,omitempty"`
	ContributionsToUnpopularRepos int64             `bson:"contributions_to_unpopular_repos,omitempty"`
	Value                         float64           `bson:"value,omitempty"`
	Ownership                     *authorizer.Scope `bson:"ownership,omitempty"`
}

func NewReputation() *Reputation {
	return &Reputation{}
}
func (reputation *Reputation) Get(field string) interface{} {
	switch field {
	case "contributions_to_popular_repos":
		return reputation.ContributionsToPopularRepos
	case "contributions_to_unpopular_repos":
		return reputation.ContributionsToUnpopularRepos
	case "value":
		return reputation.Value
	default:
		return nil
	}
}

type ReputationFilters struct {
	Ids []string
}

func (s ReputationSortBy) String() string {
	switch s {
	default:
		return "time_created"
	}
}

func GetReputationSortByFromString(s string) ReputationSortBy {
	switch s {
	default:
		return ReputationSortByDefault
	}
}

func (s ReputationSortBy) CursorType() cursor.ValueType {
	switch s {
	default:
		return cursor.ValueTypeTime
	}
}
