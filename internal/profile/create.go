package profile

import (
	"context"

	"github.com/firstcontributions/firstcontributions/internal/proto"
)

func (s *Service) CreateProfile(ctx context.Context, req *proto.Profile) (*proto.Profile, error) {
	return req, nil
}
