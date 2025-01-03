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
	db.UserStore
}

func setup(t *testing.T) *Testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	return &Testdb{
		db.NewMongoUserStore(client),
	}
}

func (d *Testdb) tearddown(t *testing.T) {
	if err := d.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
