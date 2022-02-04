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
	if cursorStr != nil {
		c := cursor.FromString(*cursorStr)
		if c != nil {
			if order == 1 {
				qb.Lt("time_created", c.TimeStamp)
				qb.Lt("_id", c.ID)
			} else {
				qb.Gt("time_created", c.TimeStamp)
				qb.Gt("_id", c.ID)
			}
		}
	}
	sortOrder := utils.GetSortOrder(order)

	options := &options.FindOptions{
		Limit: &limit,
		Sort:  sortOrder,
	}

	var firstCursor, lastCursor string

	var badges []*usersstore.Badge
	mongoCursor, err := s.getCollection(CollectionBadges).Find(ctx, qb.Build(), options)
	if err != nil {
		return nil, false, false, firstCursor, lastCursor, err
	}
	err = mongoCursor.All(ctx, &badges)
	if err != nil {
		return nil, false, false, firstCursor, lastCursor, err
	}
	count := len(badges)
	if count > 0 {
		firstCursor = cursor.NewCursor(badges[0].Id, badges[0].TimeCreated).String()
		lastCursor = cursor.NewCursor(badges[count-1].Id, badges[count-1].TimeCreated).String()
	}
	hasNextPage, hasPreviousPage := utils.CheckHasNextPrevPages(count, int(limit), order)
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
