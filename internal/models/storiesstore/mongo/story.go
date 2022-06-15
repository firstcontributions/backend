package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/firstcontributions/backend/internal/models/storiesstore"
	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/firstcontributions/backend/internal/models/utils"
	"github.com/firstcontributions/backend/pkg/cursor"
	"github.com/gokultp/go-mongoqb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *StoriesStore) CreateStory(ctx context.Context, story *storiesstore.Story) (*storiesstore.Story, error) {
	now := time.Now()
	story.TimeCreated = now
	story.TimeUpdated = now
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	story.Id = uuid.String()
	if _, err := s.getCollection(CollectionStories).InsertOne(ctx, story); err != nil {
		return nil, err
	}
	return story, nil
}

func (s *StoriesStore) GetStoryByID(ctx context.Context, id string) (*storiesstore.Story, error) {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	var story storiesstore.Story
	if err := s.getCollection(CollectionStories).FindOne(ctx, qb.Build()).Decode(&story); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &story, nil
}

func (s *StoriesStore) GetStories(
	ctx context.Context,
	ids []string,
	user *usersstore.User,
	after *string,
	before *string,
	first *int64,
	last *int64,
) (
	[]*storiesstore.Story,
	bool,
	bool,
	string,
	string,
	error,
) {
	qb := mongoqb.NewQueryBuilder()
	if len(ids) > 0 {
		qb.In("_id", ids)
	}
	if user != nil {
		qb.Eq("user_id", user.Id)
	}

	limit, order, cursorStr := utils.GetLimitAndSortOrderAndCursor(first, last, after, before)
	var c *cursor.Cursor
	if cursorStr != nil {
		c = cursor.FromString(*cursorStr)
		if c != nil {
			if order == 1 {
				qb.Lte("time_created", c.TimeStamp)
				qb.Lte("_id", c.ID)
			} else {
				qb.Gte("time_created", c.TimeStamp)
				qb.Gte("_id", c.ID)
			}
		}
	}
	sortOrder := utils.GetSortOrder(order)
	// incrementing limit by 2 to check if next, prev elements are present
	limit += 2
	options := &options.FindOptions{
		Limit: &limit,
		Sort:  sortOrder,
	}

	var firstCursor, lastCursor string
	var hasNextPage, hasPreviousPage bool

	var stories []*storiesstore.Story
	mongoCursor, err := s.getCollection(CollectionStories).Find(ctx, qb.Build(), options)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err
	}
	err = mongoCursor.All(ctx, &stories)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err
	}
	count := len(stories)
	if count == 0 {
		return stories, hasNextPage, hasPreviousPage, firstCursor, lastCursor, nil
	}

	// check if the cursor element present, if yes that can be a prev elem
	if c != nil && stories[0].Id == c.ID {
		hasPreviousPage = true
		stories = stories[1:]
		count--
	}

	// check if actual limit +1 elements are there, if yes trim it to limit
	if count >= int(limit)-1 {
		hasNextPage = true
		stories = stories[:limit-2]
		count = len(stories)
	}

	if count > 0 {
		firstCursor = cursor.NewCursor(stories[0].Id, stories[0].TimeCreated).String()
		lastCursor = cursor.NewCursor(stories[count-1].Id, stories[count-1].TimeCreated).String()
	}
	if order < 0 {
		hasNextPage, hasPreviousPage = hasPreviousPage, hasNextPage
		firstCursor, lastCursor = lastCursor, firstCursor
	}
	return stories, hasNextPage, hasPreviousPage, firstCursor, lastCursor, nil
}
func (s *StoriesStore) UpdateStory(ctx context.Context, id string, storyUpdate *storiesstore.StoryUpdate) error {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)

	now := time.Now()
	storyUpdate.TimeUpdated = &now

	u := mongoqb.NewUpdateMap().
		SetFields(storyUpdate)

	um, err := u.BuildUpdate()
	if err != nil {
		return err
	}
	if _, err := s.getCollection(CollectionStories).UpdateOne(ctx, qb.Build(), um); err != nil {
		return err
	}
	return nil
}

func (s *StoriesStore) DeleteStoryByID(ctx context.Context, id string) error {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	if _, err := s.getCollection(CollectionStories).DeleteOne(ctx, qb.Build()); err != nil {
		return err
	}
	return nil
}