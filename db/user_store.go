package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelReservetion/types"
	"log"
)

const userColl = "users"

type UserStore interface {
	GetUserById(context.Context, string) (*types.User, error)
}

type MongoUserStore struct {
	db   *mongo.Client
	coll *mongo.Collection
}

func NewMongoUserStore(db *mongo.Client) *MongoUserStore {
	coll := db.Database(DBNAME).Collection(userColl)
	return &MongoUserStore{
		db:   db,
		coll: coll,
	}
}

func (h *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user types.User
	if err := h.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		log.Printf("Error finding user: %v", err)
		return nil, err
	}

	return &user, nil
}
