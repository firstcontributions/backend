package usersstore

type CursorCheckpointsSortBy uint8

const (
	CursorCheckpointsSortByDefault = iota
	CursorCheckpointsSortByTimeCreated
)

type CursorCheckpoints struct {
	PullRequests string `bson:"pull_requests,omitempty"`
}

func NewCursorCheckpoints() *CursorCheckpoints {
	return &CursorCheckpoints{}
}

type CursorCheckpointsFilters struct {
	Ids []string
}

func (s CursorCheckpointsSortBy) String() string {
	switch s {
	case CursorCheckpointsSortByTimeCreated:
		return "time_created"
	default:
		return "time_created"
	}
}

func GetCursorCheckpointsSortByFromString(s string) CursorCheckpointsSortBy {
	switch s {
	case "time_created":
		return CursorCheckpointsSortByTimeCreated
	default:
		return CursorCheckpointsSortByDefault
	}
}
