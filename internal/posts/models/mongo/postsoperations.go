package mongo

import (
	"context"

	"github.com/firstcontributions/backend/internal/posts/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Store) CreatePost(ctx context.Context, userID string, post *proto.PostData) (*proto.Post, error) {
	postModel := NewPost()
	if err := postModel.fromProto(post); err != nil {
		return nil, err
	}
	postModel.CreatedBy = userID
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

func (s *Store) Feeds(ctx context.Context, limit, offset int64) ([]*proto.Post, error) {

	res := []Post{}
	query := bson.M{}
	options := options.FindOptions{
		Limit: &limit,
		Skip:  &offset,
	}
	c, err := s.Collection(CollectionPosts).Find(ctx, query, &options)
	if err != nil {
		return nil, err
	}
	if err := c.All(ctx, &res); err != nil {
		return nil, err
	}
	posts := []*proto.Post{}

	for _, p := range res {
		posts = append(posts, p.toProto())
	}
	return posts, nil
}
