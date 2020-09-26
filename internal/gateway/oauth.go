package gateway

import (
	"context"
	"log"
	"net/http"

	"github.com/firstcontributions/firstcontributions/internal/proto"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	if err := s.CSRFManager.Validate(ctx, state); err != nil {
		log.Printf("error on validating csrf %v", err)
		ErrorResponse(ErrForbidden(), w)
		return
	}
	conf := s.GetOAuth2Config()
	token, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Printf("error on gettimg token by code %v", err)
		ErrorResponse(ErrForbidden(), w)
		return
	}
	data, err := s.getProfileFromGithub(ctx, token)
	if err != nil {
		log.Printf("error on gettimg profile from github %v", err)
		ErrorResponse(ErrInternalServerError(), w)
		return
	}
	profile, err := s.getProfileByHandle(ctx, data.Handle)
	if err != nil && grpc.Code(err) != codes.NotFound {
		log.Printf("error on gettimg profile grpc %v", err)
		ErrorResponse(ErrInternalServerError(), w)
		return
	}
	if profile == nil || grpc.Code(err) == codes.NotFound {
		profile, err = s.createProfile(ctx, data)
		if err != nil {
			log.Printf("error on gettimg profile grpc %v", err)
			ErrorResponse(ErrInternalServerError(), w)
			return
		}
	}
	log.Print("user id", profile.Uuid)
	if err := s.setSession(w, r, profile); err != nil {
		log.Printf("error on setting session %v", err)
		ErrorResponse(ErrInternalServerError(), w)
		return
	}
	JSONResponse(w, http.StatusAccepted, data)
}

func (s *Server) getProfileFromGithub(ctx context.Context, token *oauth2.Token) (*proto.Profile, error) {

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
	return &proto.Profile{
		Name:   string(query.Viewer.Name),
		Avatar: query.Viewer.AvatarURL.String(),
		Handle: string(query.Viewer.Login),
		GithubToken: &proto.Token{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
			Expiry:       timestamppb.New(token.Expiry),
		},
	}, nil
}

func (s *Server) getProfileByHandle(ctx context.Context, handle string) (*proto.Profile, error) {
	conn, err := s.ProfileConn.Get(ctx)
	if err != nil {
		return nil, err
	}
	req := proto.GetProfileRequest{
		Handle: handle,
	}
	return proto.NewProfileServiceClient(conn).GetProfile(ctx, &req)
}

func (s *Server) createProfile(ctx context.Context, profile *proto.Profile) (*proto.Profile, error) {
	conn, err := s.ProfileConn.Get(ctx)
	if err != nil {
		return nil, err
	}

	return proto.NewProfileServiceClient(conn).CreateProfile(ctx, profile)
}
