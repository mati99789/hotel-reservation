package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelReservetion/types"
)

type BookingStore interface {
	InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error)
	UpdateBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error)
	GetBookings(ctx context.Context) (*types.Booking, error)
}

type MongoBookingStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client:     client,
		collection: client.Database(DBNAME).Collection("bookings"),
	}
}

func (h *MongoBookingStore) GetBookings(ctx context.Context) ([]*types.Booking, error) {
	var bookings []*types.Booking

	cursors, err := h.collection.Find(ctx, nil)
	if err != nil {
		return nil, err
	}

	if err := cursors.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil

}