package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/firstcontributions/backend/internal/configs"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

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
		} `graphql:"repositoriesContributedTo(first: 10, contributionTypes: COMMIT, orderBy: {field: PUSHED_AT, direction: DESC})"`
	}
}

func main() {
	c := configs.GithubConfig{}
	if err := c.DecodeEnv(); err != nil {
		panic(err)
	}
	authC := &oauth2.Config{
		ClientID:     *c.ClientID,
		ClientSecret: *c.ClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  *c.AuthRedirect,
		Scopes:       c.AuthScopes,
	}

	token := &oauth2.Token{
		AccessToken: "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",
	}
	ctx := context.Background()
	client := githubv4.NewClient(
		authC.Client(ctx, token),
	)

	query := TagsQuery{}

	if err := client.Query(ctx, &query, nil); err != nil {
		panic(err)
	}

	data, _ := json.MarshalIndent(query, "\t", "\t")
	fmt.Println(string(data))
}
