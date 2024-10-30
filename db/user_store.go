package db

import (
	"hotelReservetion/types"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
	GetUserById(string) (*types.User, error)
}

type MongoUserStore struct {
	db *mongo.Client
}

func NewMongoUserStore(db *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		db: db,
	}
}

func (h *MongoUserStore) GetUserById(id string) (*types.User, error) {
	return nil, nil
}
