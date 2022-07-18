package storiesstore

import "context"

type Store interface {
	// story methods
	GetStoryByID(ctx context.Context, id string) (*Story, error)
	GetOneStory(ctx context.Context, filters *StoryFilters) (*Story, error)
	GetStories(ctx context.Context, filters *StoryFilters, after *string, before *string, first *int64, last *int64, sortBy StorySortBy, sortOrder *string) ([]*Story, bool, bool, string, string, error)
	CountStories(ctx context.Context, filters *StoryFilters) (int64, error)
	CreateStory(ctx context.Context, story *Story) (*Story, error)
	UpdateStory(ctx context.Context, id string, update *StoryUpdate) error
	DeleteStoryByID(ctx context.Context, id string) error
	// comment methods
	GetCommentByID(ctx context.Context, id string) (*Comment, error)
	GetOneComment(ctx context.Context, filters *CommentFilters) (*Comment, error)
	GetComments(ctx context.Context, filters *CommentFilters, after *string, before *string, first *int64, last *int64, sortBy CommentSortBy, sortOrder *string) ([]*Comment, bool, bool, string, string, error)
	CountComments(ctx context.Context, filters *CommentFilters) (int64, error)
	CreateComment(ctx context.Context, comment *Comment) (*Comment, error)
	UpdateComment(ctx context.Context, id string, update *CommentUpdate) error
	DeleteCommentByID(ctx context.Context, id string) error
	// reaction methods
	GetReactionByID(ctx context.Context, id string) (*Reaction, error)
	GetOneReaction(ctx context.Context, filters *ReactionFilters) (*Reaction, error)
	GetReactions(ctx context.Context, filters *ReactionFilters, after *string, before *string, first *int64, last *int64, sortBy ReactionSortBy, sortOrder *string) ([]*Reaction, bool, bool, string, string, error)
	CountReactions(ctx context.Context, filters *ReactionFilters) (int64, error)
	CreateReaction(ctx context.Context, reaction *Reaction) (*Reaction, error)
	UpdateReaction(ctx context.Context, id string, update *ReactionUpdate) error
	DeleteReactionByID(ctx context.Context, id string) error
}
