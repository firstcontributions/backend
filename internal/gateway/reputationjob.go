package gateway

import (
	"context"

	"github.com/firstcontributions/backend/internal/gateway/session"
	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/firstcontributions/backend/internal/reputation"
)

func (s *Server) UpdateProfileReputation(user *usersstore.User, sessionID string) {
	rs := reputation.NewReputationSynchroniser(*s.GithubConfig, s.Store.UsersStore)
	go func() {
		ctx := session.WithContext(context.Background(), session.NewMetaData(user))
		rs.SyncBadges(ctx, user)
		rs.SyncTags(ctx, user)
		rs.SyncContributionStats(ctx, user)
	}()
}
