package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// DBProfile saves profile related informations
	DBProfile         = "profile"
	CollectionProfile = "profile"
	CollectionTokens  = "tokens"
)

func getCollection(conn *mongo.Client) {

}
