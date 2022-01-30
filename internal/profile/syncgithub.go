package profile

import (
	"context"
	"log"
	"strings"

	"github.com/firstcontributions/backend/internal/profile/models/mongo"
	"github.com/firstcontributions/backend/internal/profile/proto"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	go s.updateBadgesAndReputation(newCtx, profile)
	return &proto.EmptyRespose{}, nil
}

func (s *Service) updateBadgesAndReputation(newCtx context.Context, profile *mongo.Profile) {
	token := &oauth2.Token{
		AccessToken:  profile.Token.AccessToken,
		RefreshToken: profile.Token.RefreshToken,
		TokenType:    profile.Token.TokenType,
		Expiry:       profile.Token.Expiry,
	}

	client := githubv4.NewClient(
		s.GetOAuth2Config().Client(newCtx, token),
	)
	var p, prs []PullRequest
	hasNextPage := true
	var prCursor, lastValidCursor *githubv4.String
	var err error
	if profile.CursorCheckPoints != nil {
		tmp := githubv4.String(profile.CursorCheckPoints.PullRequest)
		prCursor = &tmp
	}
	protoProfile := profile.Proto(true)

	for hasNextPage {
		p, hasNextPage, prCursor, err = getPullRequestDataFromGitHub(client, prCursor)
		if err != nil {
			log.Printf("error on grtting data from github, %v", err)
			break
		}
		lastValidCursor = prCursor
		prs = append(prs, p...)
	}
	if lastValidCursor != nil && *lastValidCursor != "" {
		if protoProfile.CursorCheckPoints == nil {
			protoProfile.CursorCheckPoints = &proto.CursorCheckPoints{}
		}
		protoProfile.CursorCheckPoints.PullRequest = string(*lastValidCursor)
	}
	profileBadgeMap := BadgeMapFromBadges(protoProfile.Badges)
	for _, pr := range prs {
		for _, f := range pr.Files {
			profileBadgeMap.Add(f)
		}
		protoProfile.Reputation++
	}
	protoProfile.Badges = profileBadgeMap.ToBadges()

	mongo.UpdateProfile(newCtx, s.MongoClient, protoProfile)
}

type PullRequest struct {
	Files []string
}

type BadgeMap struct {
	data map[string]*proto.Badge
}

func (b *BadgeMap) Add(file string) {
	filesplit := strings.Split(file, ".")
	if len(filesplit) < 2 {
		return
	}
	switch filesplit[1] {
	case "go", "py", ".c", "cpp", "js", "ts", "java", "ruby", "rs", "sh", "php":
		if b.data[filesplit[1]] == nil {
			b.data[filesplit[1]] = &proto.Badge{
				AssignedOn: timestamppb.Now(),
				Uuid:       filesplit[1],
				Name:       filesplit[1],
			}
		}
		b.data[filesplit[1]].Progress++
	}
}

func (b *BadgeMap) Union(b2 *BadgeMap) {
	for ext, badge := range b2.data {
		if b.data[ext] == nil {
			b.data[ext] = badge
		} else {
			b.data[ext].Progress += badge.Progress
		}
	}
}

func (b *BadgeMap) ToBadges() []*proto.Badge {
	badges := []*proto.Badge{}

	for _, badge := range b.data {
		badges = append(badges, badge)
	}
	return badges
}

func BadgeMapFromBadges(badges []*proto.Badge) *BadgeMap {
	b := &BadgeMap{
		data: map[string]*proto.Badge{},
	}
	for _, bd := range badges {
		b.data[bd.Uuid] = bd
	}
	return b
}

type GitQuery struct {
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
			PageInfo struct {
				HasNextPage githubv4.Boolean
				EndCursor   githubv4.String
			}
		} `graphql:"pullRequests(first: 100, after: $prCursor)"`
	}
}

func getPullRequestDataFromGitHub(
	client *githubv4.Client,
	cursor *githubv4.String,
) (
	[]PullRequest,
	bool,
	*githubv4.String,
	error,
) {
	var query GitQuery
	params := map[string]interface{}{
		"prCursor": cursor,
	}

	if err := client.Query(context.Background(), &query, params); err != nil {
		return nil, false, nil, err
	}
	pullRequests := []PullRequest{}
	for _, pr := range query.Viewer.PullRequests.Edges {
		files := []string{}

		for _, f := range pr.Node.Files.Edges {
			files = append(files, string(f.Node.Path))
		}
		pullRequests = append(pullRequests, PullRequest{
			Files: files,
		})
	}
	return pullRequests,
		bool(query.Viewer.PullRequests.PageInfo.HasNextPage),
		&query.Viewer.PullRequests.PageInfo.EndCursor,
		nil
}
