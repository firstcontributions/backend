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

func (s *UsersStore) CreateUser(ctx context.Context, user *usersstore.User) (*usersstore.User, error) {
	now := time.Now()
	user.TimeCreated = now
	user.TimeUpdated = now
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	user.Id = uuid.String()
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

func (s *UsersStore) GetUsers(
	ctx context.Context,
	ids []string,
	search *string,
	handle *string,
	after *string,
	before *string,
	first *int64,
	last *int64,
) (
	[]*usersstore.User,
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
	if handle != nil {
		qb.Eq("handle", handle)
	}
	if search != nil {
		qb.Search(*search)
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

	var users []*usersstore.User
	mongoCursor, err := s.getCollection(CollectionUsers).Find(ctx, qb.Build(), options)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err
	}
	err = mongoCursor.All(ctx, &users)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err
	}
	count := len(users)
	if count == 0 {
		return users, hasNextPage, hasPreviousPage, firstCursor, lastCursor, nil
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

	if count > 0 {
		firstCursor = cursor.NewCursor(users[0].Id, users[0].TimeCreated).String()
		lastCursor = cursor.NewCursor(users[count-1].Id, users[count-1].TimeCreated).String()
	}
	if order < 0 {
		hasNextPage, hasPreviousPage = hasPreviousPage, hasNextPage
		firstCursor, lastCursor = lastCursor, firstCursor
	}
	return users, hasNextPage, hasPreviousPage, firstCursor, lastCursor, nil
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
