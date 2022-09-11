package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/firstcontributions/backend/internal/models/utils"
	"github.com/firstcontributions/backend/pkg/authorizer"
	"github.com/firstcontributions/backend/pkg/cursor"
	"github.com/gokultp/go-mongoqb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func userFiltersToQuery(filters *usersstore.UserFilters) *mongoqb.QueryBuilder {
	qb := mongoqb.NewQueryBuilder()
	if len(filters.Ids) > 0 {
		qb.In("_id", filters.Ids)
	}
	if filters.Handle != nil {
		qb.Eq("handle", filters.Handle)
	}
	if filters.Search != nil {
		qb.Search(*filters.Search)
	}
	return qb
}
func (s *UsersStore) CreateUser(ctx context.Context, user *usersstore.User, ownership *authorizer.Scope) (*usersstore.User, error) {
	now := time.Now()
	user.TimeCreated = now
	user.TimeUpdated = now
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	user.Id = uuid.String()
	user.Permissions = []authorizer.Permission{
		{
			Role:  "admin",
			Scope: authorizer.Scope{Users: []string{user.Id}},
		},
	}
	user.Ownership = &authorizer.Scope{Users: []string{user.Id}}
	if _, err := s.getCollection(CollectionUsers).InsertOne(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UsersStore) GetUserByID(ctx context.Context, id string) (*usersstore.User, error) {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	var user usersstore.User
	if err := s.getCollection(CollectionUsers).FindOne(ctx, qb.Build()).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *UsersStore) GetOneUser(ctx context.Context, filters *usersstore.UserFilters) (*usersstore.User, error) {
	qb := userFiltersToQuery(filters)
	var user usersstore.User
	if err := s.getCollection(CollectionUsers).FindOne(ctx, qb.Build()).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *UsersStore) CountUsers(ctx context.Context, filters *usersstore.UserFilters) (
	int64,
	error,
) {
	qb := userFiltersToQuery(filters)

	count, err := s.getCollection(CollectionUsers).CountDocuments(ctx, qb.Build())
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *UsersStore) GetUsers(
	ctx context.Context,
	filters *usersstore.UserFilters,
	after *string,
	before *string,
	first *int64,
	last *int64,
	sortBy usersstore.UserSortBy,
	sortOrder *string,
) (
	[]*usersstore.User,
	bool,
	bool,
	[]string,
	error,
) {
	qb := userFiltersToQuery(filters)
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
						Eq(usersstore.UserSortBy(c.SortBy).String(), c.OffsetValue).
						Gt("_id", c.ID),
					mongoqb.NewQueryBuilder().
						Gt(usersstore.UserSortBy(c.SortBy).String(), c.OffsetValue),
				)
			} else {
				qb.Or(
					mongoqb.NewQueryBuilder().
						Eq(usersstore.UserSortBy(c.SortBy).String(), c.OffsetValue).
						Lt("_id", c.ID),
					mongoqb.NewQueryBuilder().
						Lt(usersstore.UserSortBy(c.SortBy).String(), c.OffsetValue),
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

	var users []*usersstore.User
	mongoCursor, err := s.getCollection(CollectionUsers).Find(ctx, qb.Build(), options)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, nil, err
	}
	err = mongoCursor.All(ctx, &users)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, nil, err
	}
	count := len(users)
	if count == 0 {
		return users, hasNextPage, hasPreviousPage, nil, nil
	}

	// check if the cursor element present, if yes that can be a prev elem
	if c != nil && users[0].Id == c.ID {
		hasPreviousPage = true
		users = users[1:]
		count--
	}

	// check if actual limit +1 elements are there, if yes trim it to limit
	if count >= int(limit)-1 {
		hasNextPage = true
		users = users[:limit-2]
		count = len(users)
	}

	cursors := make([]string, count)
	for i, user := range users {
		cursors[i] = cursor.NewCursor(user.Id, uint8(sortBy), user.Get(sortBy.String()), sortBy.CursorType()).String()
	}

	if paginationSortOrder < 0 {
		hasNextPage, hasPreviousPage = hasPreviousPage, hasNextPage
		users = utils.ReverseList(users)
	}
	return users, hasNextPage, hasPreviousPage, cursors, nil
}

func (s *UsersStore) UpdateUser(ctx context.Context, id string, userUpdate *usersstore.UserUpdate) error {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)

	now := time.Now()
	userUpdate.TimeUpdated = &now

	u := mongoqb.NewUpdateMap().
		SetFields(userUpdate)

	um, err := u.BuildUpdate()
	if err != nil {
		return err
	}
	if _, err := s.getCollection(CollectionUsers).UpdateOne(ctx, qb.Build(), um); err != nil {
		return err
	}
	return nil
}

func (s *UsersStore) DeleteUserByID(ctx context.Context, id string) error {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	if _, err := s.getCollection(CollectionUsers).DeleteOne(ctx, qb.Build()); err != nil {
		return err
	}
	return nil
}
