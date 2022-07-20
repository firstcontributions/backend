package schema

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"

	"github.com/firstcontributions/backend/internal/gateway/session"
	"github.com/firstcontributions/backend/internal/storemanager"
	"github.com/google/uuid"
	graphql "github.com/graph-gophers/graphql-go"
)

type Resolver struct {
}

func (r *Resolver) Viewer(ctx context.Context) (*User, error) {
	session := session.FromContext(ctx)
	if session == nil {
		return nil, errors.New("Unauthorized")
	}
	store := storemanager.FromContext(ctx)

	data, err := store.UsersStore.GetUserByID(ctx, session.UserID())
	if err != nil {
		return nil, err
	}
	return NewUser(data), nil
}

type IDMarshaller struct {
	Type   NodeType
	ID     string
	IsUUID bool
}

func NewIDMarshaller(t NodeType, id string, isUUID bool) *IDMarshaller {
	return &IDMarshaller{
		Type:   t,
		ID:     id,
		IsUUID: isUUID,
	}
}

type PageInfo struct {
	HasNextPage     bool
	HasPreviousPage bool
	StartCursor     *string
	EndCursor       *string
}

func ParseGraphqlID(gid graphql.ID) (*IDMarshaller, error) {
	sDec, err := base64.StdEncoding.DecodeString(string(gid))
	if err != nil {
		return nil, errors.New("invalid ID")
	}
	parts := bytes.Split(sDec, []byte{'|'})
	if len(parts) < 3 {
		return nil, errors.New("invalid ID")
	}
	id := IDMarshaller{}
	id.Type = NodeType(parts[0][0])
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

func (id *IDMarshaller) String() string {
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

func (id *IDMarshaller) ToGraphqlID() graphql.ID {
	return graphql.ID(id.String())
}
