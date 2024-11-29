package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Booking struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	RoomID     primitive.ObjectID `json:"room_id" bson:"room_id,omitempty"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	NumPersons int                `json:"num_persons" bson:"num_persons,omitempty"`
	From       time.Time          `json:"from" bson:"from,omitempty"`
	To         time.Time          `json:"to" bson:"to,omitempty"`
}
