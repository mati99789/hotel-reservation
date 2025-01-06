package db

import "go.mongodb.org/mongo-driver/mongo"

func NewMongoBookingStoreWithDB(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client:     client,
		collection: client.Database(TESTDBName).Collection("bookings"),
	}
}
