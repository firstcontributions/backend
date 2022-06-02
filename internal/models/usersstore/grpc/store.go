package grpc

import (
	"context"
	"time"

	pool "github.com/processout/grpc-go-pool"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UsersStore struct {
	pool *pool.Pool
}

// NewUsersStore makes and keeps connection pool to given grpc server
// and return an instance of the client
func NewUsersStore(ctx context.Context, url string, initConnections, connectionCapacity, ttl int) (*UsersStore, error) {
	pool, err := pool.New(
		func() (*grpc.ClientConn, error) {
			return grpc.Dial(
				url,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
			)
		},
		initConnections,
		connectionCapacity,
		time.Duration(ttl)*time.Minute,
	)
	if err != nil {
		return nil, err
	}
	return &UsersStore{
		pool: pool,
	}, nil
}
