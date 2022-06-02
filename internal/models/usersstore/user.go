package usersstore

import (
	"time"

	"github.com/firstcontributions/backend/internal/grpc/users/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	CursorCheckpoints *CursorCheckpoints `bson:"cursor_checkpoints,omitempty"`
	Handle            string             `bson:"handle,omitempty"`
	Id                string             `bson:"_id"`
	Name              string             `bson:"name,omitempty"`
	Tags              *Tags              `bson:"tags,omitempty"`
	TimeCreated       time.Time          `bson:"time_created,omitempty"`
	TimeUpdated       time.Time          `bson:"time_updated,omitempty"`
	Token             *Token             `bson:"token,omitempty"`
}

func NewUser() *User {
	return &User{}
}
func (user *User) ToProto() *proto.User {
	return &proto.User{
		CursorCheckpoints: user.CursorCheckpoints.ToProto(),
		Handle:            user.Handle,
		Id:                user.Id,
		Name:              user.Name,
		Tags:              user.Tags.ToProto(),
		TimeCreated:       timestamppb.New(user.TimeCreated),
		TimeUpdated:       timestamppb.New(user.TimeUpdated),
		Token:             user.Token.ToProto(),
	}
}

func (user *User) FromProto(protoUser *proto.User) *User {
	user.CursorCheckpoints = NewCursorCheckpoints().FromProto(protoUser.CursorCheckpoints)
	user.Handle = protoUser.Handle
	user.Id = protoUser.Id
	user.Name = protoUser.Name
	user.Tags = NewTags().FromProto(protoUser.Tags)
	user.TimeCreated = protoUser.TimeCreated.AsTime()
	user.TimeUpdated = protoUser.TimeUpdated.AsTime()
	user.Token = NewToken().FromProto(protoUser.Token)
	return user
}

type UserUpdate struct {
	CursorCheckpoints *CursorCheckpoints `bson:"cursor_checkpoints,omitempty"`
	Name              *string            `bson:"name,omitempty"`
	Tags              *Tags              `bson:"tags,omitempty"`
	TimeUpdated       *time.Time         `bson:"time_updated,omitempty"`
	Token             *Token             `bson:"token,omitempty"`
}

func (user *UserUpdate) ToProto() *proto.UpdateUserRequest {
	p := &proto.UpdateUserRequest{}
	if user.CursorCheckpoints != nil {
		p.CursorCheckpoints = user.CursorCheckpoints.ToProto()
	}
	if user.Name != nil {
		p.Name = user.Name
	}
	if user.Tags != nil {
		p.Tags = user.Tags.ToProto()
	}
	if user.TimeUpdated != nil {
		p.TimeUpdated = timestamppb.New(*user.TimeUpdated)
	}
	if user.Token != nil {
		p.Token = user.Token.ToProto()
	}
	return p
}

func (user *UserUpdate) FromProto(protoUser *proto.UpdateUserRequest) {
	user.CursorCheckpoints = NewCursorCheckpoints().FromProto(protoUser.CursorCheckpoints)
	user.Name = protoUser.Name
	user.Tags = NewTags().FromProto(protoUser.Tags)
	time_updated := protoUser.TimeUpdated.AsTime()
	user.TimeUpdated = &time_updated
	user.Token = NewToken().FromProto(protoUser.Token)
}
