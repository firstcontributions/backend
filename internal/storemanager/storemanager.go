package storemanager

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/issuesstore"
	"github.com/firstcontributions/backend/internal/models/storiesstore"
	"github.com/firstcontributions/backend/internal/models/usersstore"
)

type storeCtxKey int

const (
	store storeCtxKey = iota
)

type Store struct {
	IssuesStore  issuesstore.Store
	StoriesStore storiesstore.Store
	UsersStore   usersstore.Store
}

func NewStore(
	issuesStore issuesstore.Store,
	storiesStore storiesstore.Store,
	usersStore usersstore.Store,
) *Store {
	return &Store{
		IssuesStore:  issuesStore,
		StoriesStore: storiesStore,
		UsersStore:   usersStore,
	}
}

func ContextWithStore(ctx context.Context, s *Store) context.Context {
	return context.WithValue(ctx, store, s)
}

func FromContext(ctx context.Context) *Store {
	return ctx.Value(store).(*Store)
}
