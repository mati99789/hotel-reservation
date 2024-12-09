package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	RoomID     primitive.ObjectID `bson:"room_id" json:"room_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	NumPersons int                `bson:"num_persons" json:"num_persons"`
	From       time.Time          `bson:"from" json:"from"`
	To         time.Time          `bson:"to" json:"to"`
	Canceled   bool               `bson:"canceled" json:"canceled"`
}
