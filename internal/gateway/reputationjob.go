package gateway

import (
	"github.com/firstcontributions/backend/internal/models/usersstore"
)

func (s *Server) UpdateProfileReputation(user *usersstore.User) {
	// rs := reputation.NewReputationSynchroniser(*s.GithubConfig, s.Store.UsersStore)
	// ctx := session.WithContext(context.Background(), session.NewMetaData(user))
	// go rs.SyncBadges(ctx, user)
}
