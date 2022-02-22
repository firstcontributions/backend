package issuesstore

type Issue struct {
	Id                string `bson:"_id"`
	IssueType         string `bson:"issue_type,omitempty"`
	Repository        string `bson:"repository,omitempty"`
	RespositoryAvatar string `bson:"respository_avatar,omitempty"`
	Title             string `bson:"title,omitempty"`
	Url               string `bson:"url,omitempty"`
	Cursor            string
}

func NewIssue() *Issue {
	return &Issue{}
}
