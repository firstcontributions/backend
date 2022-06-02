package grpc

import (
	"context"

	"github.com/firstcontributions/backend/internal/grpc/users/proto"
	"github.com/firstcontributions/backend/internal/models/usersstore"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *UsersStore) CreateBadge(ctx context.Context, badge *usersstore.Badge) (*usersstore.Badge, error) {

	request := badge.ToProto()
	request.TimeCreated = timestamppb.Now()
	request.TimeUpdated = timestamppb.Now()

	conn, err := s.pool.Get(ctx)
	if err != nil {
		return nil, err
	}

	response, err := proto.NewUsersServiceClient(conn).CreateBadge(ctx, request)
	if err != nil {
		return nil, err
	}
	return usersstore.NewBadge().FromProto(response), nil
}

func (s *UsersStore) GetBadgeByID(ctx context.Context, id string) (*usersstore.Badge, error) {
	conn, err := s.pool.Get(ctx)
	if err != nil {
		return nil, err
	}

	response, err := proto.NewUsersServiceClient(conn).GetBadgeByID(ctx, &proto.RefByIDRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return usersstore.NewBadge().FromProto(response), nil
}
func getBadgesRequest(
	ctx context.Context,
	ids []string,
	user *usersstore.User,
	after *string,
	before *string,
	first *int64,
	last *int64,
) *proto.GetBadgesRequest {
	request := &proto.GetBadgesRequest{
		Ids: ids,
	}
	if user != nil {
		request.UserId = &user.Id
	}
	if after != nil {
		request.After = after
	}
	if before != nil {
		request.After = before
	}
	if first != nil {
		request.First = first
	}
	if last != nil {
		request.Last = last
	}
	return request
}
func (s *UsersStore) GetBadges(
	ctx context.Context,
	ids []string,
	user *usersstore.User,
	after *string,
	before *string,
	first *int64,
	last *int64,
) (
	[]*usersstore.Badge,
	bool,
	bool,
	string,
	string,
	error,
) {
	request := getBadgesRequest(
		ctx,
		ids,
		user,
		after,
		before,
		first,
		last,
	)
	conn, err := s.pool.Get(ctx)
	if err != nil {
		return nil, false, false, "", "", err
	}
	response, err := proto.NewUsersServiceClient(conn).GetBadges(ctx, request)
	if err != nil {
		return nil, false, false, "", "", err
	}

	badges := []*usersstore.Badge{}
	for _, d := range response.Data {
		badges = append(badges, usersstore.NewBadge().FromProto(d))
	}

	return badges, response.HasNext, response.HasPrevious, response.FirstCursor, response.LastCursor, nil
}
func (s *UsersStore) UpdateBadge(ctx context.Context, id string, badgeUpdate *usersstore.BadgeUpdate) error {
	conn, err := s.pool.Get(ctx)
	if err != nil {
		return err
	}
	input := badgeUpdate.ToProto()
	input.Id = id

	_, err = proto.NewUsersServiceClient(conn).UpdateBadge(ctx, input)
	return err
}
func (s *UsersStore) DeleteBadgeByID(ctx context.Context, id string) error {
	conn, err := s.pool.Get(ctx)
	if err != nil {
		return err
	}

	_, err = proto.NewUsersServiceClient(conn).DeleteBadge(ctx, &proto.RefByIDRequest{Id: id})
	return err
}
