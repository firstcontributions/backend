package usersstore

import (
	"time"

	"github.com/firstcontributions/backend/pkg/authorizer"
	"github.com/firstcontributions/backend/pkg/cursor"
)

type TokenSortBy uint8

const (
	TokenSortByDefault TokenSortBy = iota
)

type Token struct {
	AccessToken  string            `bson:"access_token,omitempty"`
	Expiry       time.Time         `bson:"expiry,omitempty"`
	RefreshToken string            `bson:"refresh_token,omitempty"`
	TokenType    string            `bson:"token_type,omitempty"`
	Ownership    *authorizer.Scope `bson:"ownership,omitempty"`
}

func NewToken() *Token {
	return &Token{}
}
func (token *Token) Get(field string) interface{} {
	switch field {
	case "access_token":
		return token.AccessToken
	case "expiry":
		return token.Expiry
	case "refresh_token":
		return token.RefreshToken
	case "token_type":
		return token.TokenType
	default:
		return nil
	}
}

type TokenFilters struct {
	Ids []string
}

func (s TokenSortBy) String() string {
	switch s {
	default:
		return "time_created"
	}
}

func GetTokenSortByFromString(s string) TokenSortBy {
	switch s {
	default:
		return TokenSortByDefault
	}
}

func (s TokenSortBy) CursorType() cursor.ValueType {
	switch s {
	default:
		return cursor.ValueTypeTime
	}
}
