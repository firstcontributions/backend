package schema

import (
	"context"
	"errors"

	"github.com/firstcontributions/backend/internal/gateway/session"
	"github.com/firstcontributions/backend/internal/models/storiesstore"
	"github.com/firstcontributions/backend/internal/storemanager"
	"github.com/firstcontributions/backend/pkg/cursor"
	graphql "github.com/graph-gophers/graphql-go"
)

type Comment struct {
	ref             *storiesstore.Comment
	AbstractContent string
	ContentJson     string
	createdBy       string
	Id              string
	TimeCreated     graphql.Time
	TimeUpdated     graphql.Time
}

func NewComment(m *storiesstore.Comment) *Comment {
	if m == nil {
		return nil
	}
	return &Comment{
		ref:             m,
		AbstractContent: m.AbstractContent,
		ContentJson:     m.ContentJson,
		createdBy:       m.CreatedBy,
		Id:              m.Id,
		TimeCreated:     graphql.Time{Time: m.TimeCreated},
		TimeUpdated:     graphql.Time{Time: m.TimeUpdated},
	}
}
func (n *Comment) CreatedBy(ctx context.Context) (*User, error) {

	data, err := storemanager.FromContext(ctx).UsersStore.GetUserByID(ctx, n.createdBy)
	if err != nil {
		return nil, err
	}
	return NewUser(data), nil
}

type CreateCommentInput struct {
	AbstractContent string
	ContentJson     string
	StoryID         graphql.ID
}

func (n *CreateCommentInput) ToModel() (*storiesstore.Comment, error) {
	if n == nil {
		return nil, nil
	}
	storyID, err := ParseGraphqlID(n.StoryID)
	if err != nil {
		return nil, err
	}

	return &storiesstore.Comment{
		AbstractContent: n.AbstractContent,
		ContentJson:     n.ContentJson,
		StoryID:         storyID.ID,
	}, nil
}

type UpdateCommentInput struct {
	ID graphql.ID
}

func (n *UpdateCommentInput) ToModel() *storiesstore.CommentUpdate {
	if n == nil {
		return nil
	}
	return &storiesstore.CommentUpdate{}
}
func (n *Comment) ID(ctx context.Context) graphql.ID {
	return NewIDMarshaller("comment", n.Id).
		ToGraphqlID()
}

type CommentsConnection struct {
	Edges    []*CommentEdge
	PageInfo *PageInfo
	filters  *storiesstore.CommentFilters
}

func NewCommentsConnection(
	filters *storiesstore.CommentFilters,
	data []*storiesstore.Comment,
	hasNextPage bool,
	hasPreviousPage bool,
	firstCursor *string,
	lastCursor *string,
) *CommentsConnection {
	edges := []*CommentEdge{}
	for _, d := range data {
		node := NewComment(d)

		edges = append(edges, &CommentEdge{
			Node:   node,
			Cursor: cursor.NewCursor(d.Id, "time_created", d.TimeCreated).String(),
		})
	}
	return &CommentsConnection{
		filters: filters,
		Edges:   edges,
		PageInfo: &PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
			StartCursor:     firstCursor,
			EndCursor:       lastCursor,
		},
	}
}

func (c CommentsConnection) TotalCount(ctx context.Context) (int32, error) {
	count, err := storemanager.FromContext(ctx).StoriesStore.CountComments(ctx, c.filters)
	return int32(count), err
}
func (c CommentsConnection) HasViewerAssociation(ctx context.Context) (bool, error) {
	session := session.FromContext(ctx)
	if session == nil {
		return false, errors.New("Unauthorized")
	}
	userID := session.UserID()

	newFilter := *c.filters
	newFilter.CreatedBy = &userID

	data, err := storemanager.FromContext(ctx).StoriesStore.GetOneComment(ctx, c.filters)
	if err != nil {
		return false, err
	}
	return data != nil, nil
}

type CommentEdge struct {
	Node   *Comment
	Cursor string
}
