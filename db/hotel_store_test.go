package db

import "go.mongodb.org/mongo-driver/mongo"

func NewMongoHotelStoreWithDB(client *mongo.Client, dbName string) *MongoHotelStore {
	return &MongoHotelStore{
		Client:     client,
		Collection: client.Database(dbName).Collection("hotels"),
	}
}
