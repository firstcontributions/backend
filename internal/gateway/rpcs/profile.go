package rpcs

import (
	"context"

	"github.com/firstcontributions/backend/internal/proto"
	pool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
)

type ProfileManager struct {
	ConnectionFactory
}

type ConnectionFactory interface {
	Get(ctx context.Context) (*pool.ClientConn, error)
}

func NewProfileManager(c ConnectionFactory) *ProfileManager {
	return &ProfileManager{
		ConnectionFactory: c,
	}
}

func (p *ProfileManager) GetProfile(ctx context.Context, handle string) (*proto.Profile, error) {
	c, err := p.Get(ctx)
	if err != nil {
		return nil, err
	}
	return proto.NewProfileServiceClient(c).GetProfile(ctx, &proto.GetProfileRequest{Handle: handle})
}
func (p *ProfileManager) CreateProfile(ctx context.Context, in *proto.Profile, opts ...grpc.CallOption) (*proto.Profile, error) {
	c, err := p.Get(ctx)
	if err != nil {
		return nil, err
	}
	return proto.NewProfileServiceClient(c).CreateProfile(ctx, in, opts...)
}
