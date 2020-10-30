package mongo

import (
	"context"

	"github.com/firstcontributions/backend/internal/posts/proto"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Store) CreatePost(ctx context.Context, post *proto.PostData) (*proto.Post, error) {
	postModel := NewPost()
	if err := postModel.fromProto(post); err != nil {
		return nil, err
	}
	if _, err := s.Collection(CollectionPosts).InsertOne(ctx, postModel); err != nil {
		return nil, err
	}
	return postModel.toProto(), nil
}

func (s *Store) GetPostByUUID(ctx context.Context, uuid string) (*proto.Post, error) {
	post := NewPost()
	query := bson.M{
		"_id": uuid,
	}
	if err := s.Collection(CollectionPosts).FindOne(ctx, query).Decode(post); err != nil {
		return nil, err
	}
	return post.toProto(), nil
}
