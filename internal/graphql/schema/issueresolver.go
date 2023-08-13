package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/issuesstore"
	"github.com/firstcontributions/backend/internal/storemanager"
	graphql "github.com/graph-gophers/graphql-go"
)

type Issue struct {
	ref                 *issuesstore.Issue
	Body                string
	CommentCount        int32
	Id                  string
	IssueType           string
	Labels              []*string
	Repository          string
	RepositoryAvatar    string
	RepositoryUpdatedAt graphql.Time
	Title               string
	Url                 string
}

func NewIssue(m *issuesstore.Issue) *Issue {
	if m == nil {
		return nil
	}
	return &Issue{
		ref:                 m,
		Body:                m.Body,
		CommentCount:        int32(m.CommentCount),
		Id:                  m.Id,
		IssueType:           m.IssueType,
		Labels:              m.Labels,
		Repository:          m.Repository,
		RepositoryAvatar:    m.RepositoryAvatar,
		RepositoryUpdatedAt: graphql.Time{Time: m.RepositoryUpdatedAt},
		Title:               m.Title,
		Url:                 m.Url,
	}
}

type CreateIssueInput struct {
	Body                string
	CommentCount        int32
	IssueType           string
	Labels              []*string
	Repository          string
	RepositoryAvatar    string
	RepositoryUpdatedAt graphql.Time
	Title               string
	Url                 string
	StoryID             graphql.ID
	UserID              graphql.ID
}

func (n *CreateIssueInput) ToModel() (*issuesstore.Issue, error) {
	if n == nil {
		return nil, nil
	}
	storyID, err := ParseGraphqlID(n.StoryID)
	if err != nil {
		return nil, err
	}

	return &issuesstore.Issue{
		Body:                n.Body,
		CommentCount:        int64(n.CommentCount),
		IssueType:           n.IssueType,
		Labels:              n.Labels,
		Repository:          n.Repository,
		RepositoryAvatar:    n.RepositoryAvatar,
		RepositoryUpdatedAt: n.RepositoryUpdatedAt.Time,
		Title:               n.Title,
		Url:                 n.Url,
		StoryID:             storyID.ID,
	}, nil
}
func (n *Issue) ID(ctx context.Context) graphql.ID {
	return NewIDMarshaller(NodeTypeIssue, n.Id, false).
		ToGraphqlID()
}

type IssuesConnection struct {
	Edges    []*IssueEdge
	PageInfo *PageInfo
	filters  *issuesstore.IssueFilters
}

func NewIssuesConnection(
	filters *issuesstore.IssueFilters,
	data []*issuesstore.Issue,
	hasNextPage bool,
	hasPreviousPage bool,
	cursors []string,
) *IssuesConnection {
	edges := []*IssueEdge{}
	for i, d := range data {
		node := NewIssue(d)

		edges = append(edges, &IssueEdge{
			Node:   node,
			Cursor: cursors[i],
		})
	}
	var startCursor, endCursor *string
	if len(cursors) > 0 {
		startCursor = &cursors[0]
		endCursor = &cursors[len(cursors)-1]
	}
	return &IssuesConnection{
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

func (c IssuesConnection) TotalCount(ctx context.Context) (int32, error) {
	count, err := storemanager.FromContext(ctx).IssuesStore.CountIssues(ctx, c.filters)
	return int32(count), err
}

type IssueEdge struct {
	Node   *Issue
	Cursor string
}
