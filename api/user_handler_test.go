package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"hotelReservetion/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
