package gateway

import (
	"context"
	"log"
	"net/http"
	"time"

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
	origin := "http://explorer.firstcontributions.com"
	if originFromQuery := r.URL.Query().Get("origin"); originFromQuery != "" {
		origin = originFromQuery
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "fc_origin",
		Value:   origin,
		Expires: time.Now().Add(5 * time.Minute * 5),
		Path:    "/",
		Domain:  "firstcontributions.com",
	})
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
	sessionID, sessionErr := s.setSession(w, r, p)
	if sessionErr != nil {
		log.Println(sessionErr)
		ErrorResponse(ErrInternalServerError(), w)
	}
	go s.UpdateProfileReputation(p, sessionID)
	redirect := "http://explorer.firstcontributions.com"
	cookie, _ := r.Cookie("fc_origin")
	if cookie != nil && cookie.Value != "" {
		redirect = cookie.Value
	}
	http.Redirect(w, r, redirect, http.StatusSeeOther)
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
	profile, err := s.getProfileFromGithub(ctx, token)
	if err != nil {
		log.Printf("error on gettimg profile from github %v", err)
		return nil, ErrInternalServerError()
	}
	filters := &usersstore.UserFilters{Handle: &profile.Handle}
	user, err := s.Store.UsersStore.GetOneUser(ctx, filters)
	if err != nil {
		log.Printf("error on gettimg profile grpc %v", err)
		return nil, ErrInternalServerError()
	}
	if user == nil {
		data, err := s.Store.UsersStore.CreateUser(ctx, profile, nil)
		if err != nil {
			log.Printf("error on creating profile grpc %v", err)
			return nil, ErrInternalServerError()
		}
		return data, nil
	}
	data := user
	data.Token = profile.Token
	data.Avatar = profile.Avatar
	go s.Store.UsersStore.UpdateUser(ctx, data.Id, &usersstore.UserUpdate{
		Token:  profile.Token,
		Avatar: &profile.Avatar,
	})

	return data, nil
}

func (s *Server) getProfileFromGithub(ctx context.Context, token *oauth2.Token) (*usersstore.User, error) {

	var query struct {
		Viewer struct {
			Login     githubv4.String
			AvatarURL githubv4.URI
			Bio       githubv4.String
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
		Avatar: query.Viewer.AvatarURL.String(),
		Bio:    string(query.Viewer.Bio),
		Token: &usersstore.Token{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
			Expiry:       token.Expiry,
		},
	}, nil
}
