package main

import (
	"context"
	"flag"
	"hotelReservetion/api"
	"hotelReservetion/db"
	"hotelReservetion/types"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	// package for command-line argument parsing.
	listenAddr := flag.String("listenAddr", ":8080", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		userStore    = db.NewMongoUserStore(client)
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookingStore(client)
	)

	var (
		store = db.Store{
			User:    userStore,
			Room:    roomStore,
			Hotel:   hotelStore,
			Booking: bookingStore,
		}
	)

	// Handlers initialization
	var (
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(&store)
		authHandler    = api.NewAuthHandler(userStore)
		bookingHandler = api.NewBookingHandler(&store)
		roomHandler    = api.NewRoomHandler(&store)
	)

	app := fiber.New(config)
	auth := app.Group("/api")

	// Public routes (no authentication needed)
	auth.Post("/auth", authHandler.HandleAuthenticate)
	auth.Post("/register", userHandler.HandlePostUser)
	auth.Post("/refresh", authHandler.HandleRefresh)

	// Protected routes require authentication
	apiv1 := app.Group("/api/v1", api.JWTAuthentications)

	// Admin routes
	admin := apiv1.Group("/admin", api.AuthorizeRole(types.AdminRole))
	admin.Get("/users", userHandler.HandleGetUsers)

	// user handler
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	//hotel handlers
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotelByID)
	apiv1.Post("/hotel/:id", hotelHandler.HandleHotelUpdate)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// Rooms handlers
	apiv1.Get("/rooms", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// Bookings handlers
	apiv1.Get("/booking", bookingHandler.HandleGetBookings)
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	log.Fatal(app.Listen(*listenAddr))
}
