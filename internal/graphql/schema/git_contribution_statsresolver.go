package schema

import (
	"github.com/firstcontributions/backend/internal/models/usersstore"
)

type GitContributionStats struct {
	ref          *usersstore.GitContributionStats
	Issues       int32
	PullRequests int32
}

func NewGitContributionStats(m *usersstore.GitContributionStats) *GitContributionStats {
	if m == nil {
		return nil
	}
	return &GitContributionStats{
		ref:          m,
		Issues:       int32(m.Issues),
		PullRequests: int32(m.PullRequests),
	}
}
func (n *GitContributionStats) ToModel() *usersstore.GitContributionStats {
	if n == nil {
		return nil
	}
	return &usersstore.GitContributionStats{
		Issues:       int64(n.Issues),
		PullRequests: int64(n.PullRequests),
	}
}
