package service

import (
	"context"

	"github.com/firstcontributions/backend/internal/grpc/users/proto"
	"github.com/firstcontributions/backend/internal/models/usersstore"
)

func (s *Service) CreateUser(ctx context.Context, in *proto.User) (*proto.User, error) {

	user := usersstore.NewUser()
	user.FromProto(in)
	res, err := s.Store.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return res.ToProto(), nil
}

func (s *Service) GetUserByID(ctx context.Context, in *proto.RefByIDRequest) (*proto.User, error) {
	user, err := s.Store.GetUserByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return user.ToProto(), nil
}

func (s *Service) GetUsers(ctx context.Context, in *proto.GetUsersRequest) (*proto.GetUsersResponse, error) {
	data, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err := s.Store.GetUsers(
		ctx,
		in.Ids,
		in.Search, in.Handle,
		in.After,
		in.Before,
		in.First,
		in.Last,
	)
	if err != nil {
		return nil, err
	}
	res := []*proto.User{}
	for _, d := range data {
		res = append(res, d.ToProto())
	}
	return &proto.GetUsersResponse{
		HasNext:     hasNextPage,
		HasPrevious: hasPreviousPage,
		FirstCursor: firstCursor,
		LastCursor:  lastCursor,
		Data:        res,
	}, nil
}
func (s *Service) UpdateUser(ctx context.Context, in *proto.UpdateUserRequest) (*proto.StatusResponse, error) {
	updateUser := &usersstore.UserUpdate{}
	updateUser.FromProto(in)

	err := s.Store.UpdateUser(ctx, in.Id, updateUser)
	if err != nil {
		return nil, err
	}
	return &proto.StatusResponse{Status: true}, nil
}
func (s *Service) DeleteUser(ctx context.Context, in *proto.RefByIDRequest) (*proto.StatusResponse, error) {

	err := s.Store.DeleteUserByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &proto.StatusResponse{Status: true}, nil
}
