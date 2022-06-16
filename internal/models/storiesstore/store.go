package storiesstore

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/usersstore"
)

type Store interface {
	// comment methods
	GetCommentByID(ctx context.Context, id string) (*Comment, error)
	GetComments(ctx context.Context, ids []string,
		story *Story, after *string, before *string, first *int64, last *int64) ([]*Comment, bool, bool, string, string, error)
	CreateComment(ctx context.Context, comment *Comment) (*Comment, error)
	UpdateComment(ctx context.Context, id string, update *CommentUpdate) error
	DeleteCommentByID(ctx context.Context, id string) error
	// reaction methods
	GetReactionByID(ctx context.Context, id string) (*Reaction, error)
	GetReactions(ctx context.Context, ids []string,
		comment *Comment,
		story *Story, after *string, before *string, first *int64, last *int64) ([]*Reaction, bool, bool, string, string, error)
	CreateReaction(ctx context.Context, reaction *Reaction) (*Reaction, error)
	UpdateReaction(ctx context.Context, id string, update *ReactionUpdate) error
	DeleteReactionByID(ctx context.Context, id string) error
	// story methods
	GetStoryByID(ctx context.Context, id string) (*Story, error)
	GetStories(ctx context.Context, ids []string,
		user *usersstore.User, after *string, before *string, first *int64, last *int64) ([]*Story, bool, bool, string, string, error)
	CreateStory(ctx context.Context, story *Story) (*Story, error)
	UpdateStory(ctx context.Context, id string, update *StoryUpdate) error
	DeleteStoryByID(ctx context.Context, id string) error
}
