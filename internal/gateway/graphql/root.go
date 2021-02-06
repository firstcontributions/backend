package graphql

import (
	"context"
	"fmt"

	"github.com/firstcontributions/backend/internal/gateway/rpcs"
	"github.com/firstcontributions/backend/internal/gateway/session"
	"google.golang.org/grpc/metadata"
)

type Resolver struct {
	*rpcs.ProfileManager
}

func (r *Resolver) Viewer(ctx context.Context) (*profileResolver, error) {
	meta := ctx.Value(session.CxtKeySession).(session.MetaData)
	fmt.Println("meta", meta.Proto())
	ctx = metadata.NewOutgoingContext(ctx, meta.Proto())
	profile, err := r.GetProfile(ctx, meta.Handle())
	if err != nil {
		return nil, err
	}
	return &profileResolver{profile: profile}, nil
}
