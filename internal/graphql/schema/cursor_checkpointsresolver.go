package schema

import (
	"github.com/firstcontributions/backend/internal/models/usersstore"
	graphql "github.com/graph-gophers/graphql-go"
)

type CursorCheckpoints struct {
	PullRequests string
	TimeCreated  graphql.Time
	TimeUpdated  graphql.Time
}

func NewCursorCheckpoints(m *usersstore.CursorCheckpoints) *CursorCheckpoints {
	if m == nil {
		return nil
	}
	return &CursorCheckpoints{
		PullRequests: m.PullRequests,
		TimeCreated:  graphql.Time{Time: m.TimeCreated},
		TimeUpdated:  graphql.Time{Time: m.TimeUpdated},
	}
}
