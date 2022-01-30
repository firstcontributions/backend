package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/firstcontributions/backend/internal/profile/models/mongo"
	"github.com/firstcontributions/backend/internal/profile/proto"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     *s.GithubConfig.ClientID,
		ClientSecret: *s.GithubConfig.ClientSecret,
		Scopes:       s.GithubConfig.AuthScopes,
		Endpoint:     github.Endpoint,
		RedirectURL:  *s.GithubConfig.AuthRedirect,
	}

}

func (s *Service) SyncGitHubData(ctx context.Context, in *proto.GetByUUIDRequest) (*proto.EmptyRespose, error) {

	newCtx := context.Background()

	profile, err := mongo.GetProfileByUUID(newCtx, s.MongoClient, in.Uuid)
	if err != nil {
		log.Printf("error on reading profile data, %v", err)
		return nil, status.Errorf(codes.Internal, "error on mongo query %w", err)
	}

	token := &oauth2.Token{
		AccessToken:  profile.Token.AccessToken,
		RefreshToken: profile.Token.RefreshToken,
		TokenType:    profile.Token.TokenType,
		Expiry:       profile.Token.Expiry,
	}

	client := githubv4.NewClient(
		s.GetOAuth2Config().Client(newCtx, token),
	)

	query, params := buildGraphQLQuery(profile)
	if err := client.Query(context.Background(), query, params); err != nil {
		log.Printf("error on getting user data from github %v", err)
		return nil, status.Errorf(codes.Internal, "error on github query %w", err)
	}
	j, _ := json.MarshalIndent(query, "", "	")

	fmt.Println(string(j))
	return &proto.EmptyRespose{}, nil
}

func buildGraphQLQuery(profile *mongo.Profile) (interface{}, map[string]interface{}) {
	var q struct {
		Viewer struct {
			Login githubv4.String

			PullRequests struct {
				Edges []struct {
					Node struct {
						Files struct {
							Edges []struct {
								Node struct {
									Path githubv4.String
								}
							}
						} `graphql:"files(first: 20)"`
					}
				}
			} `graphql:"pullRequests(first: 100, after: $prCursor)"`
		}
	}

	var prCusrsor *githubv4.String

	if profile.CursorCheckPoints != nil {
		tmp := githubv4.String(profile.CursorCheckPoints.PullRequest)
		prCusrsor = &tmp
	}
	params := map[string]interface{}{
		"prCursor": prCusrsor,
	}

	return &q, params
}
