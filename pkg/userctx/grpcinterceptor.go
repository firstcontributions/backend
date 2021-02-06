package userctx

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Authorize unary interceptor function to handle authorize per RPC call
func Authorize(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (
	interface{}, error,
) {
	m := FromIncomingCtx(ctx)
	if m == nil || m.UserID() == "" {
		return nil, grpc.Errorf(codes.Unauthenticated, "no user context provided with the request")
	}
	log.Printf("from the interceptor")
	// Calls the handler
	h, err := handler(ctx, req)
	return h, err
}
