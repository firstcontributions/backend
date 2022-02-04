package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	DBUsers          = "users"
	CollectionUsers  = "users"
	CollectionBadges = "badges"
	CollectionTokens = "tokens"
)

type UsersStore struct {
	client *mongo.Client
}

// NewUsersStore makes connection to mongo server by provided url
// and return an instance of the client
func NewUsersStore(ctx context.Context, url string) (*UsersStore, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return &UsersStore{
		client: client,
	}, nil
}

func (s *UsersStore) getCollection(collection string) *mongo.Collection {
	return s.client.Database(DBUsers).Collection(collection)
}

func (s *UsersStore) getPagination(
	ctx context.Context,
	collection string,
	query bson.M,
	revQuery bson.M,
	limit *int64,
) (
	*bool,
	*bool,
	error,
) {
	var limitVal int64 = 10
	if limit != nil {
		limitVal = int64(*limit)
	}
	count, err := s.getCollection(collection).CountDocuments(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	hasNextPage := count > limitVal

	limitOne := int64(1)
	options := &options.FindOptions{
		Limit: &limitOne,
	}

	cursor, err := s.getCollection(collection).Find(ctx, revQuery, options)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)
	hasPreviousPage := cursor.Next(ctx)
	return &hasNextPage, &hasPreviousPage, nil
}
