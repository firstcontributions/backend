package reputation

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/shurcooL/githubv4"
)

type ContributionStatsQuery struct {
	Viewer struct {
		PullRequests struct {
			TotalCount githubv4.Int
		}
		Issues struct {
			TotalCount githubv4.Int
		}
	}
}

func (r *ReputationSynchroniser) SyncContributionStats(ctx context.Context, user *usersstore.User) error {
	query := &ContributionStatsQuery{}
	if err := r.Query(ctx, query, nil); err != nil {
		return err
	}
	update := &usersstore.UserUpdate{
		GitContributionStats: &usersstore.GitContributionStats{
			PullRequests: int64(query.Viewer.PullRequests.TotalCount),
			Issues:       int64(query.Viewer.Issues.TotalCount),
		},
	}
	return r.userStore.UpdateUser(ctx, user.Id, update)
}
