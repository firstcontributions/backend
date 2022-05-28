package githubstore

import (
	"context"
	"fmt"

	"github.com/firstcontributions/backend/internal/gateway/session"
	"github.com/firstcontributions/backend/internal/models/issuesstore"
	"github.com/shurcooL/githubv4"
)

const (
	IssueTypeLastRepo   = "last_repo_issues"
	IssueTypeRecentRepo = "issues_from_other_recent_repos"
	IssueTypeRelevant   = "relevant_issues"
)

func getQuery(ctx context.Context, issueType string) string {
	meta := session.FromContext(ctx)

	switch issueType {
	case IssueTypeLastRepo:
		return fmt.Sprintf("is:issue is:open  label:\"help wanted\",\"good first issue\",\"goodfirstissue\"  repo:%s", meta.Tags.RecentRepos[0])
	case IssueTypeRecentRepo:
		ln := len(meta.Tags.RecentRepos)
		count := min(7, ln)
		repos := meta.Tags.RecentRepos[1:count]
		query := "is:issue is:open  label:\"help wanted\",\"good first issue\",\"goodfirstissue\""
		for _, repo := range repos {
			query += fmt.Sprintf(" repo:%s", repo)
		}
		return query
	default:
		languageCount := min(3, len(meta.Tags.Languages))
		languages := meta.Tags.Languages[:languageCount]

		query := "is:issue is:open  label:\"help wanted\",\"good first issue\",\"goodfirstissue\""
		for _, lng := range languages {
			query += fmt.Sprintf(" language:%s", lng)
		}
		return query
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type IssueQuery struct {
	Search struct {
		Edges []struct {
			Node struct {
				Issue struct {
					Id         githubv4.ID
					Url        githubv4.String
					Title      githubv4.String
					Repository struct {
						NameWithOwner githubv4.String
						Owner         struct {
							AvatarUrl githubv4.String
						}
					}
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
	fmt.Println(getQuery(ctx, *issueType))
	params := map[string]interface{}{
		"q":      githubv4.String(getQuery(ctx, *issueType)),
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
			Id:                i.Node.Issue.Id.(string),
			Title:             string(i.Node.Issue.Title),
			Url:               string(i.Node.Issue.Url),
			Repository:        string(i.Node.Issue.Repository.NameWithOwner),
			RespositoryAvatar: string(i.Node.Issue.Repository.Owner.AvatarUrl),
		})
	}
	return issues,
		bool(queryData.Search.PageInfo.HasNextPage),
		bool(queryData.Search.PageInfo.HasPreviousPage),
		string(queryData.Search.PageInfo.StartCursor),
		string(queryData.Search.PageInfo.EndCursor),
		nil
}
