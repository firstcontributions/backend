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
	query := "is:issue is:open  label:\"help wanted\",\"good first issue\",\"goodfirstissue\" no:assignee"

	switch issueType {
	case IssueTypeLastRepo:
		return fmt.Sprintf("%s repo:%s", query, meta.Tags.RecentRepos[0])
	case IssueTypeRecentRepo:
		ln := len(meta.Tags.RecentRepos)
		count := min(7, ln)
		repos := meta.Tags.RecentRepos[1:count]
		for _, repo := range repos {
			query += fmt.Sprintf(" repo:%s", repo)
		}
		return query
	default:
		languageCount := min(3, len(meta.Tags.Languages))
		languages := meta.Tags.Languages[:languageCount]

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
					Id       githubv4.ID
					Url      githubv4.String
					Title    githubv4.String
					Body     githubv4.String
					Comments struct {
						TotalCount githubv4.Int
					}
					Labels struct {
						Edges []struct {
							Node struct {
								Name githubv4.String
							}
						}
					} `graphql:"labels (first: 10)"`
					Repository struct {
						NameWithOwner githubv4.String
						Owner         struct {
							AvatarUrl githubv4.String
						}
						UpdatedAt githubv4.DateTime
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
		labels := []*string{}
		for _, label := range i.Node.Issue.Labels.Edges {
			strLabel := string(label.Node.Name)
			labels = append(labels, &strLabel)
		}
		issues = append(issues, &issuesstore.Issue{
			Id:                  i.Node.Issue.Id.(string),
			Title:               string(i.Node.Issue.Title),
			Body:                string(i.Node.Issue.Body),
			Url:                 string(i.Node.Issue.Url),
			Labels:              labels,
			CommentCount:        int64(i.Node.Issue.Comments.TotalCount),
			Repository:          string(i.Node.Issue.Repository.NameWithOwner),
			RespositoryAvatar:   string(i.Node.Issue.Repository.Owner.AvatarUrl),
			RepositoryUpdatedAt: i.Node.Issue.Repository.UpdatedAt.Time,
		})
	}
	return issues,
		bool(queryData.Search.PageInfo.HasNextPage),
		bool(queryData.Search.PageInfo.HasPreviousPage),
		string(queryData.Search.PageInfo.StartCursor),
		string(queryData.Search.PageInfo.EndCursor),
		nil
}
