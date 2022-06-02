package service

import (
	"context"

	"github.com/firstcontributions/backend/internal/grpc/users/proto"
	"github.com/firstcontributions/backend/internal/models/usersstore"
)

func (s *Service) CreateBadge(ctx context.Context, in *proto.Badge) (*proto.Badge, error) {

	badge := usersstore.NewBadge()
	badge.FromProto(in)
	res, err := s.Store.CreateBadge(ctx, badge)
	if err != nil {
		return nil, err
	}
	return res.ToProto(), nil
}

func (s *Service) GetBadgeByID(ctx context.Context, in *proto.RefByIDRequest) (*proto.Badge, error) {
	badge, err := s.Store.GetBadgeByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return badge.ToProto(), nil
}

func (s *Service) GetBadges(ctx context.Context, in *proto.GetBadgesRequest) (*proto.GetBadgesResponse, error) {
	var user *usersstore.User
	if in.UserId != nil {
		user = &usersstore.User{Id: *in.UserId}
	}
	data, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err := s.Store.GetBadges(
		ctx,
		in.Ids,
		user,
		in.After,
		in.Before,
		in.First,
		in.Last,
	)
	if err != nil {
		return nil, err
	}
	res := []*proto.Badge{}
	for _, d := range data {
		res = append(res, d.ToProto())
	}
	return &proto.GetBadgesResponse{
		HasNext:     hasNextPage,
		HasPrevious: hasPreviousPage,
		FirstCursor: firstCursor,
		LastCursor:  lastCursor,
		Data:        res,
	}, nil
}
func (s *Service) UpdateBadge(ctx context.Context, in *proto.UpdateBadgeRequest) (*proto.StatusResponse, error) {
	updateBadge := &usersstore.BadgeUpdate{}
	updateBadge.FromProto(in)

	err := s.Store.UpdateBadge(ctx, in.Id, updateBadge)
	if err != nil {
		return nil, err
	}
	return &proto.StatusResponse{Status: true}, nil
}
func (s *Service) DeleteBadge(ctx context.Context, in *proto.RefByIDRequest) (*proto.StatusResponse, error) {

	err := s.Store.DeleteBadgeByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &proto.StatusResponse{Status: true}, nil
}
