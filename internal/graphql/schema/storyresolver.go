package schema

import (
	"context"
	"errors"

	"github.com/firstcontributions/backend/internal/gateway/session"
	"github.com/firstcontributions/backend/internal/models/storiesstore"
	"github.com/firstcontributions/backend/internal/storemanager"
	graphql "github.com/graph-gophers/graphql-go"
)

type Story struct {
	ref             *storiesstore.Story
	AbstractContent string
	ContentJson     string
	createdBy       string
	Id              string
	Thumbnail       string
	TimeCreated     graphql.Time
	TimeUpdated     graphql.Time
	Title           string
	UrlSuffix       string
}

func NewStory(m *storiesstore.Story) *Story {
	if m == nil {
		return nil
	}
	return &Story{
		ref:             m,
		AbstractContent: m.AbstractContent,
		ContentJson:     m.ContentJson,
		createdBy:       m.CreatedBy,
		Id:              m.Id,
		Thumbnail:       m.Thumbnail,
		TimeCreated:     graphql.Time{Time: m.TimeCreated},
		TimeUpdated:     graphql.Time{Time: m.TimeUpdated},
		Title:           m.Title,
		UrlSuffix:       m.UrlSuffix,
	}
}
func (n *Story) CreatedBy(ctx context.Context) (*User, error) {

	data, err := storemanager.FromContext(ctx).UsersStore.GetUserByID(ctx, n.createdBy)
	if err != nil {
		return nil, err
	}
	return NewUser(data), nil
}

type CreateStoryInput struct {
	AbstractContent string
	ContentJson     string
	Thumbnail       string
	Title           string
	UrlSuffix       string
	UserID          graphql.ID
}

func (n *CreateStoryInput) ToModel() (*storiesstore.Story, error) {
	if n == nil {
		return nil, nil
	}

	return &storiesstore.Story{
		AbstractContent: n.AbstractContent,
		ContentJson:     n.ContentJson,
		Thumbnail:       n.Thumbnail,
		Title:           n.Title,
		UrlSuffix:       n.UrlSuffix,
	}, nil
}

type UpdateStoryInput struct {
	ID graphql.ID
}

func (n *UpdateStoryInput) ToModel() *storiesstore.StoryUpdate {
	if n == nil {
		return nil
	}
	return &storiesstore.StoryUpdate{}
}
func (n *Story) ID(ctx context.Context) graphql.ID {
	return NewIDMarshaller(NodeTypeStory, n.Id, true).
		ToGraphqlID()
}

type StoriesConnection struct {
	Edges    []*StoryEdge
	PageInfo *PageInfo
	filters  *storiesstore.StoryFilters
}

func NewStoriesConnection(
	filters *storiesstore.StoryFilters,
	data []*storiesstore.Story,
	hasNextPage bool,
	hasPreviousPage bool,
	cursors []string,
) *StoriesConnection {
	edges := []*StoryEdge{}
	for i, d := range data {
		node := NewStory(d)

		edges = append(edges, &StoryEdge{
			Node:   node,
			Cursor: cursors[i],
		})
	}
	var startCursor, endCursor *string
	if len(cursors) > 0 {
		startCursor = &cursors[0]
		endCursor = &cursors[len(cursors)-1]
	}
	return &StoriesConnection{
		filters: filters,
		Edges:   edges,
		PageInfo: &PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		},
	}
}

func (c StoriesConnection) TotalCount(ctx context.Context) (int32, error) {
	count, err := storemanager.FromContext(ctx).StoriesStore.CountStories(ctx, c.filters)
	return int32(count), err
}
func (c StoriesConnection) HasViewerAssociation(ctx context.Context) (bool, error) {
	session := session.FromContext(ctx)
	if session == nil {
		return false, errors.New("Unauthorized")
	}
	userID := session.UserID()

	newFilter := *c.filters
	newFilter.CreatedBy = &userID

	data, err := storemanager.FromContext(ctx).StoriesStore.GetOneStory(ctx, c.filters)
	if err != nil {
		return false, err
	}
	return data != nil, nil
}

type StoryEdge struct {
	Node   *Story
	Cursor string
}
