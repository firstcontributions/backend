package csrf

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Store interface {
	Insert(context.Context, string) error
	Get(context.Context, string) (string, error)
	Invalidate(context.Context, string) error
}

type Manager struct {
	Store
}

func NewManager(store Store) *Manager {
	return &Manager{
		Store: store,
	}
}

func (m *Manager) Generate(ctx context.Context) (string, error) {
	token, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	if err := m.Insert(ctx, token.String()); err != nil {
		return "", err
	}
	return token.String(), nil
}

func (m *Manager) Validate(ctx context.Context, token string) error {
	_, err := m.Get(ctx, token)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("could not find the token, err: %w", err)
		}
		return err
	}
	return m.Invalidate(ctx, token)
}
