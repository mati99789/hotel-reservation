package db

import "go.mongodb.org/mongo-driver/mongo"

func NewMongoBookingStoreWithDB(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client:     client,
		collection: client.Database(TESTDBName).Collection("bookings"),
	}
}

func NewMongoHotelStoreWithDB(client *mongo.Client, dbName string) *MongoHotelStore {
	return &MongoHotelStore{
		Client:     client,
		Collection: client.Database(dbName).Collection("hotels"),
	}
}

func NewMongoRoomStoreWithDB(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		collection: client.Database(TESTDBName).Collection("rooms"),
		HotelStore: hotelStore,
	}
}

func NewMongoUserStoreWithDB(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		db:   client,
		coll: client.Database(TESTDBName).Collection("users"),
	}
}
