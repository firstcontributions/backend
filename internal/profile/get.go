package profile

import (
	"context"
	"log"

	"github.com/firstcontributions/backend/internal/profile/models/mongo"
	"github.com/firstcontributions/backend/internal/profile/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// GetProfile implemets GetProfile RPC call, will get the profile by reference
func (s *Service) GetProfile(ctx context.Context, req *proto.GetProfileRequest) (*proto.Profile, error) {
	data, err := mongo.GetProfileByHandle(ctx, s.MongoClient, req.Handle)
	if err != nil {
		log.Printf("error on getting profile %v", err)
		return nil, grpc.Errorf(codes.Internal, "error on mongo query %w", err)
	}
	if data == nil {
		return nil, grpc.Errorf(codes.NotFound, "profile does not found")
	}
	return data.Proto(), nil
}
