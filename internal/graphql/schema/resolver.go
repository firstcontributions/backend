package schema

import (
	"context"
	"errors"

	"github.com/firstcontributions/backend/internal/gateway/session"
	"github.com/firstcontributions/backend/internal/storemanager"
	"github.com/firstcontributions/backend/pkg/graphqlid"
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
	*graphqlid.GraphqlID
}

func NewIDMarshaller(t NodeType, id string, isUUID bool) *IDMarshaller {
	return &IDMarshaller{
		GraphqlID: graphqlid.NewGraphqlID(uint8(t), id, isUUID),
	}
}

func (id *IDMarshaller) NodeType() NodeType {
	return NodeType(id.Type)
}

func ParseGraphqlID(gid graphql.ID) (*IDMarshaller, error) {
	id, err := graphqlid.ParseGraphqlID(gid)
	if err != nil {
		return nil, err
	}
	return &IDMarshaller{GraphqlID: id}, nil
}

type PageInfo struct {
	HasNextPage     bool
	HasPreviousPage bool
	StartCursor     *string
	EndCursor       *string
}
