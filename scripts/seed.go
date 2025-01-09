package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelReservetion/db"
	"hotelReservetion/db/fixtures"
	"hotelReservetion/types"
	"log"
	"math/rand"
	"strconv"
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
	tokens, err := types.CreateTokenPair(user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User Access Token => ", tokens.AccessToken)
	fmt.Println("User Refresh Token => ", tokens.RefreshToken)

	admin := fixtures.AddUser(&store, "admin@o2.pl", "Admin", "Admin", types.AdminRole)
	adminTokens, err := types.CreateTokenPair(admin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Admin Access Token => ", adminTokens.AccessToken)
	fmt.Println("Admin Refresh Token => ", adminTokens.RefreshToken)

	hotel := fixtures.AddHotel(&store, "Bellucia", "France", 2, nil)
	room := fixtures.AddRoom(&store, types.Normal, false, 129, hotel.ID)
	booking := fixtures.AddBooking(&store, room.ID, user.ID, 2, time.Now(), time.Now().AddDate(0, 0, 2))
	fmt.Println("Booking ->", booking.ID)

	for i := 0; i < 100; i++ {
		hotelName := "Hotel " + strconv.Itoa(i)
		location := "Location " + strconv.Itoa(i)
		random := rand.Intn(5) + 1

		fixtures.AddHotel(&store, hotelName, location, random, nil)
	}
}
