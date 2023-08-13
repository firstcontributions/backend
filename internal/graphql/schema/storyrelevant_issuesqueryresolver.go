package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/issuesstore"
	"github.com/firstcontributions/backend/internal/storemanager"
)

type StoryRelevantIssuesInput struct {
	First     *int32
	Last      *int32
	After     *string
	Before    *string
	SortOrder *string
	SortBy    *string
}

func (n *Story) RelevantIssues(ctx context.Context, in *StoryRelevantIssuesInput) (*IssuesConnection, error) {
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
	issueType := "relevant_issues_story"

	filters := &issuesstore.IssueFilters{
		IssueType: &issueType,
		Story:     n.ref,
	}
	sortByStr := ""
	if in.SortBy != nil {
		sortByStr = *in.SortBy
	}
	data, hasNextPage, hasPreviousPage, cursors, err := store.IssuesStore.GetIssues(
		ctx,
		filters,
		in.After,
		in.Before,
		first,
		last,
		issuesstore.GetIssueSortByFromString(sortByStr),
		in.SortOrder,
	)
	if err != nil {
		return nil, err
	}
	return NewIssuesConnection(filters, data, hasNextPage, hasPreviousPage, cursors), nil
}
