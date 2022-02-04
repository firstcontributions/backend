package gateway

import (
	"context"
	"log"
	"net/http"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func (s *Server) GetOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     *s.GithubConfig.ClientID,
		ClientSecret: *s.GithubConfig.ClientSecret,
		Scopes:       s.GithubConfig.AuthScopes,
		Endpoint:     github.Endpoint,
		RedirectURL:  *s.GithubConfig.AuthRedirect,
	}
}
func (s *Server) AuthRedirect(w http.ResponseWriter, r *http.Request) {
	conf := s.GetOAuth2Config()
	state, err := s.CSRFManager.Generate(r.Context())
	if err != nil {
		log.Printf("error on generating csrf state %v", err)
		ErrorResponse(ErrInternalServerError(), w)
		return
	}
	http.Redirect(w, r, conf.AuthCodeURL(state), http.StatusSeeOther)
}

func (s *Server) AuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	if code == "" || state == "" {
		ErrorResponse(ErrForbidden(), w)
		return
	}
	p, err := s.handleAuthCallback(ctx, code, state)
	if err != nil {
		ErrorResponse(err, w)
		return
	}
	s.setSession(w, r, p)
	http.Redirect(w, r, "http://explorer.firstcontributions.com", http.StatusSeeOther)
	// JSONResponse(w, http.StatusOK, p)

}

func (s *Server) handleAuthCallback(ctx context.Context, code, state string) (*usersstore.User, *Error) {
	if err := s.CSRFManager.Validate(ctx, state); err != nil {
		log.Printf("error on validating csrf %v", err)
		return nil, ErrForbidden()
	}
	conf := s.GetOAuth2Config()
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Printf("error on gettimg token by code %v", err)
		return nil, ErrForbidden()
	}
	data, err := s.getProfileFromGithub(ctx, token)
	if err != nil {
		log.Printf("error on gettimg profile from github %v", err)
		return nil, ErrInternalServerError()
	}
	users, _, _, _, _, err := s.Store.UsersStore.GetUsers(ctx, nil, nil, &data.Handle, nil, nil, nil, nil)
	if err != nil {
		log.Printf("error on gettimg profile grpc %v", err)
		return nil, ErrInternalServerError()
	}
	if len(users) == 0 {
		data, err = s.Store.UsersStore.CreateUser(ctx, data)
		if err != nil {
			log.Printf("error on creating profile grpc %v", err)
			return nil, ErrInternalServerError()
		}
	} else {
		data = users[0]
	}
	// go s.UpdateProfileReputation(profile)
	return data, nil
}

func (s *Server) getProfileFromGithub(ctx context.Context, token *oauth2.Token) (*usersstore.User, error) {

	var query struct {
		Viewer struct {
			Login     githubv4.String
			AvatarURL githubv4.URI
			Name      githubv4.String
		}
	}
	src := oauth2.StaticTokenSource(token)
	client := githubv4.NewClient(oauth2.NewClient(ctx, src))
	if err := client.Query(context.Background(), &query, nil); err != nil {
		return nil, err
	}
	return &usersstore.User{
		Name:   string(query.Viewer.Name),
		Handle: string(query.Viewer.Login),
		Token: &usersstore.Token{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
			Expiry:       token.Expiry,
		},
	}, nil
}
