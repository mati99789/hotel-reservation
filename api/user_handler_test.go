package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hotelReservetion/db"
	"hotelReservetion/types"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburi  = "mongodb://localhost:27017"
	dbName = "hotel-reservation-test"
)

type Testdb struct {
	db.UserStore
}

func setup(t *testing.T) *Testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	return &Testdb{
		db.NewMongoUserStore(client, dbName),
	}
}

func (d *Testdb) tearddown(t *testing.T) {
	if err := d.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func TestPostUser(t *testing.T) {

	db := setup(t)

	defer db.tearddown(t)

	app := fiber.New()
	userHandler := NewUserHandler(db.UserStore)

	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	fmt.Println(resp.Status)

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	if len(user.ID) == 0 {
		t.Errorf("Expecting a user id to be set")
	}

	if len(user.EncryptedPassword) > 0 {
		t.Errorf("Expecting the ecryptedpassword not to be included in the json response.")
	}

	if user.FirstName != params.FirstName {
		t.Errorf("First name should be %s", params.FirstName)
	}

	if user.LastName != params.LastName {
		t.Errorf("Last name should be %s", params.LastName)
	}

	if user.Email != params.Email {
		t.Errorf("Email should be %s", params.Email)
	}
}
