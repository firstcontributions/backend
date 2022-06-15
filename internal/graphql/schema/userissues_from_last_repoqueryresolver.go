package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/storemanager"
)

type UserIssuesFromLastRepoInputInput struct {
	First  *int32
	Last   *int32
	After  *string
	Before *string
}

func (n *User) IssuesFromLastRepo(ctx context.Context, in *UserIssuesFromLastRepoInputInput) (*IssuesConnection, error) {
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
	issueType := "last_repo_issues"
	data, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err := store.IssuesStore.GetIssues(
		ctx,
		nil,
		&issueType,
		n.ref,
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