package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/firstcontributions/backend/internal/models/storiesstore"
	"github.com/firstcontributions/backend/internal/models/utils"
	"github.com/firstcontributions/backend/pkg/cursor"
	"github.com/gokultp/go-mongoqb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func commentFiltersToQuery(filters *storiesstore.CommentFilters) *mongoqb.QueryBuilder {
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
func (s *StoriesStore) CreateComment(ctx context.Context, comment *storiesstore.Comment) (*storiesstore.Comment, error) {
	now := time.Now()
	comment.TimeCreated = now
	comment.TimeUpdated = now
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	comment.Id = uuid.String()
	if _, err := s.getCollection(CollectionComments).InsertOne(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *StoriesStore) GetCommentByID(ctx context.Context, id string) (*storiesstore.Comment, error) {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	var comment storiesstore.Comment
	if err := s.getCollection(CollectionComments).FindOne(ctx, qb.Build()).Decode(&comment); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

func (s *StoriesStore) GetOneComment(ctx context.Context, filters *storiesstore.CommentFilters) (*storiesstore.Comment, error) {
	qb := commentFiltersToQuery(filters)
	var comment storiesstore.Comment
	if err := s.getCollection(CollectionComments).FindOne(ctx, qb.Build()).Decode(&comment); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

func (s *StoriesStore) CountComments(ctx context.Context, filters *storiesstore.CommentFilters) (
	int64,
	error,
) {
	qb := commentFiltersToQuery(filters)

	count, err := s.getCollection(CollectionComments).CountDocuments(ctx, qb.Build())
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *StoriesStore) GetComments(
	ctx context.Context,
	filters *storiesstore.CommentFilters,
	after *string,
	before *string,
	first *int64,
	last *int64,
	sortBy storiesstore.CommentSortBy,
	sortOrder *string,
) (
	[]*storiesstore.Comment,
	bool,
	bool,
	[]string,
	error,
) {
	qb := commentFiltersToQuery(filters)
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
						Eq(storiesstore.CommentSortBy(c.SortBy).String(), c.OffsetValue).
						Gt("_id", c.ID),
					mongoqb.NewQueryBuilder().
						Gt(storiesstore.CommentSortBy(c.SortBy).String(), c.OffsetValue),
				)
			} else {
				qb.Or(
					mongoqb.NewQueryBuilder().
						Eq(storiesstore.CommentSortBy(c.SortBy).String(), c.OffsetValue).
						Lt("_id", c.ID),
					mongoqb.NewQueryBuilder().
						Lt(storiesstore.CommentSortBy(c.SortBy).String(), c.OffsetValue),
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

	var comments []*storiesstore.Comment
	mongoCursor, err := s.getCollection(CollectionComments).Find(ctx, qb.Build(), options)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, nil, err
	}
	err = mongoCursor.All(ctx, &comments)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, nil, err
	}
	count := len(comments)
	if count == 0 {
		return comments, hasNextPage, hasPreviousPage, nil, nil
	}

	// check if the cursor element present, if yes that can be a prev elem
	if c != nil && comments[0].Id == c.ID {
		hasPreviousPage = true
		comments = comments[1:]
		count--
	}

	// check if actual limit +1 elements are there, if yes trim it to limit
	if count >= int(limit)-1 {
		hasNextPage = true
		comments = comments[:limit-2]
		count = len(comments)
	}

	cursors := make([]string, count)
	for i, comment := range comments {
		cursors[i] = cursor.NewCursor(comment.Id, uint8(sortBy), comment.Get(sortBy.String()), sortBy.CursorType()).String()
	}

	if paginationSortOrder < 0 {
		hasNextPage, hasPreviousPage = hasPreviousPage, hasNextPage
		comments = utils.ReverseList(comments)
	}
	return comments, hasNextPage, hasPreviousPage, cursors, nil
}

func (s *StoriesStore) UpdateComment(ctx context.Context, id string, commentUpdate *storiesstore.CommentUpdate) error {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)

	now := time.Now()
	commentUpdate.TimeUpdated = &now

	u := mongoqb.NewUpdateMap().
		SetFields(commentUpdate)

	um, err := u.BuildUpdate()
	if err != nil {
		return err
	}
	if _, err := s.getCollection(CollectionComments).UpdateOne(ctx, qb.Build(), um); err != nil {
		return err
	}
	return nil
}

func (s *StoriesStore) DeleteCommentByID(ctx context.Context, id string) error {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	if _, err := s.getCollection(CollectionComments).DeleteOne(ctx, qb.Build()); err != nil {
		return err
	}
	return nil
}
