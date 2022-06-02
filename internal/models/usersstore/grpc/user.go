package grpc

import (
	"context"

	"github.com/firstcontributions/backend/internal/grpc/users/proto"
	"github.com/firstcontributions/backend/internal/models/usersstore"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *UsersStore) CreateUser(ctx context.Context, user *usersstore.User) (*usersstore.User, error) {

	request := user.ToProto()
	request.TimeCreated = timestamppb.Now()
	request.TimeUpdated = timestamppb.Now()

	conn, err := s.pool.Get(ctx)
	if err != nil {
		return nil, err
	}

	response, err := proto.NewUsersServiceClient(conn).CreateUser(ctx, request)
	if err != nil {
		return nil, err
	}
	return usersstore.NewUser().FromProto(response), nil
}

func (s *UsersStore) GetUserByID(ctx context.Context, id string) (*usersstore.User, error) {
	conn, err := s.pool.Get(ctx)
	if err != nil {
		return nil, err
	}

	response, err := proto.NewUsersServiceClient(conn).GetUserByID(ctx, &proto.RefByIDRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return usersstore.NewUser().FromProto(response), nil
}
func getUsersRequest(
	ctx context.Context,
	ids []string,
	search *string,
	handle *string,
	after *string,
	before *string,
	first *int64,
	last *int64,
) *proto.GetUsersRequest {
	request := &proto.GetUsersRequest{
		Ids: ids,
	}
	if search != nil {
		request.Search = search
	}
	if handle != nil {
		request.Handle = handle
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
func (s *UsersStore) GetUsers(
	ctx context.Context,
	ids []string,
	search *string,
	handle *string,
	after *string,
	before *string,
	first *int64,
	last *int64,
) (
	[]*usersstore.User,
	bool,
	bool,
	string,
	string,
	error,
) {
	request := getUsersRequest(
		ctx,
		ids,
		search,
		handle,
		after,
		before,
		first,
		last,
	)
	conn, err := s.pool.Get(ctx)
	if err != nil {
		return nil, false, false, "", "", err
	}
	response, err := proto.NewUsersServiceClient(conn).GetUsers(ctx, request)
	if err != nil {
		return nil, false, false, "", "", err
	}

	users := []*usersstore.User{}
	for _, d := range response.Data {
		users = append(users, usersstore.NewUser().FromProto(d))
	}

	return users, response.HasNext, response.HasPrevious, response.FirstCursor, response.LastCursor, nil
}
func (s *UsersStore) UpdateUser(ctx context.Context, id string, userUpdate *usersstore.UserUpdate) error {
	conn, err := s.pool.Get(ctx)
	if err != nil {
		return err
	}
	input := userUpdate.ToProto()
	input.Id = id

	_, err = proto.NewUsersServiceClient(conn).UpdateUser(ctx, input)
	return err
}
func (s *UsersStore) DeleteUserByID(ctx context.Context, id string) error {
	conn, err := s.pool.Get(ctx)
	if err != nil {
		return err
	}

	_, err = proto.NewUsersServiceClient(conn).DeleteUser(ctx, &proto.RefByIDRequest{Id: id})
	return err
}
