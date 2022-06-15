package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/storemanager"
)

type StoryReactionsInputInput struct {
	First  *int32
	Last   *int32
	After  *string
	Before *string
}

func (n *Story) Reactions(ctx context.Context, in *StoryReactionsInputInput) (*ReactionsConnection, error) {
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
	data, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err := store.StoriesStore.GetReactions(
		ctx,
		nil,
		nil,
		n.ref,
		in.After,
		in.Before,
		first,
		last,
	)
	if err != nil {
		return nil, err
	}
	return NewReactionsConnection(data, hasNextPage, hasPreviousPage, &firstCursor, &lastCursor), nil
}
