package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/firstcontributions/backend/internal/storemanager"
)

type UserBadgesInput struct {
	First     *int32
	Last      *int32
	After     *string
	Before    *string
	SortBy    *string
	SortOrder *string
}

func (n *User) Badges(ctx context.Context, in *UserBadgesInput) (*BadgesConnection, error) {
	var first, last *int64
	if in.First != nil {
		tmp := int64(*in.First)
		first = &tmp
	}
	if in.Last != nil {
		tmp := int64(*in.Last)
		last = &tmp
	}
	store := storemanager.FromContext(ctx)

	filters := &usersstore.BadgeFilters{
		User: n.ref,
	}
	sortByStr := ""
	if in.SortBy != nil {
		sortByStr = *in.SortBy
	}
	data, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err := store.UsersStore.GetBadges(
		ctx,
		filters,
		in.After,
		in.Before,
		first,
		last,
		usersstore.GetBadgeSortByFromString(sortByStr),
		in.SortOrder,
	)
	if err != nil {
		return nil, err
	}
	return NewBadgesConnection(filters, data, hasNextPage, hasPreviousPage, &firstCursor, &lastCursor), nil
}
