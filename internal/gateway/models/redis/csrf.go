package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type CSRFStore struct {
	client *redis.Client
	ttl    time.Duration
}

// NewCSRFStore returns an instance of redis based session store
func NewCSRFStore(host, port, password string, sessionTime time.Duration) *CSRFStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password, // no password set
		DB:       0,        // use default DB
	})
	return &CSRFStore{
		client: rdb,
		ttl:    sessionTime,
	}
}

// Get gets token from the store
func (s *CSRFStore) Get(ctx context.Context, token string) (string, error) {
	return s.client.Get(ctx, token).Result()
}

// Insert adds token in the store
func (s *CSRFStore) Insert(ctx context.Context, token string) error {
	return s.client.SetNX(ctx, token, token, s.ttl).Err()
}

//Invalidate will remove token from the store
func (s *CSRFStore) Invalidate(ctx context.Context, token string) error {
	return s.client.Del(ctx, token).Err()
}
