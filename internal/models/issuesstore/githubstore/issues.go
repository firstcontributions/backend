package githubstore

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/issuesstore"
	"github.com/shurcooL/githubv4"
)

type IssueQuery struct {
	Search struct {
		Edges []struct {
			Node struct {
				Issue struct {
					Id    githubv4.ID
					Url   githubv4.String
					Title githubv4.String
				} `graphql:"... on Issue"`
			}
		}
		PageInfo struct {
			HasNextPage     githubv4.Boolean
			HasPreviousPage githubv4.Boolean
			EndCursor       githubv4.String
			StartCursor     githubv4.String
		}
	} `graphql:"search(query:$q, type:ISSUE, first:$first, last:$last, after:$after, before:$before)"`
}

func (g *GitHubStore) GetIssues(
	ctx context.Context,
	ids []string,
	issueType *string,
	after *string,
	before *string,
	first *int64,
	last *int64,
) (
	[]*issuesstore.Issue,
	bool,
	bool,
	string,
	string,
	error,
) {
	query := "is:issue is:open language:Go"

	var firstGql, lastGql *githubv4.Int
	if first != nil {
		tmp := githubv4.Int(*first)
		firstGql = &tmp
	}
	if last != nil {
		tmp := githubv4.Int(*last)
		lastGql = &tmp
	}

	var afterGql, beforeGql *githubv4.String
	if after != nil {
		tmp := githubv4.String(*after)
		afterGql = &tmp
	}
	if before != nil {
		tmp := githubv4.String(*before)
		beforeGql = &tmp
	}
	params := map[string]interface{}{
		"q":      githubv4.String(query),
		"after":  afterGql,
		"before": beforeGql,
		"first":  firstGql,
		"last":   lastGql,
	}
	queryData := IssueQuery{}
	if err := g.Query(ctx, &queryData, params); err != nil {
		return nil, false, false, "", "", err
	}
	issues := []*issuesstore.Issue{}

	for _, i := range queryData.Search.Edges {
		issues = append(issues, &issuesstore.Issue{
			Id:    i.Node.Issue.Id.(string),
			Title: string(i.Node.Issue.Title),
			Url:   string(i.Node.Issue.Url),
		})
	}
	return issues,
		bool(queryData.Search.PageInfo.HasNextPage),
		bool(queryData.Search.PageInfo.HasPreviousPage),
		string(queryData.Search.PageInfo.StartCursor),
		string(queryData.Search.PageInfo.EndCursor),
		nil
}
