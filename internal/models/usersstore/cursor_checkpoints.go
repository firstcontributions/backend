package usersstore

import "github.com/firstcontributions/backend/internal/grpc/users/proto"

type CursorCheckpoints struct {
	PullRequests string `bson:"pull_requests,omitempty"`
}

func NewCursorCheckpoints() *CursorCheckpoints {
	return &CursorCheckpoints{}
}
func (cursor_checkpoints *CursorCheckpoints) ToProto() *proto.CursorCheckpoints {
	return &proto.CursorCheckpoints{
		PullRequests: cursor_checkpoints.PullRequests,
	}
}

func (cursor_checkpoints *CursorCheckpoints) FromProto(protoCursorCheckpoints *proto.CursorCheckpoints) *CursorCheckpoints {
	cursor_checkpoints.PullRequests = protoCursorCheckpoints.PullRequests
	return cursor_checkpoints
}
