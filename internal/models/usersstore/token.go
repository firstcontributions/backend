package usersstore

import "time"

type TokenSortBy uint8

const (
	TokenSortByDefault = iota
	TokenSortByTimeCreated
)

type Token struct {
	AccessToken  string    `bson:"access_token,omitempty"`
	Expiry       time.Time `bson:"expiry,omitempty"`
	RefreshToken string    `bson:"refresh_token,omitempty"`
	TokenType    string    `bson:"token_type,omitempty"`
}

func NewToken() *Token {
	return &Token{}
}

type TokenFilters struct {
	Ids []string
}

func (s TokenSortBy) String() string {
	switch s {
	case TokenSortByTimeCreated:
		return "time_created"
	default:
		return "time_created"
	}
}

func GetTokenSortByFromString(s string) TokenSortBy {
	switch s {
	case "time_created":
		return TokenSortByTimeCreated
	default:
		return TokenSortByDefault
	}
}
