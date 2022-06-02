package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/firstcontributions/backend/pkg/cursor"
	graphql "github.com/graph-gophers/graphql-go"
)

type Badge struct {
	ref                           *usersstore.Badge
	CurrentLevel                  int32
	DisplayName                   string
	Id                            string
	Points                        int32
	ProgressPercentageToNextLevel int32
	TimeCreated                   graphql.Time
	TimeUpdated                   graphql.Time
}

func NewBadge(m *usersstore.Badge) *Badge {
	if m == nil {
		return nil
	}
	return &Badge{
		ref:                           m,
		CurrentLevel:                  int32(m.CurrentLevel),
		DisplayName:                   m.DisplayName,
		Id:                            m.Id,
		Points:                        int32(m.Points),
		ProgressPercentageToNextLevel: int32(m.ProgressPercentageToNextLevel),
		TimeCreated:                   graphql.Time{Time: m.TimeCreated},
		TimeUpdated:                   graphql.Time{Time: m.TimeUpdated},
	}
}

type CreateBadgeInput struct {
	CurrentLevel                  int32
	DisplayName                   string
	Points                        int32
	ProgressPercentageToNextLevel int32
}

func (n *CreateBadgeInput) ToModel() *usersstore.Badge {
	if n == nil {
		return nil
	}
	return &usersstore.Badge{
		CurrentLevel:                  int64(n.CurrentLevel),
		DisplayName:                   n.DisplayName,
		Points:                        int64(n.Points),
		ProgressPercentageToNextLevel: int64(n.ProgressPercentageToNextLevel),
	}
}
func (n *Badge) ID(ctx context.Context) graphql.ID {
	return NewIDMarshaller("badge", n.Id).
		ToGraphqlID()
}

type BadgesConnection struct {
	Edges    []*BadgeEdge
	PageInfo *PageInfo
}

func NewBadgesConnection(
	data []*usersstore.Badge,
	hasNextPage bool,
	hasPreviousPage bool,
	firstCursor *string,
	lastCursor *string,
) *BadgesConnection {
	edges := []*BadgeEdge{}
	for _, d := range data {
		node := NewBadge(d)

		edges = append(edges, &BadgeEdge{
			Node:   node,
			Cursor: cursor.NewCursor(d.Id, d.TimeCreated).String(),
		})
	}
	return &BadgesConnection{
		Edges: edges,
		PageInfo: &PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
			StartCursor:     firstCursor,
			EndCursor:       lastCursor,
		},
	}
}

type BadgeEdge struct {
	Node   *Badge
	Cursor string
}
