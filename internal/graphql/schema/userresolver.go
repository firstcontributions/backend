package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/usersstore"
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
	return NewIDMarshaller("user", n.Id).
		ToGraphqlID()
}
