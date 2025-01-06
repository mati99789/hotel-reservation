package db

import "go.mongodb.org/mongo-driver/mongo"

func NewMongoUserStoreWithDB(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		db:   client,
		coll: client.Database(TESTDBName).Collection("users"),
	}
}
