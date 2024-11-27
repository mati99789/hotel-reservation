package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelReservetion/types"
)

type HotelStore interface {
	Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error)
	Update(context.Context, bson.M, bson.M) error
	GetHotels(context.Context, bson.M) ([]*types.Hotel, error)
	GetHotelByID(context.Context, string) (*types.Hotel, error)
	UpdateHotel(context.Context, *types.Hotel, string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{client: client, collection: client.Database(DBNAME).Collection("hotels")}
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error) {
	resp, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var hotels []*types.Hotel
	if err := resp.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, nil
}

func (s *MongoHotelStore) Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	response, err := s.collection.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = response.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, hotelID string) (*types.Hotel, error) {
	var hotel types.Hotel

	hotelObjectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return nil, err
	}

	if err := s.collection.FindOne(ctx, bson.M{"_id": hotelObjectID}).Decode(&hotel); err != nil {
		return nil, err
	}

	return &hotel, nil
}

func (s *MongoHotelStore) UpdateHotel(ctx context.Context, hotel *types.Hotel, hotelID string) (*types.Hotel, error) {
	objectId, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": hotel}

	// Create options to return the updated document
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedHotel types.Hotel
	err = s.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedHotel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("hotel not found")
		}
		return nil, err
	}

	return &updatedHotel, nil
}
