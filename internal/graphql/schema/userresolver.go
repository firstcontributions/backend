package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/firstcontributions/backend/internal/storemanager"
	graphql "github.com/graph-gophers/graphql-go"
)

type User struct {
	ref                  *usersstore.User
	Avatar               string
	Bio                  string
	GitContributionStats *GitContributionStats
	Handle               string
	Id                   string
	Name                 string
	Reputation           *Reputation
	TimeCreated          graphql.Time
	TimeUpdated          graphql.Time
}

func NewUser(m *usersstore.User) *User {
	if m == nil {
		return nil
	}
	return &User{
		ref:                  m,
		Avatar:               m.Avatar,
		Bio:                  m.Bio,
		GitContributionStats: NewGitContributionStats(m.GitContributionStats),
		Handle:               m.Handle,
		Id:                   m.Id,
		Name:                 m.Name,
		Reputation:           NewReputation(m.Reputation),
		TimeCreated:          graphql.Time{Time: m.TimeCreated},
		TimeUpdated:          graphql.Time{Time: m.TimeUpdated},
	}
}

type CreateUserInput struct {
	Avatar               string
	Bio                  string
	GitContributionStats *GitContributionStats
	Handle               string
	Name                 string
	Reputation           *Reputation
}

func (n *CreateUserInput) ToModel() (*usersstore.User, error) {
	if n == nil {
		return nil, nil
	}

	return &usersstore.User{
		Avatar:               n.Avatar,
		Bio:                  n.Bio,
		GitContributionStats: n.GitContributionStats.ToModel(),
		Handle:               n.Handle,
		Name:                 n.Name,
		Reputation:           n.Reputation.ToModel(),
	}, nil
}

type UpdateUserInput struct {
	ID                   graphql.ID
	Avatar               *string
	Bio                  *string
	GitContributionStats *GitContributionStats
	Name                 *string
	Reputation           *Reputation
}

func (n *UpdateUserInput) ToModel() *usersstore.UserUpdate {
	if n == nil {
		return nil
	}
	return &usersstore.UserUpdate{
		Avatar:               n.Avatar,
		Bio:                  n.Bio,
		GitContributionStats: n.GitContributionStats.ToModel(),
		Name:                 n.Name,
		Reputation:           n.Reputation.ToModel(),
	}
}
func (n *User) ID(ctx context.Context) graphql.ID {
	return NewIDMarshaller(NodeTypeUser, n.Id, true).
		ToGraphqlID()
}

type UsersConnection struct {
	Edges    []*UserEdge
	PageInfo *PageInfo
	filters  *usersstore.UserFilters
}

func NewUsersConnection(
	filters *usersstore.UserFilters,
	data []*usersstore.User,
	hasNextPage bool,
	hasPreviousPage bool,
	cursors []string,
) *UsersConnection {
	edges := []*UserEdge{}
	for i, d := range data {
		node := NewUser(d)

		edges = append(edges, &UserEdge{
			Node:   node,
			Cursor: cursors[i],
		})
	}
	var startCursor, endCursor *string
	if len(cursors) > 0 {
		startCursor = &cursors[0]
		endCursor = &cursors[len(cursors)-1]
	}
	return &UsersConnection{
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

func (c UsersConnection) TotalCount(ctx context.Context) (int32, error) {
	count, err := storemanager.FromContext(ctx).UsersStore.CountUsers(ctx, c.filters)
	return int32(count), err
}

type UserEdge struct {
	Node   *User
	Cursor string
}
