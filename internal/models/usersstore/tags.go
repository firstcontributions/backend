package usersstore

type TagsSortBy uint8

const (
	TagsSortByDefault = iota
	TagsSortByTimeCreated
)

type Tags struct {
	Languages   []*string `bson:"languages,omitempty"`
	RecentRepos []*string `bson:"recent_repos,omitempty"`
	Topics      []*string `bson:"topics,omitempty"`
}

func NewTags() *Tags {
	return &Tags{}
}

type TagsFilters struct {
	Ids []string
}

func (s TagsSortBy) String() string {
	switch s {
	case TagsSortByTimeCreated:
		return "time_created"
	default:
		return "time_created"
	}
}

func GetTagsSortByFromString(s string) TagsSortBy {
	switch s {
	case "time_created":
		return TagsSortByTimeCreated
	default:
		return TagsSortByDefault
	}
}
