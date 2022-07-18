package usersstore

type ReputationSortBy uint8

const (
	ReputationSortByDefault = iota
	ReputationSortByTimeCreated
)

type Reputation struct {
	ContributionsToPopularRepos   int64   `bson:"contributions_to_popular_repos,omitempty"`
	ContributionsToUnpopularRepos int64   `bson:"contributions_to_unpopular_repos,omitempty"`
	Value                         float64 `bson:"value,omitempty"`
}

func NewReputation() *Reputation {
	return &Reputation{}
}

type ReputationFilters struct {
	Ids []string
}

func (s ReputationSortBy) String() string {
	switch s {
	case ReputationSortByTimeCreated:
		return "time_created"
	default:
		return "time_created"
	}
}

func GetReputationSortByFromString(s string) ReputationSortBy {
	switch s {
	case "time_created":
		return ReputationSortByTimeCreated
	default:
		return ReputationSortByDefault
	}
}
