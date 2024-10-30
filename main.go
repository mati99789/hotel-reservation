package main

import (
	"context"
	"flag"
	"hotelReservetion/api"
	"hotelReservetion/db"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburi    = "mongodb://localhost:27017"
	dbname   = "hotel-reservertion"
	userColl = "users"
)

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

	// Handlers initilizaion
	userHandler := api.NewUserHandler(mongoDBInit)

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	log.Fatal(app.Listen(*listenAddr))
}
