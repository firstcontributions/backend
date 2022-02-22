package reputation

import (
	"context"
	"log"
	"time"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/shurcooL/githubv4"
)

func (r ReputationSynchroniser) SyncBadges(ctx context.Context, user *usersstore.User) error {
	start := time.Now()
	existingBadges, _, _, _, _, err := r.userStore.GetBadges(ctx, nil, &user.Id, nil, nil, nil, nil)
	if err != nil {
		return nil
	}

	badgeMap := BadgeMapFromBadges(existingBadges)

	fileChanges, cursor := r.getPRFileChangesFromGitHub(ctx, user)

	for _, f := range fileChanges {
		badgeMap.Add(f.path, f.additions)
	}
	for _, badge := range badgeMap.data {
		if err := r.updateBadge(ctx, badge, user); err != nil {
			return nil
		}
	}
	log.Printf("look %d ms\n", (time.Since(start))*time.Millisecond)
	if cursor != nil {
		userUpdate := &usersstore.UserUpdate{
			CursorCheckpoints: &usersstore.CursorCheckpoints{PullRequests: string(*cursor)},
		}
		return r.userStore.UpdateUser(ctx, user.Id, userUpdate)
	}
	return nil
}

func (r ReputationSynchroniser) updateBadge(ctx context.Context, badge *usersstore.Badge, user *usersstore.User) error {
	badge.CurrentLevel = int64(GetLevelFromPoints(int(badge.Points)))
	badge.UserID = user.Id
	badge.ProgressPercentageToNextLevel = GetProgressPercentageToNextLevel(int(badge.Points))

	if badge.Id == "" {
		_, err := r.userStore.CreateBadge(ctx, badge)
		return err
	}
	update := &usersstore.BadgeUpdate{
		CurrentLevel:                  &badge.CurrentLevel,
		Points:                        &badge.Points,
		ProgressPercentageToNextLevel: &badge.ProgressPercentageToNextLevel,
	}
	return r.userStore.UpdateBadge(ctx, badge.Id, update)
}

func (r ReputationSynchroniser) getPRFileChangesFromGitHub(ctx context.Context, user *usersstore.User) ([]FileChange, *githubv4.String) {

	var f, fchanges []FileChange
	hasNextPage := true
	var prCursor, lastValidCursor *githubv4.String
	var err error

	if user.CursorCheckpoints != nil && user.CursorCheckpoints.PullRequests != "" {
		tmp := githubv4.String(user.CursorCheckpoints.PullRequests)
		prCursor = &tmp
	}
	for hasNextPage {
		f, hasNextPage, prCursor, err = r.getPullRequestDataFromGitHub(ctx, prCursor)
		if err != nil {
			log.Printf("error on grtting data from github, %v", err)
			break
		}
		lastValidCursor = prCursor
		fchanges = append(fchanges, f...)
	}
	return fchanges, lastValidCursor
}

type FileChange struct {
	path      string
	additions int
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
								Path      githubv4.String
								Additions githubv4.Int
							}
						}
					} `graphql:"files(first: 100)"`
				}
			}
			PageInfo struct {
				HasNextPage githubv4.Boolean
				EndCursor   githubv4.String
			}
		} `graphql:"pullRequests(first: 100, after: $prCursor)"`
	}
}

func (r ReputationSynchroniser) getPullRequestDataFromGitHub(
	ctx context.Context,
	cursor *githubv4.String,
) (
	[]FileChange,
	bool,
	*githubv4.String,
	error,
) {
	var query GitQuery
	params := map[string]interface{}{
		"prCursor": cursor,
	}

	if err := r.Query(context.Background(), &query, params); err != nil {
		return nil, false, nil, err
	}
	files := []FileChange{}
	for _, pr := range query.Viewer.PullRequests.Edges {

		for _, f := range pr.Node.Files.Edges {
			files = append(files, FileChange{
				path:      string(f.Node.Path),
				additions: int(f.Node.Additions),
			})
		}
	}
	return files,
		bool(query.Viewer.PullRequests.PageInfo.HasNextPage),
		&query.Viewer.PullRequests.PageInfo.EndCursor,
		nil
}
