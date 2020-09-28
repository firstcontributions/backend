package graphql

import (
	"context"

	"github.com/firstcontributions/backend/internal/proto"
	"github.com/graph-gophers/graphql-go"
)

type profileResolver struct {
	profile *proto.Profile
}

type badgeResolver struct {
	badge *proto.Badge
}

type badgeArgs struct {
	Page *Pagination
}

func (p *profileResolver) UUID() graphql.ID {
	return graphql.ID(p.profile.Uuid)
}
func (p *profileResolver) Name() string {
	return p.profile.Name
}
func (p *profileResolver) Handle() string {
	return p.profile.Handle
}
func (p *profileResolver) Email() string {
	return p.profile.Email
}
func (p *profileResolver) Avatar() string {
	return p.profile.Avatar
}

func (p *profileResolver) Reputation() int32 {
	return int32(p.profile.Reputation)
}

func (p *profileResolver) Badges(ctx context.Context, args badgeArgs) *[]*badgeResolver {
	resolvers := []*badgeResolver{}
	for _, badge := range p.profile.Badges {
		resolvers = append(resolvers, &badgeResolver{badge: badge})
	}
	return &resolvers
}

func (b *badgeResolver) UUID() graphql.ID {
	return graphql.ID(b.badge.Uuid)
}
func (b *badgeResolver) Name() string {
	return b.badge.Name
}

func (b *badgeResolver) AssignedOn() graphql.Time {
	return graphql.Time{
		Time: b.badge.AssignedOn.AsTime(),
	}
}

func (b *badgeResolver) Progress() int32 {
	return int32(b.badge.Progress)
}
