package reputation

import (
	"github.com/firstcontributions/backend/internal/configs"
	"github.com/firstcontributions/backend/internal/githubclient"
	"github.com/firstcontributions/backend/internal/models/usersstore"
)

type ReputationSynchroniser struct {
	*githubclient.GitHubClient
	userStore usersstore.Store
}

func NewReputationSynchroniser(
	gitConfigs configs.GithubConfig,
	userStore usersstore.Store,
) *ReputationSynchroniser {

	return &ReputationSynchroniser{
		GitHubClient: githubclient.NewGitHubClient(gitConfigs),
		userStore:    userStore,
	}
}
