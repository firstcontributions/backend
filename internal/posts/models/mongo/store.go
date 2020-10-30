package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DBPosts         = "posts"
	CollectionPosts = "posts"
)

type Store struct {
	*mongo.Client
}

func NewStore(client *mongo.Client) *Store {
	return &Store{
		Client: client,
	}
}

func (s *Store) Collection(c string) *mongo.Collection {
	return s.Database(DBPosts).Collection(c)
}
