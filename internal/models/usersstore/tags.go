package usersstore

type Tags struct {
	Languages   []string `bson:"languages,omitempty"`
	RecentRepos []string `bson:"recent_repos,omitempty"`
	Topics      []string `bson:"topics,omitempty"`
}

func NewTags() *Tags {
	return &Tags{}
}
