package schema

import "context"

func (r *Resolver) IssuesFeed(ctx context.Context) (*IssuesFeed, error) {
	return NewIssuesFeed(), nil
}
