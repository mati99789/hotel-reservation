package db

import (
	"context"
	"errors"
	"hotelReservetion/shared"
	"hotelReservetion/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelStore interface {
	Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error)
	Update(context.Context, shared.Map, shared.Map) error
	GetHotels(context.Context, shared.Map, *types.PaginationOptions) ([]*types.Hotel, error)
	GetHotelByID(context.Context, string) (*types.Hotel, error)
	UpdateHotel(context.Context, *types.Hotel, string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{Client: client, Collection: client.Database(DBNAME).Collection("hotels")}
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter shared.Map, pagination *types.PaginationOptions) ([]*types.Hotel, error) {
	opts := options.Find()

	if pagination != nil {
		opts.SetSkip((pagination.Page - 1) * pagination.PageSize)
		opts.SetLimit(pagination.PageSize)

		if pagination.SortBy != "" {
			sortDirection := 1
			if pagination.SortDesc {
				sortDirection = -1
			}
			opts.SetSort(bson.D{{Key: pagination.SortBy, Value: sortDirection}})
		}
	}
	resp, err := s.Collection.Find(ctx, filter, opts)
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
	response, err := s.Collection.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = response.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) Update(ctx context.Context, filter shared.Map, update shared.Map) error {
	_, err := s.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, hotelID string) (*types.Hotel, error) {
	var hotel types.Hotel

	hotelObjectID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return nil, err
	}

	if err := s.Collection.FindOne(ctx, shared.Map{"_id": hotelObjectID}).Decode(&hotel); err != nil {
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
	err = s.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedHotel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("hotel not found")
		}
		return nil, err
	}

	return &updatedHotel, nil
}
