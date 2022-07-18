package usersstore

type GitContributionStatsSortBy uint8

const (
	GitContributionStatsSortByDefault = iota
	GitContributionStatsSortByTimeCreated
)

type GitContributionStats struct {
	Issues       int64 `bson:"issues,omitempty"`
	PullRequests int64 `bson:"pull_requests,omitempty"`
}

func NewGitContributionStats() *GitContributionStats {
	return &GitContributionStats{}
}

type GitContributionStatsFilters struct {
	Ids []string
}

func (s GitContributionStatsSortBy) String() string {
	switch s {
	case GitContributionStatsSortByTimeCreated:
		return "time_created"
	default:
		return "time_created"
	}
}

func GetGitContributionStatsSortByFromString(s string) GitContributionStatsSortBy {
	switch s {
	case "time_created":
		return GitContributionStatsSortByTimeCreated
	default:
		return GitContributionStatsSortByDefault
	}
}
