package usersstore

type CursorCheckpoints struct {
	PullRequests string `bson:"pull_requests,omitempty"`
}

func NewCursorCheckpoints() *CursorCheckpoints {
	return &CursorCheckpoints{}
}
