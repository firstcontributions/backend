package session

import "context"

// Manager implements a session manager
type Manager struct {
	Store
}

// Store is an interface of functions for managing a session store
type Store interface {
	Set(context.Context, string, interface{}) error
	Get(context.Context, string, interface{}) error
}

func NewManager(store Store) *Manager {
	return &Manager{
		Store: store,
	}
}
