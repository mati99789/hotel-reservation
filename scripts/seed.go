package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelReservetion/api"
	"hotelReservetion/db"
	"hotelReservetion/db/fixtures"
	"hotelReservetion/types"
	"log"
	"time"
)

var (
	client *mongo.Client
	ctx    = context.Background()
)

func main() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)

	store := db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   hotelStore,
	}

	user := fixtures.AddUser(&store, "matteus.urbaniak@hotmail.com", "Mateusz", "Urbaniak", types.GuestRole)
	fmt.Println("User => ", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(&store, "admin@o2.pl", "Admin", "Admin", types.AdminRole)
	fmt.Println("Admin => ", api.CreateTokenFromUser(admin))

	hotel := fixtures.AddHotel(&store, "Bellucia", "France", 2, nil)
	room := fixtures.AddRoom(&store, types.Normal, false, 129, hotel.ID)
	booking := fixtures.AddBooking(&store, room.ID, user.ID, 2, time.Now(), time.Now().AddDate(0, 0, 2))
	fmt.Println("Booking ->", booking.ID)
}
