package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelReservetion/db"
	"hotelReservetion/types"
	"log"
	"time"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	bookStore  db.BookingStore
	ctx        = context.Background()
)

func seedBooking(roomID, userID primitive.ObjectID, from, till time.Time, numPerson int) *types.Booking {
	booking := &types.Booking{
		RoomID:     roomID,
		UserID:     userID,
		NumPersons: numPerson,
		From:       from,
		To:         till,
		Canceled:   false,
	}

	insertedBooking, err := bookStore.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking
}

func seedHotel(name, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)

	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func seedUser(fname, lname, email string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  "supersecurepassword",
		Role:      types.GuestRole,
	})

	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := userStore.InsertUser(ctx, user)

	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}

func seedRoom(size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    types.RoomType(size),
		Seaside: ss,
		Price:   price,
		HotelId: hotelID,
	}

	insertedRoom, err := roomStore.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom

}

func main() {
	james := seedUser("James", "Baz", "james@bar.com")

	seedHotel("Bellucia", "France", 3)
	seedHotel("Radisson", "Poland", 2)
	hotel := seedHotel("Ice Hotel", "Sweden", 5)
	seedRoom("small", true, 89.99, hotel.ID)
	seedRoom("medium", false, 129.99, hotel.ID)
	room := seedRoom("large", true, 289.99, hotel.ID)

	seedBooking(room.ID, james.ID, time.Now(), time.Now().Add(time.Hour*24*2), 2)
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookStore = db.NewMongoBookingStore(client)
}
