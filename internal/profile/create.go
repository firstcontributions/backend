package profile

import (
	"context"
	"log"

	"github.com/firstcontributions/firstcontributions/internal/profile/models/mongo"
	"github.com/firstcontributions/firstcontributions/internal/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (s *Service) CreateProfile(ctx context.Context, req *proto.Profile) (*proto.Profile, error) {
	log.Printf("creating profile [handle: %s]", req.Handle)
	id, err := uuid.NewUUID()
	if err != nil {
		log.Printf("error on creating uuid in create profile %v", err)
		return nil, grpc.Errorf(codes.Internal, "error on creating uuid  %w", err)
	}
	req.Uuid = id.String()
	if err := mongo.CreateProfile(ctx, s.MongoClient, req); err != nil {
		log.Printf("error on creating profile %v", err)
		return nil, grpc.Errorf(codes.Internal, "error on creating profile  %w", err)
	}
	return req, nil
}
