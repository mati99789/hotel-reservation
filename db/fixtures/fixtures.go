package fixtures

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelReservetion/db"
	"hotelReservetion/types"
	"log"
	"time"
)

func AddUser(store *db.Store, email, fname, lname string, role types.UserRole) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  fmt.Sprintf("%s_%s", fname, lname),
		Role:      role,
	})

	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := store.User.InsertUser(context.TODO(), user)

	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}

func AddHotel(store *db.Store, name, loc string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomsIDS = rooms
	if rooms == nil {
		roomsIDS = []primitive.ObjectID{}
	}

	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomsIDS,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.Insert(context.TODO(), &hotel)

	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddRoom(store *db.Store, size types.RoomType, seaSide bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: seaSide,
		Price:   price,
		HotelId: hotelID,
	}

	room.HotelId = hotelID
	insertRoom, err := store.Room.InsertRoom(context.TODO(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertRoom

}

func AddBooking(store *db.Store, roomID, userID primitive.ObjectID, numPerson int, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		RoomID:     roomID,
		UserID:     userID,
		From:       from,
		To:         till,
		NumPersons: numPerson,
		Canceled:   false,
	}

	insertedBooking, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking
}
