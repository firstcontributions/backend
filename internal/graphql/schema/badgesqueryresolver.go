package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/storemanager"
)

type BadgesInput struct {
	First  *int64
	Last   *int64
	After  *string
	Before *string
}

func (n *User) Badges(ctx context.Context, in *BadgesInput) (*BadgesConnection, error) {
	store := storemanager.FromContext(ctx)
	var first, last *int64
	if in.First != nil {
		tmp := int64(*in.First)
		first = &tmp
	}
	if in.Last != nil {
		tmp := int64(*in.Last)
		last = &tmp
	}
	data, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err := store.UsersStore.GetBadges(
		ctx,
		nil,
		&n.Id,
		in.After,
		in.Before,
		first,
		last,
	)
	if err != nil {
		return nil, err
	}
	return NewBadgesConnection(data, hasNextPage, hasPreviousPage, &firstCursor, &lastCursor), nil
}
