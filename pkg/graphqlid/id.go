package graphqlid

import (
	"bytes"
	"encoding/base64"
	"errors"

	"github.com/google/uuid"
	graphql "github.com/graph-gophers/graphql-go"
)

type GraphqlID struct {
	Type   uint8
	ID     string
	IsUUID bool
}

func NewGraphqlID(t uint8, id string, isUUID bool) *GraphqlID {
	return &GraphqlID{
		Type:   t,
		ID:     id,
		IsUUID: isUUID,
	}
}

func ParseGraphqlID(gid graphql.ID) (*GraphqlID, error) {
	sDec, err := base64.StdEncoding.DecodeString(string(gid))
	if err != nil {
		return nil, errors.New("invalid ID")
	}
	parts := bytes.Split(sDec, []byte{'|'})
	if len(parts) < 3 {
		return nil, errors.New("invalid ID")
	}
	id := GraphqlID{}
	id.Type = uint8(parts[0][0])
	if parts[1][0] == 1 {
		id.IsUUID = true
	}
	if id.IsUUID {
		uid, _ := uuid.FromBytes(parts[2])
		id.ID = uid.String()
	} else {
		id.ID = string(parts[2])
	}
	return &id, nil
}

func (id *GraphqlID) String() string {
	bts := []byte{byte(id.Type), '|'}
	if id.IsUUID {
		bts = append(bts, 1, '|')
		uid, _ := uuid.Parse(id.ID)
		binuuid, _ := uid.MarshalBinary()
		bts = append(bts, binuuid...)
	} else {
		bts = append(bts, 0, '|')
		bts = append(bts, []byte(id.ID)...)
	}
	return base64.StdEncoding.EncodeToString(bts)
}

func (id *GraphqlID) ToGraphqlID() graphql.ID {
	return graphql.ID(id.String())
}
