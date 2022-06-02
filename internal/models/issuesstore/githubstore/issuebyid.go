package githubstore

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/issuesstore"
)

type GetIssueByIDQuery struct {
	Node struct {
		Issue GitIssue `graphql:"... on Issue"`
	} `graphql:"node(id:$id)"`
}

func (i *GitHubStore) GetIssueByID(ctx context.Context, id string) (*issuesstore.Issue, error) {
	query := &GetIssueByIDQuery{}
	if err := i.Query(ctx, query, map[string]interface{}{
		"id": id,
	}); err != nil {
		return nil, err
	}
	return issueFromGithubIssue(query.Node.Issue), nil
}
