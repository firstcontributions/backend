package issuesstore

import "context"

type Store interface {
	// issue methods
	GetIssueByID(ctx context.Context, id string) (*Issue, error)
	GetIssues(ctx context.Context, ids []string, issue_type *string, after *string, before *string, first *int64, last *int64) ([]*Issue, bool, bool, string, string, error)
}
