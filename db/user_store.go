package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelReservetion/types"
	"log"
)

const userColl = "users"

type UserStore interface {
	GetUserById(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUserById(context.Context, string) (*types.User, error)
	UpdateUserById(ctx context.Context, filter bson.M, update bson.M) error
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

func (h *MongoUserStore) UpdateUserById(ctx context.Context, filter bson.M, update bson.M) error {
	updateDoc := bson.M{
		"$set": update,
	}

	_, err := h.coll.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}

	return err
}

func (h *MongoUserStore) DeleteUserById(ctx context.Context, userId string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	res, err := h.coll.DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		return nil, err
	}

	if res.DeletedCount == 0 {
		return nil, errors.New("user not found")
	}

	return nil, nil
}

func (h *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := h.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (h *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {

	cur, err := h.coll.Find(ctx, bson.M{})

	defer cur.Close(ctx)

	if err != nil {
		return nil, err
	}

	var users []*types.User

	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
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
