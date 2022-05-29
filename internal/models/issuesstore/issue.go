package issuesstore

import "time"

type Issue struct {
	Body                string    `bson:"body,omitempty"`
	CommentCount        int64     `bson:"comment_count,omitempty"`
	Id                  string    `bson:"_id"`
	IssueType           string    `bson:"issue_type,omitempty"`
	Labels              []*string  `bson:"labels,omitempty"`
	Repository          string    `bson:"repository,omitempty"`
	RepositoryUpdatedAt time.Time `bson:"repository_updated_at,omitempty"`
	RespositoryAvatar   string    `bson:"respository_avatar,omitempty"`
	Title               string    `bson:"title,omitempty"`
	Url                 string    `bson:"url,omitempty"`
	Cursor              string
}

func NewIssue() *Issue {
	return &Issue{}
}
