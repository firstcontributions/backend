package issuesstore

type IssueRecomendations struct {
	Cursor string
}

func NewIssueRecomendations() *IssueRecomendations {
	return &IssueRecomendations{}
}
