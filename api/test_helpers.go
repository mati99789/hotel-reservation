package api

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelReservetion/db"
	"log"
	"testing"
)

type Testdb struct {
	client *mongo.Client
	*db.Store
}

func setup() *Testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := &db.MongoHotelStore{
		Client:     client,
		Collection: client.Database(db.TESTDBName).Collection("hotels"),
	}

	return &Testdb{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStoreWithDB(client),
			Booking: db.NewMongoBookingStoreWithDB(client),
			Room:    db.NewMongoRoomStoreWithDB(client, hotelStore),
			Hotel:   hotelStore,
		},
	}
}

func (d *Testdb) tearddown(t *testing.T) {
	if err := d.client.Database(db.TESTDBName).Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
}
