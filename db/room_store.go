package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelReservetion/types"
)

type RoomStore interface {
	InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client     *mongo.Client
	collection *mongo.Collection

	HotelStore HotelStore
}

func NewMongoRoomStore(client *mongo.Client, dbName string) *MongoRoomStore {
	return &MongoRoomStore{client: client, collection: client.Database(dbName).Collection("rooms")}
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	response, err := s.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = response.InsertedID.(primitive.ObjectID)

	filter := bson.M{"_id": room.ID}
	update := bson.M{"$push": bson.M{"rooms": room}}

	return room, nil
}