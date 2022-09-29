package schema

import (
	"context"
	"errors"

	"github.com/firstcontributions/backend/internal/gateway/session"
	"github.com/firstcontributions/backend/internal/storemanager"
	"github.com/firstcontributions/backend/pkg/authorizer"
)

func (m *Resolver) CreateStory(
	ctx context.Context,
	args struct {
		Story *CreateStoryInput
	},
) (*Story, error) {
	session := session.FromContext(ctx)
	if session == nil {
		return nil, errors.New("unauthorized")
	}

	storyModelInput, err := args.Story.ToModel()
	if err != nil {
		return nil, err
	}
	storyModelInput.CreatedBy = session.UserID()
	storyModelInput.UserID = session.UserID()

	ownership := &authorizer.Scope{
		Users: []string{session.UserID()},
	}
	story, err := storemanager.FromContext(ctx).StoriesStore.CreateStory(ctx, storyModelInput, ownership)
	if err != nil {
		return nil, err
	}
	return NewStory(story), nil
}
func (m *Resolver) UpdateStory(
	ctx context.Context,
	args struct {
		Story *UpdateStoryInput
	},
) (*Story, error) {
	session := session.FromContext(ctx)
	if session == nil {
		return nil, errors.New("unauthorized")
	}

	store := storemanager.FromContext(ctx)

	id, err := ParseGraphqlID(args.Story.ID)
	if err != nil {
		return nil, err
	}

	story, err := store.StoriesStore.GetStoryByID(ctx, id.ID)
	if err != nil {
		return nil, err
	}

	if !authorizer.IsAuthorized(session.Permissions, story.Ownership, authorizer.Story, authorizer.OperationUpdate) {
		return nil, errors.New("forbidden")
	}
	if err := store.StoriesStore.UpdateStory(ctx, id.ID, args.Story.ToModel()); err != nil {
		return nil, err
	}
	story, err = store.StoriesStore.GetStoryByID(ctx, id.ID)
	if err != nil {
		return nil, err
	}
	return NewStory(story), nil
}
