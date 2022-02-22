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

type GitQuery struct {
	Viewer struct {
		Login githubv4.String
	}
	RateLimit struct {
		Limit     githubv4.Int
		Remaining githubv4.Int
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

	query := GitQuery{}

	if err := client.Query(ctx, &query, nil); err != nil {
		panic(err)
	}

	data, _ := json.MarshalIndent(query, "\t", "\t")
	fmt.Println(string(data))
}
