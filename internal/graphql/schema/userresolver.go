package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	graphql "github.com/graph-gophers/graphql-go"
)

type User struct {
	CursorCheckpoints *CursorCheckpoints
	Handle            string
	Id                string
	Name              string
	TimeCreated       graphql.Time
	TimeUpdated       graphql.Time
	Token             *Token
}

func NewUser(m *usersstore.User) *User {
	if m == nil {
		return nil
	}
	return &User{
		CursorCheckpoints: NewCursorCheckpoints(m.CursorCheckpoints),
		Handle:            m.Handle,
		Id:                m.Id,
		Name:              m.Name,
		TimeCreated:       graphql.Time{Time: m.TimeCreated},
		TimeUpdated:       graphql.Time{Time: m.TimeUpdated},
		Token:             NewToken(m.Token),
	}
}
func (n *User) ID(ctx context.Context) graphql.ID {
	return NewIDMarshaller("user", n.Id).
		ToGraphqlID()
}
