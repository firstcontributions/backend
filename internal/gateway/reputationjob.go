package gateway

import (
	"context"

	"github.com/firstcontributions/backend/internal/profile/proto"
)

func (s *Server) UpdateProfileReputation(ctx context.Context, profile *proto.Profile) {
	profile.Reputation++
	s.ProfileManager.UpdateProfile(ctx, profile)
}
