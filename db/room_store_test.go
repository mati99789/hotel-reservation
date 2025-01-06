package db

import "go.mongodb.org/mongo-driver/mongo"

func NewMongoRoomStoreWithDB(client *mongo.Client, hotelStore HotelStore) *RoomStore {
	return &MongoRoomStore{
		client:     client,
		collection: client.Database(TESTDBName).Collection("rooms"),
		HotelStore: hotelStore,
	}
}
