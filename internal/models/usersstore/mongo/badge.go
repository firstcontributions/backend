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

func badgeFiltersToQuery(filters *usersstore.BadgeFilters) *mongoqb.QueryBuilder {
	qb := mongoqb.NewQueryBuilder()
	if len(filters.Ids) > 0 {
		qb.In("_id", filters.Ids)
	}
	if filters.User != nil {
		qb.Eq("user_id", filters.User.Id)
	}
	return qb
}
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

func (s *UsersStore) GetOneBadge(ctx context.Context, filters *usersstore.BadgeFilters) (*usersstore.Badge, error) {
	qb := badgeFiltersToQuery(filters)
	var badge usersstore.Badge
	if err := s.getCollection(CollectionBadges).FindOne(ctx, qb.Build()).Decode(&badge); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &badge, nil
}

func (s *UsersStore) CountBadges(ctx context.Context, filters *usersstore.BadgeFilters) (
	int64,
	error,
) {
	qb := badgeFiltersToQuery(filters)

	count, err := s.getCollection(CollectionBadges).CountDocuments(ctx, qb.Build())
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *UsersStore) GetBadges(
	ctx context.Context,
	filters *usersstore.BadgeFilters,
	after *string,
	before *string,
	first *int64,
	last *int64,
	sortBy usersstore.BadgeSortBy,
	sortOrder *string,
) (
	[]*usersstore.Badge,
	bool,
	bool,
	[]string,
	error,
) {
	qb := badgeFiltersToQuery(filters)
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
						Eq(usersstore.BadgeSortBy(c.SortBy).String(), c.OffsetValue).
						Gt("_id", c.ID),
					mongoqb.NewQueryBuilder().
						Gt(usersstore.BadgeSortBy(c.SortBy).String(), c.OffsetValue),
				)
			} else {
				qb.Or(
					mongoqb.NewQueryBuilder().
						Eq(usersstore.BadgeSortBy(c.SortBy).String(), c.OffsetValue).
						Lt("_id", c.ID),
					mongoqb.NewQueryBuilder().
						Lt(usersstore.BadgeSortBy(c.SortBy).String(), c.OffsetValue),
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

	var badges []*usersstore.Badge
	mongoCursor, err := s.getCollection(CollectionBadges).Find(ctx, qb.Build(), options)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, nil, err
	}
	err = mongoCursor.All(ctx, &badges)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, nil, err
	}
	count := len(badges)
	if count == 0 {
		return badges, hasNextPage, hasPreviousPage, nil, nil
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

	cursors := make([]string, count)
	for i, badge := range badges {
		cursors[i] = cursor.NewCursor(badge.Id, uint8(sortBy), badge.Get(sortBy.String()), sortBy.CursorType()).String()
	}

	if paginationSortOrder < 0 {
		hasNextPage, hasPreviousPage = hasPreviousPage, hasNextPage
		badges = utils.ReverseList(badges)
	}
	return badges, hasNextPage, hasPreviousPage, cursors, nil
}

func (s *UsersStore) UpdateBadge(ctx context.Context, id string, badgeUpdate *usersstore.BadgeUpdate) error {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)

	now := time.Now()
	badgeUpdate.TimeUpdated = &now

	u := mongoqb.NewUpdateMap().
		SetFields(badgeUpdate)

	um, err := u.BuildUpdate()
	if err != nil {
		return err
	}
	if _, err := s.getCollection(CollectionBadges).UpdateOne(ctx, qb.Build(), um); err != nil {
		return err
	}
	return nil
}

func (s *UsersStore) DeleteBadgeByID(ctx context.Context, id string) error {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	if _, err := s.getCollection(CollectionBadges).DeleteOne(ctx, qb.Build()); err != nil {
		return err
	}
	return nil
}
