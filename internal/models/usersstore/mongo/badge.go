package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/firstcontributions/backend/internal/models/utils"
	"github.com/firstcontributions/backend/pkg/cursor"
	"github.com/gokultp/go-mongoqb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *UsersStore) CreateBadge(ctx context.Context, badge *usersstore.Badge) (*usersstore.Badge, error) {
	now := time.Now()
	badge.TimeCreated = now
	badge.TimeUpdated = now
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	badge.Id = uuid.String()
	if _, err := s.getCollection(CollectionBadges).InsertOne(ctx, badge); err != nil {
		return nil, err
	}
	return badge, nil
}

func (s *UsersStore) GetBadgeByID(ctx context.Context, id string) (*usersstore.Badge, error) {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	var badge usersstore.Badge
	if err := s.getCollection(CollectionBadges).FindOne(ctx, qb.Build()).Decode(&badge); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &badge, nil
}

func (s *UsersStore) GetBadges(
	ctx context.Context,
	ids []string,
	userID *string,
	after *string,
	before *string,
	first *int64,
	last *int64,
) (
	[]*usersstore.Badge,
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
	if userID != nil {
		qb.Eq("user_id", userID)
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

	var badges []*usersstore.Badge
	mongoCursor, err := s.getCollection(CollectionBadges).Find(ctx, qb.Build(), options)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err
	}
	err = mongoCursor.All(ctx, &badges)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err
	}
	count := len(badges)
	if count == 0 {
		return badges, hasNextPage, hasPreviousPage, firstCursor, lastCursor, nil
	}

	// check if the cursor element present, if yes that can be a prev elem
	if c != nil && badges[0].Id == c.ID {
		hasPreviousPage = true
		badges = badges[1:]
		count--
	}

	// check if actual limit +1 elements are there, if yes trim it to limit
	if count >= int(limit)-1 {
		hasNextPage = true
		badges = badges[:limit-2]
		count = len(badges)
	}

	if count > 0 {
		firstCursor = cursor.NewCursor(badges[0].Id, badges[0].TimeCreated).String()
		lastCursor = cursor.NewCursor(badges[count-1].Id, badges[count-1].TimeCreated).String()
	}
	if order < 0 {
		hasNextPage, hasPreviousPage = hasPreviousPage, hasNextPage
		firstCursor, lastCursor = lastCursor, firstCursor
	}
	return badges, hasNextPage, hasPreviousPage, firstCursor, lastCursor, nil
}
func (s *UsersStore) DeleteBadgeByID(ctx context.Context, id string) error {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	if _, err := s.getCollection(CollectionBadges).DeleteOne(ctx, qb.Build()); err != nil {
		return err
	}
	return nil
}
