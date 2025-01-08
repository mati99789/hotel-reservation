package api

import (
	"encoding/json"
	"hotelReservetion/db/fixtures"
	"hotelReservetion/types"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestAdminGetBookings(t *testing.T) {
	db := setup()
	defer db.tearddown(t)

	var (
		admin   = fixtures.AddUser(db.Store, "admin@o2.pl", "Admin", "Admin", types.AdminRole)
		user    = fixtures.AddUser(db.Store, "user@o2.pl", "User 1", "User 1", types.GuestRole)
		hotel   = fixtures.AddHotel(db.Store, "Bellucia", "France", 2, nil)
		room    = fixtures.AddRoom(db.Store, types.Normal, false, 300, hotel.ID)
		booking = fixtures.AddBooking(db.Store, room.ID, admin.ID, 2, time.Now(), time.Now().AddDate(0, 0, 2))

		app        = fiber.New()
		apiv1      = app.Group("/api/v1", JWTAuthentications)
		adminGroup = apiv1.Group("/admin", AuthorizeRole(types.AdminRole))
	)

	bookingHandler := NewBookingHandler(db.Store)
	adminGroup.Get("/users", bookingHandler.HandleGetBookings)

	token := CreateTokenFromUser(admin)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
	req.Header.Set("X-API-Token", token)
	req.Header.Set("Content-Type", "application/json")

	response, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		t.Fatalf("Expected status code %d, got %d. Response body: %s", http.StatusOK, response.StatusCode, string(body))
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(response.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("Expected 1 booking, got %d", len(bookings))
	}

	if reflect.DeepEqual(booking, bookings[0]) {
		t.Fatalf("Expected %v, got %v", bookings[0], bookings[0])
	}

	// Test non-admin cannot access the bookings
	req = httptest.NewRequest(http.MethodGet, "/api/v1/admin/bookings", nil)
	req.Header.Set("X-API-Token", CreateTokenFromUser(user))

	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode == http.StatusUnauthorized {
		t.Fatalf("Expected status unauthorize but got %d", res.StatusCode)
	}
}
