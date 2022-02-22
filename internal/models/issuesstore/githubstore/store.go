package githubstore

import (
	"github.com/firstcontributions/backend/internal/configs"
	"github.com/firstcontributions/backend/internal/githubclient"
)

type GitHubStore struct {
	*githubclient.GitHubClient
}

func NewGitHubStore(config configs.GithubConfig) *GitHubStore {
	return &GitHubStore{
		GitHubClient: githubclient.NewGitHubClient(config),
	}
}
