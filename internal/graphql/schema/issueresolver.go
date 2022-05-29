package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/issuesstore"
	graphql "github.com/graph-gophers/graphql-go"
)

type Issue struct {
	Body                string
	CommentCount        int32
	Id                  string
	IssueType           string
	Labels              []*string
	Repository          string
	RepositoryUpdatedAt graphql.Time
	RespositoryAvatar   string
	Title               string
	Url                 string
}

func NewIssue(m *issuesstore.Issue) *Issue {
	if m == nil {
		return nil
	}
	return &Issue{
		Body:                m.Body,
		CommentCount:        int32(m.CommentCount),
		Id:                  m.Id,
		IssueType:           m.IssueType,
		Labels:              m.Labels,
		Repository:          m.Repository,
		RepositoryUpdatedAt: graphql.Time{Time: m.RepositoryUpdatedAt},
		RespositoryAvatar:   m.RespositoryAvatar,
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
	RepositoryUpdatedAt graphql.Time
	RespositoryAvatar   string
	Title               string
	Url                 string
}

func (n *CreateIssueInput) ToModel() *issuesstore.Issue {
	if n == nil {
		return nil
	}
	return &issuesstore.Issue{
		Body:                n.Body,
		CommentCount:        int64(n.CommentCount),
		IssueType:           n.IssueType,
		Labels:              n.Labels,
		Repository:          n.Repository,
		RepositoryUpdatedAt: n.RepositoryUpdatedAt.Time,
		RespositoryAvatar:   n.RespositoryAvatar,
		Title:               n.Title,
		Url:                 n.Url,
	}
}
func (n *Issue) ID(ctx context.Context) graphql.ID {
	return NewIDMarshaller("issue", n.Id).
		ToGraphqlID()
}

type IssuesConnection struct {
	Edges    []*IssueEdge
	PageInfo *PageInfo
}

func NewIssuesConnection(
	data []*issuesstore.Issue,
	hasNextPage bool,
	hasPreviousPage bool,
	firstCursor *string,
	lastCursor *string,
) *IssuesConnection {
	edges := []*IssueEdge{}
	for _, d := range data {
		node := NewIssue(d)

		edges = append(edges, &IssueEdge{
			Node:   node,
			Cursor: d.Cursor,
		})
	}
	return &IssuesConnection{
		Edges: edges,
		PageInfo: &PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
			StartCursor:     firstCursor,
			EndCursor:       lastCursor,
		},
	}
}

type IssueEdge struct {
	Node   *Issue
	Cursor string
}
