package githubstore

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/issuesstore"
)

func (i *GitHubStore) GetIssueByID(ctx context.Context, id string) (*issuesstore.Issue, error) {
	return nil, nil
}
