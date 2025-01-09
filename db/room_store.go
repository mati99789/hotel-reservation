package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelReservetion/shared"
	"hotelReservetion/types"
)

type RoomStore interface {
	InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error)
	GetRooms(ctx context.Context, m shared.Map) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client     *mongo.Client
	collection *mongo.Collection

	HotelStore HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		collection: client.Database(DBNAME).Collection("rooms"),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter shared.Map) ([]*types.Room, error) {
	response, error := s.collection.Find(ctx, filter)
	if error != nil {
		return nil, error
	}

	var rooms []*types.Room

	if err := response.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	response, err := s.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = response.InsertedID.(primitive.ObjectID)

	filter := shared.Map{"_id": room.HotelId}
	update := shared.Map{"$push": bson.M{"rooms": room.ID}}

	if err := s.HotelStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}

	return room, nil
}
