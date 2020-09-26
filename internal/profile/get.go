package profile

import (
	"context"
	"log"

	"github.com/firstcontributions/firstcontributions/internal/proto"
)

// GetProfile implemets GetProfile RPC call, will get the profile by reference
func (s *Service) GetProfile(ctx context.Context, req *proto.GetProfileRequest) (*proto.Profile, error) {
	log.Printf("%#v", req.String())
	return &proto.Profile{
		Name: "Gokul",
	}, nil
}
