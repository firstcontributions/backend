package usersstore

type Reputation struct {
	ContributionsToPopularRepos   int64   `bson:"contributions_to_popular_repos,omitempty"`
	ContributionsToUnpopularRepos int64   `bson:"contributions_to_unpopular_repos,omitempty"`
	Value                         float64 `bson:"value,omitempty"`
}

func NewReputation() *Reputation {
	return &Reputation{}
}
