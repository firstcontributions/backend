package reputation

import (
	"github.com/firstcontributions/backend/internal/configs"
	"github.com/firstcontributions/backend/internal/models/usersstore"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type ReputationSynchroniser struct {
	oauthConfig *oauth2.Config
	userStore   usersstore.Store
}

func NewReputationSynchroniser(
	gitConfigs configs.GithubConfig,
	userStore usersstore.Store,
) *ReputationSynchroniser {

	return &ReputationSynchroniser{
		oauthConfig: &oauth2.Config{
			ClientID:     *gitConfigs.ClientID,
			ClientSecret: *gitConfigs.ClientSecret,
			Endpoint:     github.Endpoint,
			RedirectURL:  *gitConfigs.AuthRedirect,
			Scopes:       gitConfigs.AuthScopes,
		},
		userStore: userStore,
	}
}
