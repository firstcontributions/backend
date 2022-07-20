package usersstore

import "github.com/firstcontributions/backend/pkg/cursor"

type CursorCheckpointsSortBy uint8

const (
	CursorCheckpointsSortByDefault CursorCheckpointsSortBy = iota
)

type CursorCheckpoints struct {
	PullRequests string `bson:"pull_requests,omitempty"`
}

func NewCursorCheckpoints() *CursorCheckpoints {
	return &CursorCheckpoints{}
}
func (cursor_checkpoints *CursorCheckpoints) Get(field string) interface{} {
	switch field {
	case "pull_requests":
		return cursor_checkpoints.PullRequests
	default:
		return nil
	}
}

type CursorCheckpointsFilters struct {
	Ids []string
}

func (s CursorCheckpointsSortBy) String() string {
	switch s {
	default:
		return "time_created"
	}
}

func GetCursorCheckpointsSortByFromString(s string) CursorCheckpointsSortBy {
	switch s {
	default:
		return CursorCheckpointsSortByDefault
	}
}

func (s CursorCheckpointsSortBy) CursorType() cursor.ValueType {
	switch s {
	default:
		return cursor.ValueTypeTime
	}
}
