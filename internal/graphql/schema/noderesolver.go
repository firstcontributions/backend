package schema

import (
	"context"
	"errors"

	"github.com/firstcontributions/backend/internal/storemanager"
	"github.com/graph-gophers/graphql-go"
)

type Node interface {
	ID(context.Context) graphql.ID
}

type NodeResolver struct {
	Node
}

type NodeIDInput struct {
	ID graphql.ID
}

func (r *Resolver) Node(ctx context.Context, in NodeIDInput) (*NodeResolver, error) {
	store := storemanager.FromContext(ctx)
	id, err := ParseGraphqlID(in.ID)
	if err != nil {
		return nil, err
	}
	switch id.Type {
	case "badge":
		badgeData, err := store.UsersStore.GetBadgeByID(ctx, id.ID)
		if err != nil {
			return nil, err
		}
		badgeNode := NewBadge(badgeData)
		return &NodeResolver{
			Node: badgeNode,
		}, nil
	case "user":
		userData, err := store.UsersStore.GetUserByID(ctx, id.ID)
		if err != nil {
			return nil, err
		}
		userNode := NewUser(userData)
		return &NodeResolver{
			Node: userNode,
		}, nil
	}
	return nil, errors.New("invalid ID")
}
func (n *NodeResolver) ToBadge() (*Badge, bool) {
	t, ok := n.Node.(*Badge)
	return t, ok
}
func (n *NodeResolver) ToUser() (*User, bool) {
	t, ok := n.Node.(*User)
	return t, ok
}
