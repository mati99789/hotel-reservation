package api

import (
	"bytes"
	"encoding/json"
	"hotelReservetion/db/fixtures"
	"hotelReservetion/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAuthenticationWrongPassword(t *testing.T) {
	// This test stays mostly the same since we're testing failure case
	tdb := setup()
	defer tdb.tearddown(t)

	fixtures.AddUser(tdb.Store, "test@test.com", "James", "Bond", types.GuestRole)
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
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
	tdb := setup()
	defer tdb.tearddown(t)

	insertedUser := fixtures.AddUser(tdb.Store, "test@test.com", "James", "Bond", types.GuestRole)
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "test@test.com",
		Password: "James_Bond",
	}

	b, _ := json.Marshal(authParams)
	request := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewReader(b))
	request.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(request)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK but got %v", resp.StatusCode)
	}

	// Get cookies from response
	cookies := resp.Cookies()
	var accessToken string
	var refreshToken string

	for _, cookie := range cookies {
		switch cookie.Name {
		case "access_token":
			accessToken = cookie.Value
		case "refresh_token":
			refreshToken = cookie.Value
		}
	}

	if accessToken == "" {
		t.Error("expected access_token cookie to be present")
	}

	if refreshToken == "" {
		t.Error("expected refresh_token cookie to be present")
	}

	// Parse the response for user data
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	// Verify user data
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

	// Validate the access token
	claims, err := ValidateToken(accessToken, types.AccessToken)
	if err != nil {
		t.Errorf("error validating token: %v", err)
	}

	if claims["user_id"] != insertedUser.ID.Hex() {
		t.Errorf("got user_id %v, want %v", claims["user_id"], insertedUser.ID)
	}

	// Optionally validate refresh token
	refreshClaims, err := ValidateToken(refreshToken, types.RefreshToken)
	if err != nil {
		t.Errorf("error validating refresh token: %v", err)
	}

	if refreshClaims["type"] != string(types.RefreshToken) {
		t.Errorf("expected refresh token type but got %v", refreshClaims["type"])
	}
}
