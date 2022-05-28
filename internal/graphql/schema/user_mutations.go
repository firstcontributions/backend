package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/storemanager"
)

func (m *Resolver) CreateUser(
	ctx context.Context,
	args struct {
		User *CreateUserInput
	},
) (*User, error) {
	store := storemanager.FromContext(ctx)
	user, err := store.UsersStore.CreateUser(ctx, args.User.ToModel())
	if err != nil {
		return nil, err
	}
	return NewUser(user), nil
}
