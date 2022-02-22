package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/storemanager"
)

type IssuesFromOtherRecentReposInput struct {
	First  *int32
	Last   *int32
	After  *string
	Before *string
}

func (n *IssuesFeed) IssuesFromOtherRecentRepos(ctx context.Context, in *IssuesFromOtherRecentReposInput) (*IssuesConnection, error) {
	var first, last *int64
	if in.First != nil {
		tmp := int64(*in.First)
		first = &tmp
	}
	if in.Last != nil {
		tmp := int64(*in.Last)
		last = &tmp
	}
	store := storemanager.FromContext(ctx)
	issueType := "issues_from_other_recent_repos"
	data, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err := store.IssuesStore.GetIssues(
		ctx,
		nil,
		&issueType,
		in.After,
		in.Before,
		first,
		last,
	)
	if err != nil {
		return nil, err
	}
	return NewIssuesConnection(data, hasNextPage, hasPreviousPage, &firstCursor, &lastCursor), nil
}
