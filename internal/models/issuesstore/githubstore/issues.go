package githubstore

import (
	"context"
	"fmt"

	"github.com/firstcontributions/backend/internal/models/issuesstore"
	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/shurcooL/githubv4"
)

const (
	IssueTypeLastRepo   = "last_repo_issues"
	IssueTypeRecentRepo = "issues_from_other_recent_repos"
	IssueTypeRelevant   = "relevant_issues"
)

func getQuery(user *usersstore.User, issueType string) (string, error) {
	query := "is:issue is:open  no:assignee"

	switch issueType {
	case IssueTypeLastRepo:
		if len(user.Tags.RecentRepos) == 0 {
			return "", fmt.Errorf("no recent repos")
		}
		return fmt.Sprintf("%s repo:%s", query, *user.Tags.RecentRepos[0]), nil
	case IssueTypeRecentRepo:
		ln := len(user.Tags.RecentRepos)
		if len(user.Tags.RecentRepos) < 2 {
			return "", fmt.Errorf("no recent repos")
		}
		count := min(7, ln)
		repos := user.Tags.RecentRepos[1:count]
		for _, repo := range repos {
			query += fmt.Sprintf(" repo:%s", *repo)
		}
		return query + " label:\"help wanted\",\"good first issue\",\"goodfirstissue\"", nil
	default:
		languageCount := min(3, len(user.Tags.Languages))
		languages := user.Tags.Languages[:languageCount]

		for _, lng := range languages {
			query += fmt.Sprintf(" language:%s", *lng)
		}
		return query + " label:\"good first issue\",\"goodfirstissue\"", nil
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type GitIssue struct {
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
}
type IssueQuery struct {
	Search struct {
		Edges []struct {
			Node struct {
				Issue GitIssue `graphql:"... on Issue"`
			}
			Cursor githubv4.String
		}
		PageInfo struct {
			HasNextPage     githubv4.Boolean
			HasPreviousPage githubv4.Boolean
			EndCursor       githubv4.String
			StartCursor     githubv4.String
		}
	} `graphql:"search(query:$q, type:ISSUE, first:$first, last:$last, after:$after, before:$before)"`
}

func (g *GitHubStore) CountIssues(
	ctx context.Context,
	filters *issuesstore.IssueFilters,
) (
	int64,
	error,
) {
	return 0, nil
}

func (g *GitHubStore) GetOneIssue(
	ctx context.Context,
	filters *issuesstore.IssueFilters,
) (
	*issuesstore.Issue,
	error,
) {
	return nil, nil
}

func (g *GitHubStore) GetIssues(
	ctx context.Context,
	filters *issuesstore.IssueFilters,
	after *string,
	before *string,
	first *int64,
	last *int64,
	sortBy issuesstore.IssueSortBy,
	sortOrder *string,
) (
	[]*issuesstore.Issue,
	bool,
	bool,
	[]string,
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
	query, err := getQuery(filters.User, *filters.IssueType)
	if err != nil {
		return nil, false, false, nil, nil
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
		return nil, false, false, nil, err
	}
	issues := []*issuesstore.Issue{}
	cursors := []string{}
	for _, i := range queryData.Search.Edges {
		issues = append(issues, issueFromGithubIssue(i.Node.Issue))
		cursors = append(cursors, string(i.Cursor))
	}
	return issues,
		bool(queryData.Search.PageInfo.HasNextPage),
		bool(queryData.Search.PageInfo.HasPreviousPage),
		cursors,
		nil
}

func issueFromGithubIssue(i GitIssue) *issuesstore.Issue {
	labels := []*string{}
	for _, label := range i.Labels.Edges {
		strLabel := string(label.Node.Name)
		labels = append(labels, &strLabel)
	}
	return &issuesstore.Issue{
		Id:                  i.Id.(string),
		Title:               string(i.Title),
		Body:                string(i.Body),
		Url:                 string(i.Url),
		Labels:              labels,
		CommentCount:        int64(i.Comments.TotalCount),
		Repository:          string(i.Repository.NameWithOwner),
		RepositoryAvatar:    string(i.Repository.Owner.AvatarUrl),
		RepositoryUpdatedAt: i.Repository.UpdatedAt.Time,
	}
}
