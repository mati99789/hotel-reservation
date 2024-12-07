package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelReservetion/db"
	"hotelReservetion/types"
	"log"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedHotel(name, location string, rating int) error {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Size:  types.Small,
			Price: 99,
		},
		{
			Size:  types.Normal,
			Price: 199,
		},
		{
			Size:  types.KingSize,
			Price: 299,
		},
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)

	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelId = insertedHotel.ID
		insertRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(insertRoom)

	}

	fmt.Println(insertedHotel)
	return nil
}

func seedUser(fname, lname, email string) {
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

	_, err = userStore.InsertUser(ctx, user)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	seedHotel("Bellucia", "France", 3)
	seedHotel("Radisson", "Poland", 2)
	seedHotel("Ice Hotel", "Sweden", 5)
	seedUser("James", "Baz", "james@bar.com")
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
}
