package posts

import (
	"context"

	"github.com/firstcontributions/backend/internal/posts/proto"
)

func (s *Service) CreateComment(ctx context.Context, comment *proto.CommentData) (*proto.Comment, error) {
	return nil, nil
}

func (s *Service) GetCommentByUUID(ctx context.Context, in *proto.GetByUUIDRequest) (*proto.Comment, error) {
	return nil, nil
}
