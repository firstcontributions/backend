package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
)

// SessionStore implments session.Store interface
type SessionStore struct {
	client *redis.Client
	ttl    time.Duration
}

// NewSessionStore returns an instance of redis based session store
func NewSessionStore(host, port, password string, sessionTime time.Duration) *SessionStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password, // no password set
		DB:       0,        // use default DB
	})
	rdb.AddHook(redisotel.NewTracingHook())
	return &SessionStore{
		client: rdb,
		ttl:    sessionTime,
	}
}

// Get gets value from the store
func (s *SessionStore) Get(ctx context.Context, key string, value interface{}) error {
	return s.client.Get(ctx, key).Scan(value)
}

// Set sets value in the store
func (s *SessionStore) Set(ctx context.Context, key string, value interface{}) error {
	return s.client.SetNX(ctx, key, value, s.ttl).Err()
}
