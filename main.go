package main

import (
	"flag"
	"github.com/gofiber/fiber/v2"
	"hotelReservetion/api"
	"log"
)

func main() {

	// package for command-line argument parsing.
	listenAddr := flag.String("listenAddr", ":8080", "The listen address of the API server")
	flag.Parse()

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/users", api.HandleGetUsers)
	apiv1.Get("/user/:id", api.HandleGetUser)
	log.Fatal(app.Listen(*listenAddr))
}
