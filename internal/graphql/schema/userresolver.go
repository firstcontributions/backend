package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	graphql "github.com/graph-gophers/graphql-go"
)

type User struct {
	ref                  *usersstore.User
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
	GitContributionStats *GitContributionStats
	Handle               string
	Name                 string
	Reputation           *Reputation
}

func (n *CreateUserInput) ToModel() *usersstore.User {
	if n == nil {
		return nil
	}
	return &usersstore.User{
		GitContributionStats: n.GitContributionStats.ToModel(),
		Handle:               n.Handle,
		Name:                 n.Name,
		Reputation:           n.Reputation.ToModel(),
	}
}

type UpdateUserInput struct {
	GitContributionStats *GitContributionStats
	Name                 *string
	Reputation           *Reputation
}

func (n *UpdateUserInput) ToModel() *usersstore.UserUpdate {
	if n == nil {
		return nil
	}
	return &usersstore.UserUpdate{
		GitContributionStats: n.GitContributionStats.ToModel(),
		Name:                 n.Name,
		Reputation:           n.Reputation.ToModel(),
	}
}
func (n *User) ID(ctx context.Context) graphql.ID {
	return NewIDMarshaller("user", n.Id).
		ToGraphqlID()
}
