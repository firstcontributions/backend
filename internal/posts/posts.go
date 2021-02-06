package posts

import (
	"context"

	"log"

	"github.com/firstcontributions/backend/internal/posts/proto"
	"github.com/firstcontributions/backend/pkg/userctx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func (s *Service) CreatePost(ctx context.Context, post *proto.PostData) (*proto.Post, error) {
	p, err := s.PostManager.CreatePost(
		ctx,
		userctx.FromIncomingCtx(ctx).UserID(),
		post,
	)
	if err != nil {
		log.Printf("error on creating post %v", err)
		return nil, grpc.Errorf(codes.Internal, "error on creating post %w", err)
	}
	return p, nil
}

func (s *Service) GetPostByUUID(ctx context.Context, in *proto.GetByUUIDRequest) (*proto.Post, error) {
	p, err := s.PostManager.GetPostByUUID(ctx, in.Uuid)
	if err != nil {
		log.Printf("error on getting post from mongno %v", err)
		return nil, grpc.Errorf(codes.Internal, "error on getting post from mongno %w", err)
	}
	return p, nil
}

func (s *Service) UpdatePost(ctx context.Context, in *proto.Post) (*proto.Post, error) {
	return nil, nil
}
