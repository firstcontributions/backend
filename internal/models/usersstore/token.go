package usersstore

import (
	"time"

	"github.com/firstcontributions/backend/internal/grpc/users/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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
func (token *Token) ToProto() *proto.Token {
	return &proto.Token{
		AccessToken:  token.AccessToken,
		Expiry:       timestamppb.New(token.Expiry),
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
	}
}

func (token *Token) FromProto(protoToken *proto.Token) *Token {
	token.AccessToken = protoToken.AccessToken
	token.Expiry = protoToken.Expiry.AsTime()
	token.RefreshToken = protoToken.RefreshToken
	token.TokenType = protoToken.TokenType
	return token
}
