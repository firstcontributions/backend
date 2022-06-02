package reputation

import (
	"context"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/firstcontributions/backend/pkg/sets"
	"github.com/shurcooL/githubv4"
)

func (r ReputationSynchroniser) SyncTags(ctx context.Context, user *usersstore.User) error {
	query := TagsQuery{}
	if err := r.Query(ctx, &query, nil); err != nil {
		return err
	}
	languages := sets.NewSet()
	topics := sets.NewSet()
	repositories := []*string{}

	for _, repo := range query.Viewer.RepositoriesContributedTo.Edges {
		repoName := string(repo.Node.NameWithOwner)
		repositories = append(repositories, &repoName)
		if repo.Node.PrimaryLanguage.Name != "" {
			languages.Add(string(repo.Node.PrimaryLanguage.Name))
		}
		for _, topic := range repo.Node.RepositoryTopics.Edges {
			topics.Add(string(topic.Node.Topic.Name))
		}
	}
	updateObject := usersstore.UserUpdate{
		Tags: &usersstore.Tags{
			Languages:   languages.Elems(),
			Topics:      topics.Elems(),
			RecentRepos: repositories,
		},
	}
	return r.userStore.UpdateUser(ctx, user.Id, &updateObject)
}

type TagsQuery struct {
	Viewer struct {
		Login githubv4.String

		RepositoriesContributedTo struct {
			Edges []struct {
				Node struct {
					NameWithOwner   githubv4.String
					PrimaryLanguage struct {
						Name githubv4.String
					}
					RepositoryTopics struct {
						Edges []struct {
							Node struct {
								Topic struct {
									Name githubv4.String
								}
							}
						}
					} `graphql:"repositoryTopics(first: 3)"`
				}
			}
		} `graphql:"repositoriesContributedTo(first: 10, includeUserRepositories: true, contributionTypes: COMMIT, orderBy: {field: PUSHED_AT, direction: DESC})"`
	}
}
