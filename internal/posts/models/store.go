package models

import (
	"context"

	"github.com/firstcontributions/backend/internal/posts/proto"
)

// Store defines the interface for datastores
type Store interface {
	CreatePost(context.Context, *proto.PostData) (*proto.Post, error)
	GetPostByUUID(context.Context, string) (*proto.Post, error)
	UpdatePost(context.Context, *proto.Post) (*proto.Post, error)

	CreateComment(context.Context, *proto.CommentData) (*proto.Comment, error)
	GetCommentByUUID(context.Context, *proto.Comment) (*proto.Comment, error)
	CreateClap(context.Context, *proto.ClapRequst) (*proto.Clap, error)
}

type StorageManager struct {
	Store
}

func NewStorageManager(store Store) *StorageManager {
	return &StorageManager{
		Store: store,
	}
}
