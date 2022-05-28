package gateway

import (
	"context"
	"fmt"

	"github.com/firstcontributions/backend/internal/gateway/session"
	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/firstcontributions/backend/internal/reputation"
)

func (s *Server) UpdateProfileReputation(user *usersstore.User, sessionID string) {
	rs := reputation.NewReputationSynchroniser(*s.GithubConfig, s.Store.UsersStore)
	go func() {
		fmt.Println("---------------", user)
		ctx := session.WithContext(context.Background(), session.NewMetaData(user))
		meta := session.FromContext(ctx)
		fmt.Println("-------------- meta", meta)
		rs.SyncBadges(ctx, user)
		rs.SyncTags(ctx, user)
		updatedUser, _ := s.Store.UsersStore.GetUserByID(ctx, user.Id)
		s.SessionManager.Set(ctx, sessionID, session.NewMetaData(updatedUser))
	}()
}
