package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/issuesstore"
	graphql "github.com/graph-gophers/graphql-go"
)

type Issue struct {
	Id                string
	IssueType         string
	Repository        string
	RespositoryAvatar string
	Title             string
	Url               string
}

func NewIssue(m *issuesstore.Issue) *Issue {
	if m == nil {
		return nil
	}
	return &Issue{
		Id:                m.Id,
		IssueType:         m.IssueType,
		Repository:        m.Repository,
		RespositoryAvatar: m.RespositoryAvatar,
		Title:             m.Title,
		Url:               m.Url,
	}
}

type CreateIssueInput struct {
	IssueType         string
	Repository        string
	RespositoryAvatar string
	Title             string
	Url               string
}

func (n *CreateIssueInput) ToModel() *issuesstore.Issue {
	if n == nil {
		return nil
	}
	return &issuesstore.Issue{
		IssueType:         n.IssueType,
		Repository:        n.Repository,
		RespositoryAvatar: n.RespositoryAvatar,
		Title:             n.Title,
		Url:               n.Url,
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
