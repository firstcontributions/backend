package schema

import (
	"github.com/firstcontributions/backend/internal/models/usersstore"
	graphql "github.com/graph-gophers/graphql-go"
)

type Token struct {
	AccessToken  string
	Expiry       graphql.Time
	RefreshToken string
	TimeCreated  graphql.Time
	TimeUpdated  graphql.Time
	TokenType    string
}

func NewToken(m *usersstore.Token) *Token {
	return &Token{
		AccessToken:  m.AccessToken,
		Expiry:       graphql.Time{Time: m.Expiry},
		RefreshToken: m.RefreshToken,
		TimeCreated:  graphql.Time{Time: m.TimeCreated},
		TimeUpdated:  graphql.Time{Time: m.TimeUpdated},
		TokenType:    m.TokenType,
	}
}
