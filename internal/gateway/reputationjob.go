package gateway

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/firstcontributions/backend/internal/reputation"
)

func (s *Server) UpdateProfileReputation(user *usersstore.User) {
	rs := reputation.NewReputationSynchroniser(*s.GithubConfig, s.Store.UsersStore)

	go rs.SyncBadges(context.Background(), user)
}
