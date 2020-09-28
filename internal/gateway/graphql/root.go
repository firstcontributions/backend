package graphql

import (
	"context"

	"github.com/firstcontributions/backend/internal/gateway/rpcs"
	"github.com/firstcontributions/backend/internal/gateway/session"
)

type Resolver struct {
	*rpcs.ProfileManager
}

func (r *Resolver) Viewer(ctx context.Context) (*profileResolver, error) {
	meta := ctx.Value(session.CxtKeySession)
	handle := meta.(session.MetaData).Handle

	profile, err := r.GetProfile(ctx, handle)
	if err != nil {
		return nil, err
	}
	return &profileResolver{profile: profile}, nil
}
