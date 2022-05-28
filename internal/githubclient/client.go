package githubclient

import (
	"context"
	"fmt"

	"github.com/firstcontributions/backend/internal/configs"
	"github.com/firstcontributions/backend/internal/gateway/session"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GitHubClient struct {
	configs *oauth2.Config
}

func NewGitHubClient(config configs.GithubConfig) *GitHubClient {
	return &GitHubClient{
		configs: &oauth2.Config{
			ClientID:     *config.ClientID,
			ClientSecret: *config.ClientSecret,
			Scopes:       config.AuthScopes,
			RedirectURL:  *config.AuthRedirect,
			Endpoint:     github.Endpoint,
		},
	}
}

func (g *GitHubClient) Query(ctx context.Context, query interface{}, params map[string]interface{}) error {
	meta := session.FromContext(ctx)
	fmt.Println("-------------- meta", meta)
	token := &oauth2.Token{
		AccessToken: meta.Token.AccessToken,
	}
	return githubv4.NewClient(g.configs.Client(ctx, token)).Query(ctx, query, params)
}
