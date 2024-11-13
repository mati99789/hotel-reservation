package main

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelReservetion/api"
	"hotelReservetion/db"
	"log"
)

const (
	dburi    = "mongodb://localhost:27017"
	dbname   = "hotel-reservation-db"
	userColl = "users"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	// package for command-line argument parsing.
	listenAddr := flag.String("listenAddr", ":8080", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	// DATABASE INIT
	mongoDBInit := db.NewMongoUserStore(client)

	// Handlers initialization
	userHandler := api.NewUserHandler(mongoDBInit)

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	log.Fatal(app.Listen(*listenAddr))
}
