package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"hotelReservetion/api/middleware"
	"hotelReservetion/db"
	"hotelReservetion/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticationWrongPassword(t *testing.T) {
	tdb := setup(t)

	defer tdb.tearddown(t)

	insertTestUser(t, tdb)

	app := fiber.New()

	authHandler := NewAuthHandler(tdb.UserStore)

	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "test@test.com",
		Password: "wrongpassword",
	}

	b, _ := json.Marshal(authParams)

	request := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	request.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(request)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var genericResp genericResp

	if err := json.NewDecoder(resp.Body).Decode(&genericResp); err != nil {
		t.Error(err)
	}

	if genericResp.Type != "error" {
		t.Errorf("got %s, want %s", genericResp.Type, "error")
	}

	if genericResp.Msg != "Invalid Credentials" {
		t.Errorf("got %s, want %s", genericResp.Msg, "Invalid Credentials")
	}
}

func TestAuthenticationSuccessPassword(t *testing.T) {
	tdb := setup(t)

	defer tdb.tearddown(t)

	insertedUser := insertTestUser(t, tdb)

	app := fiber.New()

	authHandler := NewAuthHandler(tdb.UserStore)

	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "test@test.com",
		Password: "password",
	}

	b, _ := json.Marshal(authParams)

	request := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	request.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(request)
	if err != nil {
		t.Error(err)
	}

	// Parse the response
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	if authResp.Token == "" {
		t.Error("expected JWT token to be present")
	}

	if authResp.User == nil {
		t.Fatal("expected user to be present")
	}
	if authResp.User.Email != insertedUser.Email {
		t.Errorf("got email %s, want %s", authResp.User.Email, insertedUser.Email)
	}
	if authResp.User.FirstName != insertedUser.FirstName {
		t.Errorf("got first name %s, want %s", authResp.User.FirstName, insertedUser.FirstName)
	}
	if authResp.User.LastName != insertedUser.LastName {
		t.Errorf("got last name %s, want %s", authResp.User.LastName, insertedUser.LastName)
	}

	claims, err := middleware.ValidateToken(authResp.Token)

	if err != nil {
		t.Errorf("error validating token: %v", err)
	}

	if claims["email"] != insertedUser.Email {
		t.Errorf("got email %s, want %s", claims["email"], insertedUser.Email)
	}
}

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "James",
		LastName:  "Bond",
		Email:     "test@test.com",
		Password:  "password",
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = userStore.InsertUser(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}

	return user
}
