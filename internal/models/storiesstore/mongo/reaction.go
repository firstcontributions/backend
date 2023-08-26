package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/firstcontributions/backend/internal/models/storiesstore"
	"github.com/firstcontributions/backend/internal/models/utils"
	"github.com/firstcontributions/backend/pkg/authorizer"
	"github.com/firstcontributions/backend/pkg/cursor"
	"github.com/gokultp/go-mongoqb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func reactionFiltersToQuery(filters *storiesstore.ReactionFilters) *mongoqb.QueryBuilder {
	qb := mongoqb.NewQueryBuilder()
	if len(filters.Ids) > 0 {
		qb.In("_id", filters.Ids)
	}
	if filters.CreatedBy != nil {
		qb.Eq("created_by", filters.CreatedBy)
	}
	if filters.Story != nil {
		qb.Eq("story_id", filters.Story.Id)
	}
	return qb
}
func (s *StoriesStore) CreateReaction(ctx context.Context, reaction *storiesstore.Reaction, ownership *authorizer.Scope) (*storiesstore.Reaction, error) {
	now := time.Now()
	reaction.TimeCreated = now
	reaction.TimeUpdated = now
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	if dproc, ok := interface{}(reaction).(utils.DataProcessor); ok {
		if err := dproc.Process(ctx); err != nil {
			return nil, err
		}
	}
	reaction.Id = uuid.String()
	reaction.Ownership = ownership
	if _, err := s.getCollection(CollectionReactions).InsertOne(ctx, reaction); err != nil {
		return nil, err
	}
	return reaction, nil
}

func (s *StoriesStore) GetReactionByID(ctx context.Context, id string) (*storiesstore.Reaction, error) {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	var reaction storiesstore.Reaction
	if err := s.getCollection(CollectionReactions).FindOne(ctx, qb.Build()).Decode(&reaction); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &reaction, nil
}

func (s *StoriesStore) GetOneReaction(ctx context.Context, filters *storiesstore.ReactionFilters) (*storiesstore.Reaction, error) {
	qb := reactionFiltersToQuery(filters)
	var reaction storiesstore.Reaction
	if err := s.getCollection(CollectionReactions).FindOne(ctx, qb.Build()).Decode(&reaction); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &reaction, nil
}

func (s *StoriesStore) CountReactions(ctx context.Context, filters *storiesstore.ReactionFilters) (
	int64,
	error,
) {
	qb := reactionFiltersToQuery(filters)

	count, err := s.getCollection(CollectionReactions).CountDocuments(ctx, qb.Build())
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *StoriesStore) GetReactions(
	ctx context.Context,
	filters *storiesstore.ReactionFilters,
	after *string,
	before *string,
	first *int64,
	last *int64,
	sortBy storiesstore.ReactionSortBy,
	sortOrder *string,
) (
	[]*storiesstore.Reaction,
	bool,
	bool,
	[]string,
	error,
) {
	qb := reactionFiltersToQuery(filters)
	reqSortOrder := utils.GetSortOrderFromString(sortOrder)
	limit, paginationSortOrder, cursorStr, err := utils.GetLimitAndSortOrderAndCursor(first, last, after, before)
	if err != nil {
		return nil, false, false, nil, err
	}

	effectiveSortOrder := reqSortOrder * paginationSortOrder

	var c *cursor.Cursor
	if cursorStr != nil {
		c, err = cursor.FromString(*cursorStr)
		if err != nil {
			return nil, false, false, nil, err
		}
		if c != nil {
			if effectiveSortOrder == 1 {
				qb.Or(

					mongoqb.NewQueryBuilder().
						Eq(storiesstore.ReactionSortBy(c.SortBy).String(), c.OffsetValue).
						Gt("_id", c.ID),
					mongoqb.NewQueryBuilder().
						Gt(storiesstore.ReactionSortBy(c.SortBy).String(), c.OffsetValue),
				)
			} else {
				qb.Or(
					mongoqb.NewQueryBuilder().
						Eq(storiesstore.ReactionSortBy(c.SortBy).String(), c.OffsetValue).
						Lt("_id", c.ID),
					mongoqb.NewQueryBuilder().
						Lt(storiesstore.ReactionSortBy(c.SortBy).String(), c.OffsetValue),
				)
			}
		}
	}
	// incrementing limit by 2 to check if next, prev elements are present
	limit += 2
	options := &options.FindOptions{
		Limit: &limit,
		Sort:  utils.GetSortOrder(sortBy.String(), effectiveSortOrder),
	}

	var hasNextPage, hasPreviousPage bool

	var reactions []*storiesstore.Reaction
	mongoCursor, err := s.getCollection(CollectionReactions).Find(ctx, qb.Build(), options)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, nil, err
	}
	err = mongoCursor.All(ctx, &reactions)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, nil, err
	}
	count := len(reactions)
	if count == 0 {
		return reactions, hasNextPage, hasPreviousPage, nil, nil
	}

	// check if the cursor element present, if yes that can be a prev elem
	if c != nil && reactions[0].Id == c.ID {
		hasPreviousPage = true
		reactions = reactions[1:]
		count--
	}

	// check if actual limit +1 elements are there, if yes trim it to limit
	if count >= int(limit)-1 {
		hasNextPage = true
		reactions = reactions[:limit-2]
		count = len(reactions)
	}

	cursors := make([]string, count)
	for i, reaction := range reactions {
		cursors[i] = cursor.NewCursor(reaction.Id, uint8(sortBy), reaction.Get(sortBy.String()), sortBy.CursorType()).String()
	}

	if paginationSortOrder < 0 {
		hasNextPage, hasPreviousPage = hasPreviousPage, hasNextPage
		reactions = utils.ReverseList(reactions)
	}
	return reactions, hasNextPage, hasPreviousPage, cursors, nil
}

func (s *StoriesStore) UpdateReaction(ctx context.Context, id string, reactionUpdate *storiesstore.ReactionUpdate) error {

	if dproc, ok := interface{}(reactionUpdate).(utils.DataProcessor); ok {
		if err := dproc.Process(ctx); err != nil {
			return err
		}
	}
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)

	now := time.Now()
	reactionUpdate.TimeUpdated = &now

	u := mongoqb.NewUpdateMap().
		SetFields(reactionUpdate)

	um, err := u.BuildUpdate()
	if err != nil {
		return err
	}
	if _, err := s.getCollection(CollectionReactions).UpdateOne(ctx, qb.Build(), um); err != nil {
		return err
	}
	return nil
}

func (s *StoriesStore) DeleteReactionByID(ctx context.Context, id string) error {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	if _, err := s.getCollection(CollectionReactions).DeleteOne(ctx, qb.Build()); err != nil {
		return err
	}
	return nil
}
