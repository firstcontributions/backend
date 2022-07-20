package usersstore

import "github.com/firstcontributions/backend/pkg/cursor"

type GitContributionStatsSortBy uint8

const (
	GitContributionStatsSortByDefault GitContributionStatsSortBy = iota
)

type GitContributionStats struct {
	Issues       int64 `bson:"issues,omitempty"`
	PullRequests int64 `bson:"pull_requests,omitempty"`
}

func NewGitContributionStats() *GitContributionStats {
	return &GitContributionStats{}
}
func (git_contribution_stats *GitContributionStats) Get(field string) interface{} {
	switch field {
	case "issues":
		return git_contribution_stats.Issues
	case "pull_requests":
		return git_contribution_stats.PullRequests
	default:
		return nil
	}
}

type GitContributionStatsFilters struct {
	Ids []string
}

func (s GitContributionStatsSortBy) String() string {
	switch s {
	default:
		return "time_created"
	}
}

func GetGitContributionStatsSortByFromString(s string) GitContributionStatsSortBy {
	switch s {
	default:
		return GitContributionStatsSortByDefault
	}
}

func (s GitContributionStatsSortBy) CursorType() cursor.ValueType {
	switch s {
	default:
		return cursor.ValueTypeTime
	}
}
