package issuesstore

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/usersstore"
)

type Store interface {
	// issue methods
	GetIssueByID(ctx context.Context, id string) (*Issue, error)
	GetIssues(ctx context.Context, ids []string, issue_type *string,
		user *usersstore.User, after *string, before *string, first *int64, last *int64) ([]*Issue, bool, bool, string, string, error)
}