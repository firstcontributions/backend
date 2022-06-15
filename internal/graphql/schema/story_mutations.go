package schema

import (
	"context"
	"errors"

	"github.com/firstcontributions/backend/internal/gateway/session"
	"github.com/firstcontributions/backend/internal/storemanager"
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

	story, err := storemanager.FromContext(ctx).StoriesStore.CreateStory(ctx, storyModelInput)
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
	store := storemanager.FromContext(ctx)

	id, err := ParseGraphqlID(args.Story.ID)
	if err != nil {
		return nil, err
	}
	if err := store.StoriesStore.UpdateStory(ctx, id.ID, args.Story.ToModel()); err != nil {
		return nil, err
	}
	story, err := store.StoriesStore.GetStoryByID(ctx, id.ID)
	if err != nil {
		return nil, err
	}
	return NewStory(story), nil
}