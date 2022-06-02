package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	graphql "github.com/graph-gophers/graphql-go"
)

type User struct {
	ref         *usersstore.User
	Handle      string
	Id          string
	Name        string
	TimeCreated graphql.Time
	TimeUpdated graphql.Time
}

func NewUser(m *usersstore.User) *User {
	if m == nil {
		return nil
	}
	return &User{
		ref:         m,
		Handle:      m.Handle,
		Id:          m.Id,
		Name:        m.Name,
		TimeCreated: graphql.Time{Time: m.TimeCreated},
		TimeUpdated: graphql.Time{Time: m.TimeUpdated},
	}
}

type CreateUserInput struct {
	Handle string
	Name   string
}

func (n *CreateUserInput) ToModel() *usersstore.User {
	if n == nil {
		return nil
	}
	return &usersstore.User{
		Handle: n.Handle,
		Name:   n.Name,
	}
}

type UpdateUserInput struct {
	Name *string
}

func (n *UpdateUserInput) ToModel() *usersstore.UserUpdate {
	if n == nil {
		return nil
	}
	return &usersstore.UserUpdate{
		Name: n.Name,
	}
}
func (n *User) ID(ctx context.Context) graphql.ID {
	return NewIDMarshaller("user", n.Id).
		ToGraphqlID()
}
